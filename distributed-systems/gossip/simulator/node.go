package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// NodeEvent represents an event happening inside a gossip node, used for visualization.
type NodeEvent struct {
	Timestamp int64  `json:"timestamp"`
	Node      string `json:"node"`
	Type      string `json:"type"` // "sent", "received", "status_change", "state_update", "log"
	Detail    string `json:"detail"`
	Target    string `json:"target,omitempty"`
	Key       string `json:"key,omitempty"`
	Val       string `json:"val,omitempty"`
	Version   uint64 `json:"version,omitempty"`
}

// NodeConfig details parameters for a gossip node.
type NodeConfig struct {
	Addr           string
	GossipInterval time.Duration
	SuspectTimeout time.Duration
	DeadTimeout    time.Duration
	Fanout         int
	OnEvent        func(NodeEvent)
}

// Node represents a running gossip node.
type Node struct {
	config    NodeConfig
	members   *MembershipTable
	stateMu   sync.RWMutex
	state     map[string]Value
	conn      *net.UDPConn
	stopChan  chan struct{}
	wg        sync.WaitGroup
	isKilled  bool
	killMu    sync.RWMutex
}

// NewNode initializes a new gossip node. It can accept initial membership and state to simulate recovery.
func NewNode(config NodeConfig, initialMembers []Member, initialState map[string]Value) (*Node, error) {
	node := &Node{
		config:   config,
		members:  NewMembershipTable(config.Addr),
		state:    make(map[string]Value),
		stopChan: make(chan struct{}),
	}

	// Populate initial state if any
	if initialState != nil {
		for k, v := range initialState {
			node.state[k] = v
		}
	}

	// Populate initial membership if any
	if len(initialMembers) > 0 {
		// Restore self-heartbeat count from the pre-crash snapshot.
		// This prevents heartbeat count epoch regression, ensuring other nodes in the cluster
		// accept our new gossip updates rather than dropping them as outdated.
		for _, m := range initialMembers {
			if m.Addr == config.Addr {
				node.members.SetHeartbeat(m.HeartbeatCount)
				break
			}
		}

		node.members.Merge(initialMembers, time.Now().UnixNano())
	}

	return node, nil
}

// Start binds to the UDP port and starts the node loops.
func (n *Node) Start() error {
	addr, err := net.ResolveUDPAddr("udp", n.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP address: %w", err)
	}
	n.conn = conn

	n.emit(NodeEvent{
		Type:   "log",
		Detail: fmt.Sprintf("Node started, listening on UDP %s", n.config.Addr),
	})

	n.wg.Add(3)
	go n.listen()
	go n.gossipLoop()
	go n.failureDetectorLoop()

	return nil
}

// Stop closes connections and stops goroutines.
func (n *Node) Stop() {
	n.killMu.Lock()
	if n.isKilled {
		n.killMu.Unlock()
		return
	}
	n.isKilled = true
	n.killMu.Unlock()

	close(n.stopChan)
	if n.conn != nil {
		n.conn.Close()
	}
	n.wg.Wait()
	n.emit(NodeEvent{
		Type:   "log",
		Detail: "Node stopped",
	})
}

// Set updates or inserts a key-value pair in the local node state.
func (n *Node) Set(key, val string) {
	n.stateMu.Lock()
	v, exists := n.state[key]
	newVer := uint64(1)
	if exists {
		newVer = v.Version + 1
	}
	n.state[key] = Value{
		Val:       val,
		Version:   newVer,
		Timestamp: time.Now().UnixNano(),
	}
	n.stateMu.Unlock()

	n.emit(NodeEvent{
		Type:    "state_update",
		Detail:  fmt.Sprintf("State modified: %s = %s (v%d)", key, val, newVer),
		Key:     key,
		Val:     val,
		Version: newVer,
	})
}

// Get retrieves a key value from the local node state.
func (n *Node) Get(key string) (Value, bool) {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	val, exists := n.state[key]
	return val, exists
}

// GetStateSnapshot returns a copy of the entire node's state.
func (n *Node) GetStateSnapshot() map[string]Value {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.copyState()
}

// GetMembers returns the node's current cluster view.
func (n *Node) GetMembers() []Member {
	return n.members.GetMembers()
}

// AddPeer manually adds a peer's address to start gossiping with it.
func (n *Node) AddPeer(peerAddr string) {
	n.members.Merge([]Member{
		{
			Addr:           peerAddr,
			HeartbeatCount: 0,
			LastUpdateNano: time.Now().UnixNano(),
			Status:         StatusAlive,
		},
	}, time.Now().UnixNano())

	n.emit(NodeEvent{
		Type:   "status_change",
		Detail: fmt.Sprintf("Added peer: %s", peerAddr),
		Target: peerAddr,
	})
}

