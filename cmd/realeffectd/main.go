package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	re "realeffect-cli/internal/realeffect"

	"gopkg.in/yaml.v3"
)

// Agora temos 2 jeitos de enviar a missão:
// 1) spec      -> YAML em string
// 2) spec_path -> caminho para arquivo .reff no disco
type evalRequest struct {
	Spec     string `json:"spec,omitempty"`      // YAML bruto (opcional)
	SpecPath string `json:"spec_path,omitempty"` // caminho para arquivo .reff (opcional)
	Scenario string `json:"scenario,omitempty"`  // all-accepted, missing-proof, low-acceptance
}

type evalResponse struct {
	Valid          bool    `json:"valid"`
	Ratio          float64 `json:"ratio"`
	AcceptedWeight float64 `json:"accepted_weight"`
	RejectedWeight float64 `json:"rejected_weight"`
	Reason         string  `json:"reason,omitempty"`
	Error          string  `json:"error,omitempty"`
}

func evaluateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req evalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON body: %v", err), http.StatusBadRequest)
		return
	}

	if req.Scenario == "" {
		req.Scenario = "all-accepted"
	}

	// Decidir de onde vem a spec: texto direto ou arquivo
	specText, err := loadSpec(req)
	if err != nil {
		resp := evalResponse{
			Error: err.Error(),
		}
		writeJSON(w, http.StatusBadRequest, resp)
		return
	}

	var ms re.MissionSpec
	if err := yaml.Unmarshal([]byte(specText), &ms); err != nil {
		resp := evalResponse{
			Error: fmt.Sprintf("error parsing YAML spec: %v", err),
		}
		writeJSON(w, http.StatusBadRequest, resp)
		return
	}

	if err := re.ValidateSpec(ms); err != nil {
		resp := evalResponse{
			Error: fmt.Sprintf("spec is INVALID (RealEffect core): %v", err),
		}
		writeJSON(w, http.StatusBadRequest, resp)
		return
	}

	input := re.BuildScenarioInput(ms, req.Scenario)
	result := re.Evaluate(ms, input)

	resp := evalResponse{
		Valid:          result.Valid,
		Ratio:          result.Ratio,
		AcceptedWeight: result.AcceptedWeight,
		RejectedWeight: result.RejectedWeight,
		Reason:         result.Reason,
	}

	writeJSON(w, http.StatusOK, resp)
}

func loadSpec(req evalRequest) (string, error) {
	// Prioridade:
	// 1) Spec em texto (se preenchido)
	// 2) Spec via caminho de arquivo
	if strings.TrimSpace(req.Spec) != "" {
		return req.Spec, nil
	}
	if strings.TrimSpace(req.SpecPath) != "" {
		abs, err := filepath.Abs(req.SpecPath)
		if err != nil {
			return "", fmt.Errorf("invalid spec_path: %v", err)
		}
		data, err := os.ReadFile(abs)
		if err != nil {
			return "", fmt.Errorf("cannot read spec_path %q: %v", abs, err)
		}
		return string(data), nil
	}
	return "", fmt.Errorf("either 'spec' (YAML) or 'spec_path' must be provided")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error writing JSON response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/evaluate", evaluateHandler)

	addr := ":8081"
	log.Printf("RealEffectD — HTTP service listening on %s", addr)
	log.Printf(`POST /evaluate with JSON, e.g.: {"spec_path": ".\\examples\\plant_100_trees.reff", "scenario": "all-accepted"}`)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
