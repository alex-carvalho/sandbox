package main

import (
	"testing"
	"time"
)

func TestStateConvergence(t *testing.T) {
	// Create Node 1
	n1Config := NodeConfig{
		Addr:           "127.0.0.1:19001",
		GossipInterval: 100 * time.Millisecond,
		SuspectTimeout: 500 * time.Millisecond,
		DeadTimeout:    1000 * time.Millisecond,
		Fanout:         1,
	}
	n1, err := NewNode(n1Config, nil, nil)
	if err != nil {
		t.Fatalf("failed to create node 1: %v", err)
	}

	// Create Node 2
	n2Config := NodeConfig{
		Addr:           "127.0.0.1:19002",
		GossipInterval: 100 * time.Millisecond,
		SuspectTimeout: 500 * time.Millisecond,
		DeadTimeout:    1000 * time.Millisecond,
		Fanout:         1,
	}
	n2, err := NewNode(n2Config, nil, nil)
	if err != nil {
		t.Fatalf("failed to create node 2: %v", err)
	}

	// Start nodes
	if err := n1.Start(); err != nil {
		t.Fatalf("failed to start node 1: %v", err)
	}
	defer n1.Stop()

	if err := n2.Start(); err != nil {
		t.Fatalf("failed to start node 2: %v", err)
	}
	defer n2.Stop()

	// Link nodes by adding peer
	n1.AddPeer("127.0.0.1:19002")

	// Inject value on Node 1
	testKey := "cluster_greeting"
	testVal := "hello_from_node_1"
	n1.Set(testKey, testVal)

	// Poll Node 2 state until it converges or times out
	converged := false
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for !converged {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for state to converge on Node 2")
		case <-ticker.C:
			val, exists := n2.Get(testKey)
			if exists && val.Val == testVal {
				converged = true
			}
		}
	}
}

func TestNodeRevivalRecovery(t *testing.T) {
	// Create Node 1
	n1Config := NodeConfig{
		Addr:           "127.0.0.1:19011",
		GossipInterval: 50 * time.Millisecond,
		SuspectTimeout: 200 * time.Millisecond,
		DeadTimeout:    400 * time.Millisecond,
		Fanout:         1,
	}
	n1, err := NewNode(n1Config, nil, nil)
	if err != nil {
		t.Fatalf("failed to create node 1: %v", err)
	}

	// Create Node 2
	n2Config := NodeConfig{
		Addr:           "127.0.0.1:19012",
		GossipInterval: 50 * time.Millisecond,
		SuspectTimeout: 200 * time.Millisecond,
		DeadTimeout:    400 * time.Millisecond,
		Fanout:         1,
	}
	n2, err := NewNode(n2Config, nil, nil)
	if err != nil {
		t.Fatalf("failed to create node 2: %v", err)
	}

	// Start both nodes
	if err := n1.Start(); err != nil {
		t.Fatalf("failed to start node 1: %v", err)
	}
	defer n1.Stop()

	if err := n2.Start(); err != nil {
		t.Fatalf("failed to start node 2: %v", err)
	}

	n1.AddPeer("127.0.0.1:19012")

	// 1. Wait for Node 1 to see Node 2 as Alive
	time.Sleep(300 * time.Millisecond)

	// 2. Kill Node 2 (stop it)
	n2SnapshotState := n2.GetStateSnapshot()
	n2SnapshotMembers := n2.GetMembers()
	n2.Stop()

	// 3. Wait for Node 1 to declare Node 2 Suspected and then DEAD
	deadVerified := false
	timeout := time.After(1500 * time.Millisecond)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for !deadVerified {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for Node 1 to mark Node 2 as DEAD")
		case <-ticker.C:
			members := n1.GetMembers()
			for _, m := range members {
				if m.Addr == "127.0.0.1:19012" && m.Status == StatusDead {
					deadVerified = true
					break
				}
			}
		}
	}

	// 4. Revive Node 2 with its snapshot state
	n2Revived, err := NewNode(n2Config, n2SnapshotMembers, n2SnapshotState)
	if err != nil {
		t.Fatalf("failed to recreate node 2: %v", err)
	}

	if err := n2Revived.Start(); err != nil {
		t.Fatalf("failed to restart node 2: %v", err)
	}
	defer n2Revived.Stop()

	// 5. Wait for Node 1 to transition Node 2 back to ALIVE
	aliveVerified := false
	timeout = time.After(1500 * time.Millisecond)

	for !aliveVerified {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for Node 1 to mark revived Node 2 as ALIVE")
		case <-ticker.C:
			members := n1.GetMembers()
			for _, m := range members {
				if m.Addr == "127.0.0.1:19012" && m.Status == StatusAlive {
					aliveVerified = true
					break
				}
			}
		}
	}
}
