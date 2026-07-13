package main

// MemberStatus represents the health/reachability of a node.
type MemberStatus string

const (
	StatusAlive     MemberStatus = "ALIVE"
	StatusSuspected MemberStatus = "SUSPECTED"
	StatusDead      MemberStatus = "DEAD"
)

// Member represents the metadata of a node in the cluster.
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

// GossipMessage is the packet serialized and sent over UDP.
type GossipMessage struct {
	FromAddr string           `json:"from_addr"`
	Members  []Member         `json:"members"`
	State    map[string]Value `json:"state"`
}