// Helpers

func (n *Node) emit(ev NodeEvent) {
	ev.Timestamp = time.Now().UnixMilli()
	ev.Node = n.config.Addr
	if n.config.OnEvent != nil {
		n.config.OnEvent(ev)
	}
}

func (n *Node) copyState() map[string]Value {
	m := make(map[string]Value, len(n.state))
	for k, v := range n.state {
		m[k] = v
	}
	return m
}

func (n *Node) isClosed() bool {
	n.killMu.RLock()
	defer n.killMu.RUnlock()
	return n.isKilled
}

func (n *Node) mergeState(incoming map[string]Value) {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()

	for k, incomingVal := range incoming {
		localVal, exists := n.state[k]
		if !exists || incomingVal.Version > localVal.Version ||
			(incomingVal.Version == localVal.Version && incomingVal.Timestamp > localVal.Timestamp) {
			n.state[k] = incomingVal
			n.emit(NodeEvent{
				Type:    "state_update",
				Detail:  fmt.Sprintf("State converged: %s = %s (v%d)", k, incomingVal.Val, incomingVal.Version),
				Key:     k,
				Val:     incomingVal.Val,
				Version: incomingVal.Version,
			})
		}
	}
}

func (n *Node) listen() {
	defer n.wg.Done()
	buf := make([]byte, 65535)

	for {
		select {
		case <-n.stopChan:
			return
		default:
			// Read with read-deadline to avoid locking forever on Shutdown
			_ = n.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			size, _, err := n.conn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if n.isClosed() {
					return
				}
				continue
			}

			if n.isClosed() {
				return
			}

			var msg GossipMessage
			if err := json.Unmarshal(buf[:size], &msg); err != nil {
				continue
			}

			n.emit(NodeEvent{
				Type:   "received",
				Detail: fmt.Sprintf("Received gossip from %s", msg.FromAddr),
				Target: msg.FromAddr,
			})

			// Merge membership list
			changes := n.members.Merge(msg.Members, time.Now().UnixNano())
			for _, c := range changes {
				n.emit(NodeEvent{
					Type:   "status_change",
					Detail: fmt.Sprintf("%s state changed to %s (%s)", c.Addr, c.New, c.Reason),
					Target: c.Addr,
				})
			}

			// Merge key-value state
			n.mergeState(msg.State)
		}
	}
}

func (n *Node) gossipLoop() {
	defer n.wg.Done()
	ticker := time.NewTicker(n.config.GossipInterval)
	defer ticker.Stop()

	// Seed random number generator
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	for {
		select {
		case <-n.stopChan:
			return
		case <-ticker.C:
			if n.isClosed() {
				return
			}

			// Increment own heartbeat
			n.members.IncrementSelfHeartbeat()

			// Get potential communication peers
			peers := n.members.GetActivePeers()
			if len(peers) == 0 {
				continue
			}

			// Select fanout random peers
			selected := n.selectRandomPeers(peers, n.config.Fanout, r)

			// Prepare packet payload
			n.stateMu.RLock()
			msg := GossipMessage{
				FromAddr: n.config.Addr,
				Members:  n.members.GetMembers(),
				State:    n.copyState(),
			}
			n.stateMu.RUnlock()

			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			// Send gossip packets
			for _, peerAddrStr := range selected {
				peerAddr, err := net.ResolveUDPAddr("udp", peerAddrStr)
				if err != nil {
					continue
				}

				if n.isClosed() {
					return
				}

				_, _ = n.conn.WriteToUDP(data, peerAddr)
				n.emit(NodeEvent{
					Type:   "sent",
					Detail: fmt.Sprintf("Sent gossip to %s", peerAddrStr),
					Target: peerAddrStr,
				})
			}
		}
	}
}

func (n *Node) failureDetectorLoop() {
	defer n.wg.Done()
	ticker := time.NewTicker(n.config.GossipInterval)
	defer ticker.Stop()

	for {
		select {
		case <-n.stopChan:
			return
		case <-ticker.C:
			if n.isClosed() {
				return
			}

			changes := n.members.DetectFailures(
				n.config.SuspectTimeout.Nanoseconds(),
				n.config.DeadTimeout.Nanoseconds(),
				time.Now().UnixNano(),
			)

			for _, c := range changes {
				n.emit(NodeEvent{
					Type:   "status_change",
					Detail: fmt.Sprintf("Failure detector: %s is now %s (%s)", c.Addr, c.New, c.Reason),
					Target: c.Addr,
				})
			}
		}
	}
}

func (n *Node) selectRandomPeers(peers []string, k int, r *rand.Rand) []string {
	if len(peers) <= k {
		return peers
	}
	shuffled := make([]string, len(peers))
	copy(shuffled, peers)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:k]
}
