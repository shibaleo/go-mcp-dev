package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoJSON_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("expected Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected Content-Type header for POST with body")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	client := New()
	headers := map[string]string{"Authorization": "Bearer test-token"}
	body := map[string]string{"key": "value"}

	resp, err := client.DoJSON("POST", server.URL, headers, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resp) != `{"result": "success"}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestDoJSON_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	client := New()
	_, err := client.DoJSON("GET", server.URL, nil, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestDoJSON_GetWithoutBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Error("GET without body should not have Content-Type")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := New()
	resp, err := client.DoJSON("GET", server.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resp) != `[]` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestPrettyJSON(t *testing.T) {
	input := []byte(`{"name":"test","value":123}`)
	result := PrettyJSON(input)

	expected := `{
  "name": "test",
  "value": 123
}`
	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPrettyJSON_InvalidJSON(t *testing.T) {
	input := []byte(`not json`)
	result := PrettyJSON(input)

	if result != "not json" {
		t.Errorf("expected original string, got: %s", result)
	}
}
