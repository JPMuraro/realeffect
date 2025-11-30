package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	re "realeffect-cli/internal/realeffect"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: realeffectc <mission-file.reff> [scenario]\n")
		fmt.Fprintf(os.Stderr, "scenarios: all-accepted (default), missing-proof, low-acceptance\n")
		os.Exit(1)
	}

	path := os.Args[1]
	scenario := "all-accepted"
	if len(os.Args) >= 3 {
		scenario = os.Args[2]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("error resolving path %q: %v\n", path, err)
	}

	fmt.Println("RealEffect CLI — core evaluation v0.1")
	fmt.Println("Mission file:", absPath)
	fmt.Println("Scenario    :", scenario)

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("error reading file: %v\n", err)
	}

	var ms re.MissionSpec
	if err := yaml.Unmarshal(data, &ms); err != nil {
		log.Fatalf("error parsing YAML: %v\n", err)
	}

	if err := re.ValidateSpec(ms); err != nil {
		log.Fatalf("spec is INVALID (RealEffect core): %v\n", err)
	}

	fmt.Println("Spec is structurally VALID (RealEffect core).")

	input := re.BuildScenarioInput(ms, scenario)
	result := re.Evaluate(ms, input)

	fmt.Printf(
		"Evaluation result: valid=%v ratio=%.2f accepted=%.2f rejected=%.2f\nReason: %s\n",
		result.Valid,
		result.Ratio,
		result.AcceptedWeight,
		result.RejectedWeight,
		result.Reason,
	)
}
