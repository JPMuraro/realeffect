package realeffect

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper para criar um Client apontando pro servidor de teste
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})

	return &Client{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
}

func TestClient_EvaluateFromFile_Success(t *testing.T) {
	// servidor HTTP falso que simula o realeffectd respondendo sucesso
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/evaluate" {
			t.Errorf("expected path /evaluate, got %s", r.URL.Path)
		}

		var req evalRequestHTTP
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.SpecPath == "" {
			t.Errorf("expected SpecPath to be set")
		}

		resp := evalResponseHTTP{
			Valid:          true,
			Ratio:          1.0,
			AcceptedWeight: 1.0,
			RejectedWeight: 0.0,
			Reason:         "ok-from-test",
			Error:          "",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}

	client := newTestClient(t, handler)

	res, err := client.EvaluateFromFile("./examples/plant_100_trees.reff", "all-accepted")
	if err != nil {
		t.Fatalf("EvaluateFromFile returned error: %v", err)
	}

	if !res.Valid {
		t.Errorf("expected Valid=true, got false")
	}
	if res.Ratio != 1.0 {
		t.Errorf("expected Ratio=1.0, got %v", res.Ratio)
	}
	if res.Reason != "ok-from-test" {
		t.Errorf("unexpected Reason: %s", res.Reason)
	}
}

func TestClient_EvaluateFromFile_ServerErrorField(t *testing.T) {
	// servidor que devolve "error" preenchido (simulando falha do realeffectd)
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := evalResponseHTTP{
			Valid: false,
			Error: "spec is INVALID (test-error)",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}

	client := newTestClient(t, handler)

	_, err := client.EvaluateFromFile("./examples/plant_100_trees.reff", "all-accepted")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestClient_EvaluateFromFile_EmptySpecPath(t *testing.T) {
	client := &Client{
		BaseURL: "http://example.com",
	}

	_, err := client.EvaluateFromFile("", "all-accepted")
	if err == nil {
		t.Fatalf("expected error for empty specPath, got nil")
	}
}
