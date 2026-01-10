package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shibaleo/go-mcp-dev/internal/observability"
)

type Handler struct {
	tools    map[string]ToolExecutor
	sessions map[string]*Session
	mu       sync.RWMutex
}

type Session struct {
	id       string
	writer   http.ResponseWriter
	flusher  http.Flusher
	done     chan struct{}
	messages chan []byte
}

func NewHandler() *Handler {
	return &Handler{
		tools:    make(map[string]ToolExecutor),
		sessions: make(map[string]*Session),
	}
}

func (h *Handler) RegisterTool(tool ToolExecutor) {
	def := tool.Definition()
	h.tools[def.Name] = tool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleSSE(w, r)
	case http.MethodPost:
		h.handleMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create session
	sessionID := fmt.Sprintf("%d", r.Context().Value("session_id"))
	if sessionID == "<nil>" {
		sessionID = fmt.Sprintf("%p", r)
	}

	session := &Session{
		id:       sessionID,
		writer:   w,
		flusher:  flusher,
		done:     make(chan struct{}),
		messages: make(chan []byte, 100),
	}

	h.mu.Lock()
	h.sessions[sessionID] = session
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.sessions, sessionID)
		h.mu.Unlock()
		close(session.done)
	}()

	// Send endpoint event (MCP SSE protocol)
	// The endpoint tells the client where to POST messages
	fmt.Fprintf(w, "event: endpoint\ndata: /mcp?sessionId=%s\n\n", sessionID)
	flusher.Flush()
	log.Printf("SSE connection established, session=%s", sessionID)

	// Keep connection open and send messages
	for {
		select {
		case msg := <-session.messages:
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			log.Printf("SSE connection closed, session=%s", sessionID)
			return
		}
	}
}

func (h *Handler) handleMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		// For simple single-session mode, use first available session or create inline response
		h.handleInlineMessage(w, r)
		return
	}

	h.mu.RLock()
	session, ok := h.sessions[sessionID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		h.sendToSession(session, nil, &Error{Code: ParseError, Message: "Parse error"})
		w.WriteHeader(http.StatusAccepted)
		return
	}

	log.Printf("Received request: method=%s id=%v session=%s", req.Method, req.ID, sessionID)

	result, rpcErr := h.processRequest(&req)
	if rpcErr != nil {
		h.sendToSession(session, req.ID, rpcErr)
	} else if req.ID != nil {
		h.sendResultToSession(session, req.ID, result)
	}

	w.WriteHeader(http.StatusAccepted)
}

// handleInlineMessage handles POST requests without SSE (for simple testing)
func (h *Handler) handleInlineMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{JSONRPC: "2.0", Error: &Error{Code: ParseError, Message: "Parse error"}}
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("Received inline request: method=%s id=%v", req.Method, req.ID)

	result, rpcErr := h.processRequest(&req)

	w.Header().Set("Content-Type", "application/json")
	var resp Response
	if rpcErr != nil {
		resp = Response{JSONRPC: "2.0", ID: req.ID, Error: rpcErr}
	} else {
		resp = Response{JSONRPC: "2.0", ID: req.ID, Result: result}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) sendToSession(session *Session, id interface{}, err *Error) {
	resp := Response{JSONRPC: "2.0", ID: id, Error: err}
	data, _ := json.Marshal(resp)
	select {
	case session.messages <- data:
	default:
		log.Printf("Session message buffer full")
	}
}

func (h *Handler) sendResultToSession(session *Session, id interface{}, result interface{}) {
	resp := Response{JSONRPC: "2.0", ID: id, Result: result}
	data, _ := json.Marshal(resp)
	select {
	case session.messages <- data:
	default:
		log.Printf("Session message buffer full")
	}
}

func (h *Handler) processRequest(req *Request) (interface{}, *Error) {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(req), nil
	case "initialized":
		return nil, nil
	case "tools/list":
		return h.handleToolsList(), nil
	case "tools/call":
		return h.handleToolCall(req)
	default:
		return nil, &Error{Code: MethodNotFound, Message: "Method not found"}
	}
}

func (h *Handler) handleInitialize(req *Request) *InitializeResult {
	return &InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "go-mcp-dev",
			Version: "0.1.0",
		},
	}
}

func (h *Handler) handleToolsList() *ToolsListResult {
	tools := make([]Tool, 0, len(h.tools))
	for _, executor := range h.tools {
		tools = append(tools, executor.Definition())
	}
	return &ToolsListResult{Tools: tools}
}

func (h *Handler) handleToolCall(req *Request) (*ToolCallResult, *Error) {
	start := time.Now()

	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		return nil, &Error{Code: InvalidParams, Message: "Invalid params"}
	}

	var params ToolCallParams
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return nil, &Error{Code: InvalidParams, Message: "Invalid params structure"}
	}

	executor, ok := h.tools[params.Name]
	if !ok {
		return nil, &Error{Code: InvalidParams, Message: fmt.Sprintf("Unknown tool: %s", params.Name)}
	}

	// Extract module name from tool name (e.g., "supabase_run_query" -> "supabase")
	module := "unknown"
	for i, c := range params.Name {
		if c == '_' {
			module = params.Name[:i]
			break
		}
	}

	result, err := executor.Execute(params.Arguments)
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		observability.LogToolCall(module, params.Name, durationMs, "error", err.Error())
		return &ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: err.Error()}},
			IsError: true,
		}, nil
	}

	observability.LogToolCall(module, params.Name, durationMs, "success", "")
	return &ToolCallResult{
		Content: []ContentBlock{{Type: "text", Text: result}},
	}, nil
}

