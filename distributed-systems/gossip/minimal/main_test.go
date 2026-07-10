package main

import (
	"testing"
	"time"
)

func TestDiscoveryAndMerge(t *testing.T) {
	// Reset global state for testing
	peersMu.Lock()
	peers = make(map[string]*PeerState)
	selfAddr = "127.0.0.1:9001"
	selfEpoch = 1000
	selfHeartbeat = 1
	peers[selfAddr] = &PeerState{Epoch: selfEpoch, Heartbeat: selfHeartbeat, LastUpdated: time.Now()}
	peersMu.Unlock()

	peerAddr := "127.0.0.1:9002"

	// 1. Merge new peer gossip
	incoming := map[string]PeerGossip{
		peerAddr: {Epoch: 500, Heartbeat: 10},
	}
	mergeGossip("sender", incoming)

	peersMu.RLock()
	pState, exists := peers[peerAddr]
	peersMu.RUnlock()

	if !exists {
		t.Fatalf("expected peer %s to be discovered and created", peerAddr)
	}
	if pState.Epoch != 500 || pState.Heartbeat != 10 {
		t.Errorf("expected version (epoch: 500, hb: 10), got (epoch: %d, hb: %d)", pState.Epoch, pState.Heartbeat)
	}

	// 2. Merge advanced heartbeat
	incoming = map[string]PeerGossip{
		peerAddr: {Epoch: 500, Heartbeat: 11},
	}
	mergeGossip("sender", incoming)

	if pState.Heartbeat != 11 {
		t.Errorf("expected heartbeat to advance to 11, got %d", pState.Heartbeat)
	}
}

func TestEpochRecovery(t *testing.T) {
	// Reset global state
	peersMu.Lock()
	peers = make(map[string]*PeerState)
	selfAddr = "127.0.0.1:9001"
	selfEpoch = 1000
	selfHeartbeat = 1
	peers[selfAddr] = &PeerState{Epoch: selfEpoch, Heartbeat: selfHeartbeat, LastUpdated: time.Now()}
	peersMu.Unlock()

	peerAddr := "127.0.0.1:9002"

	// Add peer with epoch 500 and heartbeat 15, marked DEAD
	peersMu.Lock()
	peers[peerAddr] = &PeerState{
		Epoch:       500,
		Heartbeat:   15,
		LastUpdated: time.Now().Add(-10 * time.Second),
		IsFailed:    true,
	}
	peersMu.Unlock()

	// Merge gossip from peer restarting with a new epoch (600) but heartbeat reset to 1
	incoming := map[string]PeerGossip{
		peerAddr: {Epoch: 600, Heartbeat: 1},
	}
	mergeGossip("sender", incoming)

	peersMu.RLock()
	pState := peers[peerAddr]
	peersMu.RUnlock()

	if pState.Epoch != 600 || pState.Heartbeat != 1 {
		t.Errorf("expected peer to be updated to (epoch: 600, hb: 1), got (epoch: %d, hb: %d)", pState.Epoch, pState.Heartbeat)
	}
	if pState.IsFailed {
		t.Error("expected peer to be recovered (IsFailed = false)")
	}
}

func TestFailureDetector(t *testing.T) {
	peersMu.Lock()
	peers = make(map[string]*PeerState)
	selfAddr = "127.0.0.1:9001"
	selfEpoch = 1000
	selfHeartbeat = 1
	peers[selfAddr] = &PeerState{Epoch: selfEpoch, Heartbeat: selfHeartbeat, LastUpdated: time.Now()}

	peerAddr := "127.0.0.1:9002"
	failTimeout = 50 * time.Millisecond

	// Peer 9002 was updated 100ms ago (stale)
	peers[peerAddr] = &PeerState{
		Epoch:       500,
		Heartbeat:   10,
		LastUpdated: time.Now().Add(-100 * time.Millisecond),
		IsFailed:    false,
	}
	peersMu.Unlock()

	// Execute failure detection check
	peersMu.Lock()
	now := time.Now()
	for addr, state := range peers {
		if addr == selfAddr || state.IsFailed {
			continue
		}
		if now.Sub(state.LastUpdated) > failTimeout {
			state.IsFailed = true
		}
	}
	peersMu.Unlock()

	peersMu.RLock()
	pState := peers[peerAddr]
	peersMu.RUnlock()

	if !pState.IsFailed {
		t.Error("expected peer to be flagged as failed due to heartbeat stall")
	}
}
