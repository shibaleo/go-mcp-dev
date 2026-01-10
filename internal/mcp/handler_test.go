package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockTool for testing
type MockTool struct {
	name   string
	result string
	err    error
}

func (m *MockTool) Definition() Tool {
	return Tool{
		Name:        m.name,
		Description: "Mock tool for testing",
		InputSchema: InputSchema{
			Type:       "object",
			Properties: map[string]Property{},
		},
	}
}

func (m *MockTool) Execute(args map[string]interface{}) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.result, nil
}

func TestHandleInlineMessage_Initialize(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"test","version":"1.0"}}}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(*InitializeResult)
	if !ok {
		// Result might be map[string]interface{} after JSON round-trip
		resultMap, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatalf("unexpected result type: %T", resp.Result)
		}
		if resultMap["protocolVersion"] != "2024-11-05" {
			t.Errorf("unexpected protocol version: %v", resultMap["protocolVersion"])
		}
	} else {
		if result.ProtocolVersion != "2024-11-05" {
			t.Errorf("expected protocol version 2024-11-05, got %s", result.ProtocolVersion)
		}
	}
}

func TestHandleInlineMessage_ToolsList(t *testing.T) {
	handler := NewHandler()
	handler.RegisterTool(&MockTool{name: "test_tool", result: "ok"})

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error: %v", resp.Error)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected result type: %T", resp.Result)
	}

	tools, ok := resultMap["tools"].([]interface{})
	if !ok {
		t.Fatalf("unexpected tools type: %T", resultMap["tools"])
	}

	if len(tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(tools))
	}
}

func TestHandleInlineMessage_ToolCall(t *testing.T) {
	handler := NewHandler()
	handler.RegisterTool(&MockTool{name: "test_tool", result: "success result"})

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"test_tool","arguments":{}}}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error: %v", resp.Error)
	}
}

func TestHandleInlineMessage_ParseError(t *testing.T) {
	handler := NewHandler()

	reqBody := `{invalid json`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}

	if resp.Error.Code != ParseError {
		t.Errorf("expected error code %d, got %d", ParseError, resp.Error.Code)
	}
}

func TestHandleInlineMessage_MethodNotFound(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"unknown/method"}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}

	if resp.Error.Code != MethodNotFound {
		t.Errorf("expected error code %d, got %d", MethodNotFound, resp.Error.Code)
	}
}

func TestHandleInlineMessage_UnknownTool(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nonexistent_tool","arguments":{}}}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}

	if resp.Error.Code != InvalidParams {
		t.Errorf("expected error code %d, got %d", InvalidParams, resp.Error.Code)
	}
}

func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest("DELETE", "/mcp", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestHandleInlineMessage_ToolCallError(t *testing.T) {
	handler := NewHandler()
	handler.RegisterTool(&MockTool{
		name: "error_tool",
		err:  fmt.Errorf("tool execution failed"),
	})

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"error_tool","arguments":{}}}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Tool errors are returned as result with isError=true, not as JSON-RPC error
	if resp.Error != nil {
		t.Errorf("unexpected JSON-RPC error: %v", resp.Error)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected result type: %T", resp.Result)
	}

	isError, ok := resultMap["isError"].(bool)
	if !ok || !isError {
		t.Error("expected isError to be true")
	}
}

func TestHandleInlineMessage_InvalidParams(t *testing.T) {
	handler := NewHandler()

	// tools/call with missing params
	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call"}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}

	if resp.Error.Code != InvalidParams {
		t.Errorf("expected error code %d, got %d", InvalidParams, resp.Error.Code)
	}
}
