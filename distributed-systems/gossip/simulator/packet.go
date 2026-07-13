package main

// MemberStatus represents the health/reachability of a node.
type MemberStatus string

const (
	StatusAlive     MemberStatus = "ALIVE"
	StatusSuspected MemberStatus = "SUSPECTED"
	StatusDead      MemberStatus = "DEAD"
)

type Member struct {
	Addr           string       `json:"addr"`
	HeartbeatCount uint64       `json:"heartbeat_count"`
	LastUpdateNano int64        `json:"last_update_nano"` // Timestamp when the node last heard of this peer
	Status         MemberStatus `json:"status"`
}

// Value represents a replicated key-value state with metadata for conflict resolution (LWW).
type Value struct {
	Val       string `json:"val"`
	Version   uint64 `json:"version"`   // Increments on each update to this key
	Timestamp int64  `json:"timestamp"` // Nanoseconds timestamp of the update
}

type GossipMessage struct {
	FromAddr string           `json:"from_addr"`
	Members  []Member         `json:"members"`
	State    map[string]Value `json:"state"`
}
