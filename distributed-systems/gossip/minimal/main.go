package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// PeerState tracks the metadata of a cluster member locally.
type PeerState struct {
	Epoch       int64
	Heartbeat   uint64
	LastUpdated time.Time
	IsFailed    bool
}

// PeerGossip is the epoch-heartbeat tuple sent over the network.
type PeerGossip struct {
	Epoch     int64  `json:"epoch"`
	Heartbeat uint64 `json:"heartbeat"`
}

// GossipPacket is the JSON structure exchanged over UDP.
type GossipPacket struct {
	FromAddr string                `json:"from_addr"`
	Peers    map[string]PeerGossip `json:"peers"`
}

var (
	selfAddr      string
	selfEpoch     int64
	selfHeartbeat uint64
	gossipIter    time.Duration
	failTimeout   time.Duration

	peersMu sync.RWMutex
	peers   = make(map[string]*PeerState)
)

func main() {
	port := flag.Int("port", 9001, "UDP port to bind on (e.g. 9001)")
	peersStr := flag.String("peers", "", "Comma-separated seed peers (e.g. 127.0.0.1:9002)")
	gossipMs := flag.Int("gossip-ms", 2000, "Gossip cycle interval in milliseconds")
	failMs := flag.Int("fail-ms", 6000, "Duration before declaring a peer dead (milliseconds)")
	flag.Parse()

	selfAddr = fmt.Sprintf("127.0.0.1:%d", *port)
	gossipIter = time.Duration(*gossipMs) * time.Millisecond
	failTimeout = time.Duration(*failMs) * time.Millisecond

	// Initialize our local epoch using startup timestamp
	selfEpoch = time.Now().Unix()
	selfHeartbeat = 1

	// Bind UDP listener
	addr, err := net.ResolveUDPAddr("udp", selfAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(fmt.Errorf("failed to listen on UDP %s: %w", selfAddr, err))
	}
	defer conn.Close()

	fmt.Printf("\033[36m[START] Standalone minimal Gossip Node running on UDP %s\033[0m\n", selfAddr)
	fmt.Printf("[INFO] Gossip Interval: %v | Failure Timeout: %v \n\n", gossipIter, failTimeout)

	// Add self to local peer list
	peers[selfAddr] = &PeerState{Epoch: selfEpoch, Heartbeat: selfHeartbeat, LastUpdated: time.Now()}

	// Parse initial bootstrap peers
	if *peersStr != "" {
		for _, peer := range strings.Split(*peersStr, ",") {
			peer = strings.TrimSpace(peer)
			if peer != "" && peer != selfAddr {
				peers[peer] = &PeerState{Epoch: 0, Heartbeat: 0, LastUpdated: time.Now()}
				fmt.Printf("\033[90m[INFO] Added bootstrap peer: %s\033[0m\n", peer)
			}
		}
	}

	// Start threads
	go listenUDP(conn)
	go gossipLoop(conn)
	go failureDetectorLoop()

	// Wait forever
	select {}
}

func listenUDP(conn *net.UDPConn) {
	buf := make([]byte, 65535)
	for {
		size, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var packet GossipPacket
		if err := json.Unmarshal(buf[:size], &packet); err != nil {
			continue
		}

		mergeGossip(packet.FromAddr, packet.Peers)
	}
}

func isNewer(incomingEpoch int64, incomingHb uint64, local *PeerState) bool {
	if incomingEpoch > local.Epoch {
		return true
	}
	if incomingEpoch == local.Epoch && incomingHb > local.Heartbeat {
		return true
	}
	return false
}

func mergeGossip(sender string, incomingPeers map[string]PeerGossip) {
	peersMu.Lock()
	defer peersMu.Unlock()

	now := time.Now()

	for addr, incoming := range incomingPeers {
		// Handle refutation: if a peer claims a higher/equal version status for us,
		// we refute it by incrementing our epoch or advancing our heartbeat count.
		if addr == selfAddr {
			if incoming.Epoch > selfEpoch {
				selfEpoch = incoming.Epoch + 1
				selfHeartbeat = 1
				peers[selfAddr].Epoch = selfEpoch
				peers[selfAddr].Heartbeat = selfHeartbeat
				peers[selfAddr].LastUpdated = now
			}
			continue
		}

		localPeer, exists := peers[addr]
		if !exists {
			// Discovered a new node
			peers[addr] = &PeerState{
				Epoch:       incoming.Epoch,
				Heartbeat:   incoming.Heartbeat,
				LastUpdated: now,
			}
			fmt.Printf("\033[32m[DISCOVER] [%s] heartbeat: %d (epoch: %d)\033[0m\n",
				shortAddr(addr), incoming.Heartbeat, incoming.Epoch)
			continue
		}

		// Update state if we received a newer epoch/heartbeat version
		if isNewer(incoming.Epoch, incoming.Heartbeat, localPeer) {
			localPeer.Epoch = incoming.Epoch
			localPeer.Heartbeat = incoming.Heartbeat
			localPeer.LastUpdated = now
			if localPeer.IsFailed {
				localPeer.IsFailed = false
				fmt.Printf("\033[32m[RECOVER] [%s] heartbeat: %d\033[0m\n",
					shortAddr(addr), incoming.Heartbeat)
			} else {
				fmt.Printf("\033[90m[GOSSIP] [%s] heartbeat: %d\033[0m\n",
					shortAddr(addr), incoming.Heartbeat)
			}
		}
	}
}

func gossipLoop(conn *net.UDPConn) {
	ticker := time.NewTicker(gossipIter)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range ticker.C {
		peersMu.Lock()
		selfHeartbeat++
		peers[selfAddr].Heartbeat = selfHeartbeat
		peers[selfAddr].LastUpdated = time.Now()

		var fanoutSize = 1
	 	if(len(peers) > 3) {
	 		fanoutSize = 2
		}

		// Select random targets
		targets := getActivePeers(r, fanoutSize)

		// Prepare packet
		packet := GossipPacket{
			FromAddr: selfAddr,
			Peers:    make(map[string]PeerGossip),
		}
		for addr, state := range peers {
			packet.Peers[addr] = PeerGossip{
				Epoch:     state.Epoch,
				Heartbeat: state.Heartbeat,
			}
		}
		peersMu.Unlock()

		data, err := json.Marshal(packet)
		if err != nil {
			continue
		}

		for _, target := range targets {
			raddr, err := net.ResolveUDPAddr("udp", target)
			if err != nil {
				continue
			}
			_, _ = conn.WriteToUDP(data, raddr)
		}
	}
}

func getActivePeers(r *rand.Rand, k int) []string {
	var active []string
	for addr, state := range peers {
		if addr != selfAddr && !state.IsFailed {
			active = append(active, addr)
		}
	}

	if len(active) <= k {
		return active
	}

	r.Shuffle(len(active), func(i, j int) {
		active[i], active[j] = active[j], active[i]
	})
	return active[:k]
}

func failureDetectorLoop() {
	ticker := time.NewTicker(gossipIter)
	for range ticker.C {
		peersMu.Lock()
		now := time.Now()
		for addr, state := range peers {
			if addr == selfAddr || state.IsFailed {
				continue
			}

			if now.Sub(state.LastUpdated) > failTimeout {
				state.IsFailed = true
				fmt.Printf("\033[31m[DEAD] [%s] has failed (no update for %v)\033[0m\n", shortAddr(addr), failTimeout)
			}
		}
		peersMu.Unlock()
	}
}

// shortAddr returns the port part of "127.0.0.1:port".
func shortAddr(addr string) string {
	if idx := strings.Index(addr, ":"); idx != -1 {
		return addr[idx+1:]
	}
	return addr
}
