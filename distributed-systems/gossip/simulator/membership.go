package main

import (
	"sync"
)

// MembershipTable holds the cluster membership state thread-safely.
type MembershipTable struct {
	mu      sync.RWMutex
	selfAddr string
	members map[string]*Member
}

// NewMembershipTable initializes a new table.
func NewMembershipTable(selfAddr string) *MembershipTable {
	table := &MembershipTable{
		selfAddr: selfAddr,
		members:  make(map[string]*Member),
	}
	// Add self to the table
	table.members[selfAddr] = &Member{
		Addr:           selfAddr,
		HeartbeatCount: 1,
		LastUpdateNano: 0, // Self never times out
		Status:         StatusAlive,
	}
	return table
}

// GetMembers returns a copy of all members in the table.
func (t *MembershipTable) GetMembers() []Member {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]Member, 0, len(t.members))
	for _, m := range t.members {
		result = append(result, *m)
	}
	return result
}

// GetActivePeers returns a list of addresses of alive or suspected peers (excluding self).
func (t *MembershipTable) GetActivePeers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var peers []string
	for addr, m := range t.members {
		if addr == t.selfAddr {
			continue
		}
		if m.Status == StatusAlive || m.Status == StatusSuspected {
			peers = append(peers, addr)
		}
	}
	return peers
}

// IncrementSelfHeartbeat increments the heartbeat of the local node.
func (t *MembershipTable) IncrementSelfHeartbeat() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	self, exists := t.members[t.selfAddr]
	if !exists {
		self = &Member{Addr: t.selfAddr, Status: StatusAlive}
		t.members[t.selfAddr] = self
	}
	self.HeartbeatCount++
	self.Status = StatusAlive
	self.LastUpdateNano = 0 // Self doesn't time out
	return self.HeartbeatCount
}

// SetHeartbeat directly overrides the heartbeat count of the local node (used during crash-recovery).
func (t *MembershipTable) SetHeartbeat(count uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	self, exists := t.members[t.selfAddr]
	if !exists {
		self = &Member{Addr: t.selfAddr, Status: StatusAlive}
		t.members[t.selfAddr] = self
	}
	self.HeartbeatCount = count
}

// Merge merges incoming membership information with the local table.
// Returns a list of changes that occurred (e.g. status changes).
type MemberChange struct {
	Addr   string
	Old    MemberStatus
	New    MemberStatus
	Reason string
}

func (t *MembershipTable) Merge(incoming []Member, nowNano int64) []MemberChange {
	t.mu.Lock()
	defer t.mu.Unlock()

	var changes []MemberChange

	self, selfExists := t.members[t.selfAddr]

	for _, incomingMem := range incoming {
		// Rule 1: If it's about us, we may need to refute claims
		if incomingMem.Addr == t.selfAddr {
			if selfExists {
				// If someone says we are SUSPECTED or DEAD, and they have an equal or higher heartbeat than ours,
				// we refute by incrementing our heartbeat and reasserting ALIVE.
				if (incomingMem.Status == StatusSuspected || incomingMem.Status == StatusDead) &&
					incomingMem.HeartbeatCount >= self.HeartbeatCount {
					self.HeartbeatCount = incomingMem.HeartbeatCount + 1
					self.Status = StatusAlive
					changes = append(changes, MemberChange{
						Addr:   t.selfAddr,
						Old:    incomingMem.Status,
						New:    StatusAlive,
						Reason: "refuting suspected/dead status claim",
					})
				}
			}
			continue
		}

		localMem, exists := t.members[incomingMem.Addr]
		if !exists {
			// New node discovered. We accept it as long as it's not DEAD.
			// (If it's already dead and we don't know it, we don't need to add it to save space).
			if incomingMem.Status != StatusDead {
				newMem := &Member{
					Addr:           incomingMem.Addr,
					HeartbeatCount: incomingMem.HeartbeatCount,
					LastUpdateNano: nowNano,
					Status:         incomingMem.Status,
				}
				t.members[incomingMem.Addr] = newMem
				changes = append(changes, MemberChange{
					Addr:   incomingMem.Addr,
					Old:    "",
					New:    incomingMem.Status,
					Reason: "discovered new node",
				})
			}
			continue
		}

		// Rule 2: Newer heartbeat overrides everything
		if incomingMem.HeartbeatCount > localMem.HeartbeatCount {
			oldStatus := localMem.Status
			localMem.HeartbeatCount = incomingMem.HeartbeatCount
			localMem.Status = incomingMem.Status
			localMem.LastUpdateNano = nowNano

			if oldStatus != incomingMem.Status {
				changes = append(changes, MemberChange{
					Addr:   incomingMem.Addr,
					Old:    oldStatus,
					New:    incomingMem.Status,
					Reason: "newer heartbeat received",
				})
			}
		} else if incomingMem.HeartbeatCount == localMem.HeartbeatCount {
			// Rule 3: Equal heartbeat, but status might be worse (Alive -> Suspected -> Dead)
			if localMem.Status == StatusAlive && incomingMem.Status == StatusSuspected {
				localMem.Status = StatusSuspected
				localMem.LastUpdateNano = nowNano
				changes = append(changes, MemberChange{
					Addr:   incomingMem.Addr,
					Old:    StatusAlive,
					New:    StatusSuspected,
					Reason: "peer marked as suspected at same heartbeat version",
				})
			} else if (localMem.Status == StatusAlive || localMem.Status == StatusSuspected) && incomingMem.Status == StatusDead {
				oldStatus := localMem.Status
				localMem.Status = StatusDead
				localMem.LastUpdateNano = nowNano
				changes = append(changes, MemberChange{
					Addr:   incomingMem.Addr,
					Old:    oldStatus,
					New:    StatusDead,
					Reason: "peer marked as dead at same heartbeat version",
				})
			}
		}
	}

	return changes
}

// DetectFailures transitions nodes to Suspected or Dead if heartbeats stall.
// timeoutSuspectNs: how long without update before suspecting
// timeoutDeadNs: how long suspected before marking dead
func (t *MembershipTable) DetectFailures(timeoutSuspectNs, timeoutDeadNs int64, nowNano int64) []MemberChange {
	t.mu.Lock()
	defer t.mu.Unlock()

	var changes []MemberChange

	for addr, m := range t.members {
		if addr == t.selfAddr {
			continue
		}

		elapsed := nowNano - m.LastUpdateNano

		if m.Status == StatusAlive && elapsed > timeoutSuspectNs {
			m.Status = StatusSuspected
			m.LastUpdateNano = nowNano // Reset timestamp to measure dead timeout from here
			changes = append(changes, MemberChange{
				Addr:   addr,
				Old:    StatusAlive,
				New:    StatusSuspected,
				Reason: "no heartbeat received within suspect timeout",
			})
		} else if m.Status == StatusSuspected && elapsed > timeoutDeadNs {
			m.Status = StatusDead
			m.LastUpdateNano = nowNano
			changes = append(changes, MemberChange{
				Addr:   addr,
				Old:    StatusSuspected,
				New:    StatusDead,
				Reason: "suspected node did not recover within dead timeout",
			})
		}
	}

	return changes
}
