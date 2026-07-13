package main

import (
	"fmt"
	"sync"
	"time"
)

type KilledState struct {
	State   map[string]Value
	Members []Member
}

type NodeStatus struct {
	Addr     string           `json:"addr"`
	IsKilled bool             `json:"is_killed"`
	Members  []Member         `json:"members"`
	State    map[string]Value `json:"state"`
}

type Simulator struct {
	mu           sync.RWMutex
	nodes        map[string]*Node
	killedStates map[string]*KilledState
	eventsChan   chan NodeEvent
	config       NodeConfig // Base config template
}

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
			select {
			case sim.eventsChan <- ev:
			default:
				// Channel full, drop event to avoid blocking nodes
			}
		},
	}

	return sim
}

func (s *Simulator) Events() <-chan NodeEvent {
	return s.eventsChan
}

func (s *Simulator) StartCluster(basePort int, size int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addresses := make([]string, size)
	for i := 0; i < size; i++ {
		addresses[i] = fmt.Sprintf("127.0.0.1:%d", basePort+i)
	}

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

	for addr, node := range s.nodes {
		for _, peerAddr := range addresses {
			if peerAddr != addr {
				node.AddPeer(peerAddr)
			}
		}
	}

	for _, node := range s.nodes {
		if err := node.Start(); err != nil {
			return fmt.Errorf("failed to start node %s: %w", node.config.Addr, err)
		}
	}

	return nil
}

func (s *Simulator) StopCluster() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, node := range s.nodes {
		node.Stop()
	}
	s.nodes = make(map[string]*Node)
}

func (s *Simulator) KillNode(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, exists := s.nodes[addr]
	if !exists {
		return fmt.Errorf("node %s not found or already killed", addr)
	}

	s.killedStates[addr] = &KilledState{
		State:   node.GetStateSnapshot(),
		Members: node.GetMembers(),
	}

	node.Stop()
	delete(s.nodes, addr)

	s.config.OnEvent(NodeEvent{
		Node:      addr,
		Type:      "status_change",
		Detail:    "Node crashed / stopped",
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

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

func (s *Simulator) AddNode(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodes[addr]; exists {
		return fmt.Errorf("node %s is already running", addr)
	}

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

func (s *Simulator) GetStatusSnapshot() []NodeStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statuses := make([]NodeStatus, 0, len(s.nodes)+len(s.killedStates))

	for addr, node := range s.nodes {
		statuses = append(statuses, NodeStatus{
			Addr:     addr,
			IsKilled: false,
			Members:  node.GetMembers(),
			State:    node.GetStateSnapshot(),
		})
	}

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
