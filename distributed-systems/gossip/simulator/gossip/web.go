package gossip

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

//go:embed index.html
var indexHTML string

// WebServer handles HTTP API and SSE events.
type WebServer struct {
	sim      *Simulator
	node     *Node // Optional: used when visualizing a standalone single node
	port     int
	clients  map[chan NodeEvent]bool
	clientMu sync.Mutex
}

// NewWebServer initializes the HTTP server wrapper.
func NewWebServer(sim *Simulator, port int) *WebServer {
	return &WebServer{
		sim:     sim,
		port:    port,
		clients: make(map[chan NodeEvent]bool),
	}
}

// SetNode binds a single standalone node to this web visualizer.
func (ws *WebServer) SetNode(node *Node) {
	ws.node = node
}

// Start launches the broadcast loop and listens for HTTP requests.
func (ws *WebServer) Start() error {
	ws.startBroadcast()

	mux := http.NewServeMux()

	// Static routes
	mux.HandleFunc("/", ws.handleIndex)

	// API routes
	mux.HandleFunc("/api/events", ws.handleSSE)
	mux.HandleFunc("/api/status", ws.handleStatus)
	mux.HandleFunc("/api/kill", ws.handleKill)
	mux.HandleFunc("/api/revive", ws.handleRevive)
	mux.HandleFunc("/api/add", ws.handleAdd)
	mux.HandleFunc("/api/inject", ws.handleInject)
	mux.HandleFunc("/api/report-event", ws.handleReportEvent)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", ws.port),
		Handler: mux,
	}

	return server.ListenAndServe()
}

func (ws *WebServer) startBroadcast() {
	if ws.sim == nil {
		return // In standalone node mode, events are fed directly via BroadcastEvent
	}
	go func() {
		for ev := range ws.sim.Events() {
			ws.BroadcastEvent(ev)
		}
	}()
}

// BroadcastEvent relays a node event to all connected SSE browser clients.
func (ws *WebServer) BroadcastEvent(ev NodeEvent) {
	ws.clientMu.Lock()
	for ch := range ws.clients {
		select {
		case ch <- ev:
		default:
			// Drop event if client queue is full
		}
	}
	ws.clientMu.Unlock()
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func (ws *WebServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := make(chan NodeEvent, 200)

	ws.clientMu.Lock()
	ws.clients[ch] = true
	ws.clientMu.Unlock()

	defer func() {
		ws.clientMu.Lock()
		delete(ws.clients, ch)
		ws.clientMu.Unlock()
		close(ch)
	}()

	// Flush opening event
	_, _ = fmt.Fprintf(w, "event: open\ndata: connected\n\n")
	flusher.Flush()

	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (ws *WebServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// If we are visualizing a single standalone node, compile status from its point of view
	if ws.node != nil {
		members := ws.node.GetMembers()
		selfState := ws.node.GetStateSnapshot()

		statuses := make([]NodeStatus, 0, len(members))
		for _, m := range members {
			isSelf := m.Addr == ws.node.config.Addr
			var state map[string]Value
			if isSelf {
				state = selfState
			} else {
				// Show our replica's view of the database
				state = selfState
			}

			isKilled := m.Status == StatusDead

			statuses = append(statuses, NodeStatus{
				Addr:     m.Addr,
				IsKilled: isKilled,
				Members:  members,
				State:    state,
			})
		}
		_ = json.NewEncoder(w).Encode(statuses)
		return
	}

	snapshot := ws.sim.GetStatusSnapshot()
	_ = json.NewEncoder(w).Encode(snapshot)
}

func (ws *WebServer) handleKill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ws.node != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Kill node action is not supported in standalone mode. Kill the process manually in your terminal!"})
		return
	}

	var req struct {
		Node string `json:"node"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ws.sim.KillNode(req.Node); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ws *WebServer) handleRevive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ws.node != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Revive node action is not supported in standalone mode. Restart the process manually in your terminal!"})
		return
	}

	var req struct {
		Node string `json:"node"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ws.sim.ReviveNode(req.Node); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ws *WebServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Node string `json:"node"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ws.node != nil {
		ws.node.AddPeer(req.Node)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	if err := ws.sim.AddNode(req.Node); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ws *WebServer) handleInject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Node  string `json:"node"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ws.node != nil {
		ws.node.Set(req.Key, req.Value)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	if err := ws.sim.InjectValue(req.Node, req.Key, req.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleReportEvent processes telemetry events posted by remote standalone nodes.
func (ws *WebServer) handleReportEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ev NodeEvent
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ws.BroadcastEvent(ev)
	w.WriteHeader(http.StatusOK)
}
