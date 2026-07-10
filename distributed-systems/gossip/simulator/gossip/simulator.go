package gossip

import (
	"fmt"
	"sync"
	"time"
)

// KilledState preserves the state of a node across restarts.
type KilledState struct {
	State   map[string]Value
	Members []Member
}

// NodeStatus represents the current visual state of a node in the simulator.
type NodeStatus struct {
	Addr     string           `json:"addr"`
	IsKilled bool             `json:"is_killed"`
	Members  []Member         `json:"members"`
	State    map[string]Value `json:"state"`
}

// Simulator manages a collection of local gossip nodes.
type Simulator struct {
	mu           sync.RWMutex
	nodes        map[string]*Node
	killedStates map[string]*KilledState
	eventsChan   chan NodeEvent
	config       NodeConfig // Base config template
}

// NewSimulator creates a simulator instance.
func NewSimulator(gossipInterval, suspectTimeout, deadTimeout time.Duration, fanout int) *Simulator {
	sim := &Simulator{
		nodes:        make(map[string]*Node),
		killedStates: make(map[string]*KilledState),
		eventsChan:   make(chan NodeEvent, 1000),
	}

	sim.config = NodeConfig{
		GossipInterval: gossipInterval,
		SuspectTimeout: suspectTimeout,
		DeadTimeout:    deadTimeout,
		Fanout:         fanout,
		OnEvent: func(ev NodeEvent) {
			// Forward all node events to the simulator channel
			select {
			case sim.eventsChan <- ev:
			default:
				// Channel full, drop event to avoid blocking nodes
			}
		},
	}

	return sim
}

// Events returns the channel emitting node events.
func (s *Simulator) Events() <-chan NodeEvent {
	return s.eventsChan
}

// StartCluster initializes a cluster of N nodes.
func (s *Simulator) StartCluster(basePort int, size int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addresses := make([]string, size)
	for i := 0; i < size; i++ {
		addresses[i] = fmt.Sprintf("127.0.0.1:%d", basePort+i)
	}

	// Create nodes
	for _, addr := range addresses {
		node, err := NewNode(NodeConfig{
			Addr:           addr,
			GossipInterval: s.config.GossipInterval,
			SuspectTimeout: s.config.SuspectTimeout,
			DeadTimeout:    s.config.DeadTimeout,
			Fanout:         s.config.Fanout,
			OnEvent:        s.config.OnEvent,
		}, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to create node %s: %w", addr, err)
		}
		s.nodes[addr] = node
	}

	// Interconnect nodes initially by providing seed peers
	for addr, node := range s.nodes {
		for _, peerAddr := range addresses {
			if peerAddr != addr {
				node.AddPeer(peerAddr)
			}
		}
	}

	// Start all nodes
	for _, node := range s.nodes {
		if err := node.Start(); err != nil {
			return fmt.Errorf("failed to start node %s: %w", node.config.Addr, err)
		}
	}

	return nil
}

// StopCluster shuts down all active nodes.
func (s *Simulator) StopCluster() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, node := range s.nodes {
		node.Stop()
	}
	s.nodes = make(map[string]*Node)
}

// KillNode simulates a node crash/network disconnection.
func (s *Simulator) KillNode(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, exists := s.nodes[addr]
	if !exists {
		return fmt.Errorf("node %s not found or already killed", addr)
	}

	// Save state before stopping
	s.killedStates[addr] = &KilledState{
		State:   node.GetStateSnapshot(),
		Members: node.GetMembers(),
	}

	// Stop node loops and UDP socket
	node.Stop()
	delete(s.nodes, addr)

	// Emit simulation event
	s.config.OnEvent(NodeEvent{
		Node:      addr,
		Type:      "status_change",
		Detail:    "Node crashed / stopped",
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

// ReviveNode recovers a previously killed node.
func (s *Simulator) ReviveNode(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, active := s.nodes[addr]; active {
		return fmt.Errorf("node %s is already running", addr)
	}

	killed, exists := s.killedStates[addr]
	var state map[string]Value
	var members []Member

	if exists {
		state = killed.State
		members = killed.Members
		delete(s.killedStates, addr)
	}

	// Reset members to Suspected/Dead to force re-discovery, but keep themselves alive
	for i := range members {
		if members[i].Addr != addr {
			// Reset last update so they will immediately go through failure detection
			// unless we hear a heartbeat soon
			members[i].LastUpdateNano = time.Now().UnixNano()
		}
	}

	node, err := NewNode(NodeConfig{
		Addr:           addr,
		GossipInterval: s.config.GossipInterval,
		SuspectTimeout: s.config.SuspectTimeout,
		DeadTimeout:    s.config.DeadTimeout,
		Fanout:         s.config.Fanout,
		OnEvent:        s.config.OnEvent,
	}, members, state)
	if err != nil {
		return fmt.Errorf("failed to revive node %s: %w", addr, err)
	}

	if err := node.Start(); err != nil {
		return fmt.Errorf("failed to start revived node %s: %w", addr, err)
	}

	s.nodes[addr] = node

	s.config.OnEvent(NodeEvent{
		Node:      addr,
		Type:      "status_change",
		Detail:    "Node revived and restarted",
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

// AddNode dynamically spawns a new node and joins it to the cluster.
func (s *Simulator) AddNode(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodes[addr]; exists {
		return fmt.Errorf("node %s is already running", addr)
	}

	// Find an active bootstrap peer address
	var bootstrapPeer string
	for activeAddr := range s.nodes {
		bootstrapPeer = activeAddr
		break
	}

	node, err := NewNode(NodeConfig{
		Addr:           addr,
		GossipInterval: s.config.GossipInterval,
		SuspectTimeout: s.config.SuspectTimeout,
		DeadTimeout:    s.config.DeadTimeout,
		Fanout:         s.config.Fanout,
		OnEvent:        s.config.OnEvent,
	}, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create dynamic node %s: %w", addr, err)
	}

	if err := node.Start(); err != nil {
		return fmt.Errorf("failed to start dynamic node %s: %w", addr, err)
	}

	s.nodes[addr] = node

	// If there's a bootstrap peer, connect it
	if bootstrapPeer != "" {
		node.AddPeer(bootstrapPeer)
		s.config.OnEvent(NodeEvent{
			Node:      addr,
			Type:      "log",
			Detail:    fmt.Sprintf("Dynamic node bootstrapped with peer %s", bootstrapPeer),
			Timestamp: time.Now().UnixMilli(),
		})
	}

	return nil
}

// InjectValue sets a key-value pair on a specific node.
func (s *Simulator) InjectValue(addr, key, val string) error {
	s.mu.RLock()
	node, exists := s.nodes[addr]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("active node %s not found", addr)
	}

	node.Set(key, val)
	return nil
}

// GetStatusSnapshot collects the state of all nodes (active and dead).
func (s *Simulator) GetStatusSnapshot() []NodeStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statuses := make([]NodeStatus, 0, len(s.nodes)+len(s.killedStates))

	// Active nodes
	for addr, node := range s.nodes {
		statuses = append(statuses, NodeStatus{
			Addr:     addr,
			IsKilled: false,
			Members:  node.GetMembers(),
			State:    node.GetStateSnapshot(),
		})
	}

	// Killed nodes
	for addr, state := range s.killedStates {
		statuses = append(statuses, NodeStatus{
			Addr:     addr,
			IsKilled: true,
			Members:  state.Members,
			State:    state.State,
		})
	}

	return statuses
}
