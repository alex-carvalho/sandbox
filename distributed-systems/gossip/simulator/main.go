package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	mode := flag.String("mode", "simulate", "Running mode: 'simulate' (dashboard) or 'node' (standalone node CLI)")
	port := flag.Int("port", 9001, "UDP Port to listen on (only in 'node' mode)")
	peersList := flag.String("peers", "", "Comma-separated list of seed UDP peers (only in 'node' mode, e.g. '127.0.0.1:9002,127.0.0.1:9003')")
	httpPort := flag.Int("http-port", 8080, "HTTP port for the dashboard visualizer (only in 'simulate' mode)")
	clusterSize := flag.Int("cluster-size", 5, "Number of nodes to spawn in simulation")
	gossipInterval := flag.Int("gossip-ms", 1000, "Gossip cycle interval in milliseconds")
	suspectTimeout := flag.Int("suspect-ms", 3000, "Heartbeat stall threshold before suspecting a peer (milliseconds)")
	deadTimeout := flag.Int("dead-ms", 6000, "Timeout before declaring suspected peer as dead (milliseconds)")
	fanout := flag.Int("fanout", 2, "Gossip fanout size (number of random peers to gossip to per round)")
	enableUI := flag.Bool("ui", false, "Start visualizer dashboard on the HTTP port in standalone 'node' mode")
	visualizerURL := flag.String("visualizer", "", "URL of remote visualizer server to send telemetry events (e.g. http://localhost:8080)")

	flag.Parse()

	nodeConfig := NodeConfig{
		GossipInterval: time.Duration(*gossipInterval) * time.Millisecond,
		SuspectTimeout: time.Duration(*suspectTimeout) * time.Millisecond,
		DeadTimeout:    time.Duration(*deadTimeout) * time.Millisecond,
		Fanout:         *fanout,
	}

	if *mode == "node" {
		runStandaloneNode(*port, *peersList, nodeConfig, *enableUI, *httpPort, *visualizerURL)
	} else if *mode == "simulate" {
		runSimulation(*httpPort, *clusterSize, nodeConfig)
	} else {
		log.Fatalf("Unknown mode: %s. Use 'simulate' or 'node'.", *mode)
	}
}

func runSimulation(httpPort int, clusterSize int, config NodeConfig) {
	sim := NewSimulator(
		config.GossipInterval,
		config.SuspectTimeout,
		config.DeadTimeout,
		config.Fanout,
	)

	basePort := 10001
	fmt.Printf("Starting Gossip Protocol Simulation with %d nodes on UDP ports %d-%d...\n", clusterSize, basePort, basePort+clusterSize-1)

	if err := sim.StartCluster(basePort, clusterSize); err != nil {
		log.Fatalf("Failed to start cluster simulation: %v", err)
	}
	defer sim.StopCluster()

	// Auto-open browser in background
	go func() {
		time.Sleep(1200 * time.Millisecond)
		url := fmt.Sprintf("http://localhost:%d", httpPort)
		fmt.Printf("\n=======================================================\n")
		fmt.Printf(">>> Web Visualizer dashboard running at: %s <<<\n", url)
		fmt.Printf("=======================================================\n\n")
		openBrowser(url)
	}()

	webServer := NewWebServer(sim, httpPort)
	if err := webServer.Start(); err != nil {
		log.Fatalf("Failed to run Web visualizer server: %v", err)
	}
}

func runStandaloneNode(port int, peersStr string, config NodeConfig, enableUI bool, httpPort int, visualizerURL string) {
	config.Addr = fmt.Sprintf("127.0.0.1:%d", port)

	var webServer *WebServer
	if enableUI {
		webServer = NewWebServer(nil, httpPort)
	}

	config.OnEvent = func(ev NodeEvent) {
		timeStr := time.UnixMilli(ev.Timestamp).Format("15:04:05.000")
		fmt.Printf("[%s] [%s] %s\n", timeStr, ev.Type, ev.Detail)
		if webServer != nil {
			webServer.BroadcastEvent(ev)
		}
		if visualizerURL != "" {
			go reportEventToVisualizer(visualizerURL, ev)
		}
	}

	node, err := NewNode(config, nil, nil)
	if err != nil {
		log.Fatalf("Failed to construct node: %v", err)
	}

	if enableUI {
		webServer.SetNode(node)
		go func() {
			time.Sleep(1000 * time.Millisecond)
			url := fmt.Sprintf("http://localhost:%d", httpPort)
			fmt.Printf("\n=======================================================\n")
			fmt.Printf(">>> Standalone Node Dashboard running at: %s <<<\n", url)
			fmt.Printf("=======================================================\n\n")
			openBrowser(url)
		}()
		go func() {
			if err := webServer.Start(); err != nil {
				log.Printf("Warning: Web Visualizer server failed to start: %v", err)
			}
		}()
	}

	if err := node.Start(); err != nil {
		log.Fatalf("Failed to initialize node UDP socket: %v", err)
	}
	defer node.Stop()

	if peersStr != "" {
		peers := strings.Split(peersStr, ",")
		for _, peer := range peers {
			p := strings.TrimSpace(peer)
			if p != "" {
				node.AddPeer(p)
			}
		}
	}

	fmt.Printf("\nNode running on UDP %s. Commands available:\n", config.Addr)
	fmt.Println("  set <key> <val>   - Insert or update local key-value state")
	fmt.Println("  get <key>         - Retrieve key value from local db")
	fmt.Println("  members           - Print known cluster membership status")
	fmt.Println("  add <addr>        - Connect a new peer address (e.g. 127.0.0.1:9002)")
	fmt.Println("  exit              - Terminate node and shutdown")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		cmd := strings.ToLower(tokens[0])
		switch cmd {
		case "exit":
			return
		case "set":
			if len(tokens) < 3 {
				fmt.Println("Syntax: set <key> <val>")
				continue
			}
			key := tokens[1]
			val := strings.Join(tokens[2:], " ")
			node.Set(key, val)
		case "get":
			if len(tokens) < 2 {
				fmt.Println("Syntax: get <key>")
				continue
			}
			key := tokens[1]
			if v, exists := node.Get(key); exists {
				fmt.Printf("State found: '%s' = '%s' (version: %d, updated: %s)\n",
					key, v.Val, v.Version, time.Unix(0, v.Timestamp).Format("15:04:05.000"))
			} else {
				fmt.Println("Key not found")
			}
		case "members":
			members := node.GetMembers()
			fmt.Printf("Membership Table (%d nodes total):\n", len(members))
			for _, m := range members {
				statusColor := "\033[32m"
				if m.Status == StatusSuspected {
					statusColor = "\033[33m"
				} else if m.Status == StatusDead {
					statusColor = "\033[31m"
				}
				fmt.Printf("  - %s : %s%s\033[0m (heartbeat: %d)\n", m.Addr, statusColor, m.Status, m.HeartbeatCount)
			}
		case "add":
			if len(tokens) < 2 {
				fmt.Println("Syntax: add <addr>")
				continue
			}
			addr := tokens[1]
			node.AddPeer(addr)
			fmt.Printf("Peer %s added to membership list\n", addr)
		default:
			fmt.Println("Unknown command. Options: set, get, members, add, exit")
		}
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Run()
}

func reportEventToVisualizer(visualizerURL string, ev NodeEvent) {
	data, err := json.Marshal(ev)
	if err != nil {
		return
	}
	resp, err := http.Post(visualizerURL+"/api/report-event", "application/json", bytes.NewBuffer(data))
	if err == nil {
		resp.Body.Close()
	}
}
