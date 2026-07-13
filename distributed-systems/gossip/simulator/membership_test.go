package main

import (
	"testing"
	"time"
)

func TestNewMembershipTable(t *testing.T) {
	selfAddr := "127.0.0.1:9001"
	table := NewMembershipTable(selfAddr)

	members := table.GetMembers()
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}

	if members[0].Addr != selfAddr {
		t.Errorf("expected self address %s, got %s", selfAddr, members[0].Addr)
	}

	if members[0].Status != StatusAlive {
		t.Errorf("expected self status %s, got %s", StatusAlive, members[0].Status)
	}
}

func TestMergeMembership(t *testing.T) {
	selfAddr := "127.0.0.1:9001"
	table := NewMembershipTable(selfAddr)
	now := time.Now().UnixNano()

	// 1. Discover new node
	peerAddr := "127.0.0.1:9002"
	incoming := []Member{
		{
			Addr:           peerAddr,
			HeartbeatCount: 10,
			Status:         StatusAlive,
		},
	}

	changes := table.Merge(incoming, now)
	if len(changes) != 1 {
		t.Fatalf("expected 1 status change, got %d", len(changes))
	}
	if changes[0].Addr != peerAddr || changes[0].New != StatusAlive {
		t.Errorf("expected change on %s to ALIVE, got %+v", peerAddr, changes[0])
	}

	// Verify heartbeat updated
	members := table.GetMembers()
	var found bool
	for _, m := range members {
		if m.Addr == peerAddr {
			found = true
			if m.HeartbeatCount != 10 {
				t.Errorf("expected heartbeat 10, got %d", m.HeartbeatCount)
			}
		}
	}
	if !found {
		t.Fatalf("peer %s not found in table", peerAddr)
	}

	// 2. Same heartbeat, suspect transition
	incoming = []Member{
		{
			Addr:           peerAddr,
			HeartbeatCount: 10,
			Status:         StatusSuspected,
		},
	}
	changes = table.Merge(incoming, now)
	if len(changes) != 1 {
		t.Fatalf("expected 1 status change, got %d", len(changes))
	}
	if changes[0].New != StatusSuspected {
		t.Errorf("expected status to change to SUSPECTED, got %s", changes[0].New)
	}

	// 3. Higher heartbeat, alive transition
	incoming = []Member{
		{
			Addr:           peerAddr,
			HeartbeatCount: 11,
			Status:         StatusAlive,
		},
	}
	changes = table.Merge(incoming, now)
	if len(changes) != 1 {
		t.Fatalf("expected 1 status change, got %d", len(changes))
	}
	if changes[0].New != StatusAlive {
		t.Errorf("expected status to change to ALIVE, got %s", changes[0].New)
	}
}

func TestRefutation(t *testing.T) {
	selfAddr := "127.0.0.1:9001"
	table := NewMembershipTable(selfAddr)

	// Increment self heartbeat to 2
	localHb := table.IncrementSelfHeartbeat()

	// Another node says we are SUSPECTED with our current heartbeat
	incoming := []Member{
		{
			Addr:           selfAddr,
			HeartbeatCount: localHb,
			Status:         StatusSuspected,
		},
	}

	changes := table.Merge(incoming, time.Now().UnixNano())
	if len(changes) != 1 {
		t.Fatalf("expected 1 change (refutation), got %d", len(changes))
	}

	if changes[0].New != StatusAlive {
		t.Errorf("expected status to remain ALIVE, got %s", changes[0].New)
	}

	members := table.GetMembers()
	for _, m := range members {
		if m.Addr == selfAddr {
			if m.HeartbeatCount != localHb+1 {
				t.Errorf("expected heartbeat to increment to %d to refute, got %d", localHb+1, m.HeartbeatCount)
			}
			if m.Status != StatusAlive {
				t.Errorf("expected self status to remain ALIVE, got %s", m.Status)
			}
		}
	}
}

func TestDetectFailures(t *testing.T) {
	selfAddr := "127.0.0.1:9001"
	table := NewMembershipTable(selfAddr)
	peerAddr := "127.0.0.1:9002"

	now := time.Now().UnixNano()
	table.Merge([]Member{
		{
			Addr:           peerAddr,
			HeartbeatCount: 1,
			Status:         StatusAlive,
		},
	}, now)

	suspectTimeout := 500 * time.Millisecond
	deadTimeout := 1000 * time.Millisecond

	// 1. Check before suspectTimeout expires -> no changes
	changes := table.DetectFailures(
		suspectTimeout.Nanoseconds(),
		deadTimeout.Nanoseconds(),
		now+int64(200*time.Millisecond),
	)
	if len(changes) != 0 {
		t.Errorf("expected no failures detected, got %d changes", len(changes))
	}

	// 2. Check after suspectTimeout expires -> transitions to SUSPECTED
	changes = table.DetectFailures(
		suspectTimeout.Nanoseconds(),
		deadTimeout.Nanoseconds(),
		now+int64(600*time.Millisecond),
	)
	if len(changes) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(changes))
	}
	if changes[0].Addr != peerAddr || changes[0].New != StatusSuspected {
		t.Errorf("expected peer to be SUSPECTED, got %+v", changes[0])
	}

	// The DetectFailures sets lastUpdate to the detection time (600ms)
	nowSuspected := now + int64(600*time.Millisecond)

	// 3. Check before deadTimeout from the suspicion time -> no changes
	changes = table.DetectFailures(
		suspectTimeout.Nanoseconds(),
		deadTimeout.Nanoseconds(),
		nowSuspected+int64(500*time.Millisecond),
	)
	if len(changes) != 0 {
		t.Errorf("expected no dead transitions, got %d changes", len(changes))
	}

	// 4. Check after deadTimeout expires -> transitions to DEAD
	changes = table.DetectFailures(
		suspectTimeout.Nanoseconds(),
		deadTimeout.Nanoseconds(),
		nowSuspected+int64(1100*time.Millisecond),
	)
	if len(changes) != 1 {
		t.Fatalf("expected 1 dead transition, got %d", len(changes))
	}
	if changes[0].Addr != peerAddr || changes[0].New != StatusDead {
		t.Errorf("expected peer to be DEAD, got %+v", changes[0])
	}
}
