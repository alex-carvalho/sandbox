package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Printf("\033[32m[DISCOVER] [%s] joined cluster (addr: %s:%d)\033[0m\n", node.Name, node.Addr, node.Port)
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Printf("\033[31m[DEAD] [%s] has left or failed\033[0m\n", node.Name)
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Printf("\033[90m[UPDATE] [%s] updated state\033[0m\n", node.Name)
}

func main() {
	port := flag.Int("port", 9001, "Port to bind memberlist on")
	peersStr := flag.String("peers", "", "Comma-separated seed peers (e.g. 127.0.0.1:9001)")
	nodeName := flag.String("name", "", "Unique name for this node")
	bindAddr := flag.String("bind", "127.0.0.1", "Address to bind memberlist to")
	flag.Parse()

	if *nodeName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Failed to get hostname: %v", err)
		}
		*nodeName = fmt.Sprintf("%s-%d", hostname, *port)
	}

	fmt.Printf("\033[36m[START] memberlist Gossip Node running on %s:%d (Node Name: %s)\033[0m\n\n", *bindAddr, *port, *nodeName)

	config := memberlist.DefaultLocalConfig()
	config.Name = *nodeName
	config.BindAddr = *bindAddr
	config.BindPort = *port
	config.AdvertisePort = *port
	config.Logger = log.New(io.Discard, "", 0) // Silence internal logs
	config.Events = &eventDelegate{}

	list, err := memberlist.Create(config)
	if err != nil {
		log.Fatalf("Failed to create memberlist: %v", err)
	}

	joinPeers(list, *peersStr)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			logMembers(list)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// if not add a graceful shutdown nodes will take more time to detect the other node is down
	gracefulShutdown(list)
}

func joinPeers(list *memberlist.Memberlist, peersStr string) {
	if peersStr == "" {
		return
	}

	var peersToJoin []string
	for _, p := range strings.Split(peersStr, ",") {
		if p = strings.TrimSpace(p); p != "" {
			peersToJoin = append(peersToJoin, p)
		}
	}

	if len(peersToJoin) > 0 {
		fmt.Printf("\033[90m[INFO] Attempting to join seed peers: %s\033[0m\n", strings.Join(peersToJoin, ", "))
		n, err := list.Join(peersToJoin)
		if err != nil {
			log.Fatalf("Failed to join cluster: %v", err)
		}
		fmt.Printf("\033[90m[INFO] Successfully contacted %d nodes\033[0m\n", n)
	}
}

func logMembers(list *memberlist.Memberlist) {
	var members []string
	for _, m := range list.Members() {
		members = append(members, fmt.Sprintf("%s (%s:%d)", m.Name, m.Addr, m.Port))
	}
	fmt.Printf("[MEMBERS] %s\n", strings.Join(members, ", "))
}

func gracefulShutdown(list *memberlist.Memberlist) {
	log.Println("Shutting down node gracefully...")
	if err := list.Leave(time.Second * 5); err != nil {
		log.Printf("Error leaving cluster: %v", err)
	}
	if err := list.Shutdown(); err != nil {
		log.Printf("Error shutting down memberlist: %v", err)
	}
	log.Println("Node stopped.")
}
