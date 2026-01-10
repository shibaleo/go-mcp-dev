package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

func init() {
	// Register a test module for testing
	modules.Register(modules.ModuleDefinition{
		Name:        "test",
		Description: "Test module",
		Tools: []modules.Tool{
			{
				Name:        "echo",
				Description: "Echo back the input",
				InputSchema: modules.InputSchema{
					Type: "object",
					Properties: map[string]modules.Property{
						"message": {Type: "string", Description: "Message to echo"},
					},
				},
			},
		},
		Handlers: map[string]modules.ToolHandler{
			"echo": func(params map[string]interface{}) (string, error) {
				msg, _ := params["message"].(string)
				return "Echo: " + msg, nil
			},
		},
	})
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

	// Should return 2 meta tools: get_module_schema and call_module_tool
	if len(tools) != 2 {
		t.Errorf("expected 2 meta tools, got %d", len(tools))
	}
}

func TestHandleInlineMessage_GetModuleSchema(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_module_schema","arguments":{"module":"test"}}}`

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

	content, ok := resultMap["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("expected content in result")
	}
}

func TestHandleInlineMessage_CallModuleTool(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"call_module_tool","arguments":{"module":"test","tool_name":"echo","params":{"message":"hello"}}}}`

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

	content, ok := resultMap["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("expected content in result")
	}

	firstContent := content[0].(map[string]interface{})
	text := firstContent["text"].(string)
	if text != "Echo: hello" {
		t.Errorf("expected 'Echo: hello', got '%s'", text)
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

func TestHandleInlineMessage_UnknownModule(t *testing.T) {
	handler := NewHandler()

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_module_schema","arguments":{"module":"nonexistent"}}}`

	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleInlineMessage(rec, req)

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected JSON-RPC error: %v", resp.Error)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected result type: %T", resp.Result)
	}

	isError, _ := resultMap["isError"].(bool)
	if !isError {
		t.Error("expected isError to be true for unknown module")
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
