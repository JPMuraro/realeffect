package realeffect

import "testing"

// helper para criar uma missão simples válida de 5 slots de 0.2
func makeValidMissionSpec() MissionSpec {
	return MissionSpec{
		SpecVersion: "0.1",
		MissionID:   "test_mission",
		Context: Context{
			Title:          "Teste",
			Description:    "Missão de teste",
			Location:       "Lugar nenhum",
			TimeframeStart: "2025-01-01",
			TimeframeEnd:   "2025-12-31",
			Owner:          "test_owner",
		},
		Participants: []ParticipantRole{
			{Role: "tester", Min: 1, Max: intPtr(1)},
		},
		EvidenceSlots: []EvidenceSlot{
			{ID: "slot1", Description: "Slot 1", Category: "doc", Weight: 0.2, Required: true},
			{ID: "slot2", Description: "Slot 2", Category: "doc", Weight: 0.2, Required: true},
			{ID: "slot3", Description: "Slot 3", Category: "doc", Weight: 0.2, Required: true},
			{ID: "slot4", Description: "Slot 4", Category: "doc", Weight: 0.2, Required: true},
			{ID: "slot5", Description: "Slot 5", Category: "doc", Weight: 0.2, Required: true},
		},
	}
}

func intPtr(v int) *int {
	return &v
}

func TestValidateSpec_ValidSpec(t *testing.T) {
	ms := makeValidMissionSpec()

	if err := ValidateSpec(ms); err != nil {
		t.Fatalf("expected spec to be valid, got error: %v", err)
	}
}

func TestValidateSpec_InvalidWeightPerSlot(t *testing.T) {
	ms := makeValidMissionSpec()
	ms.EvidenceSlots[0].Weight = 0.8 // maior que MaxWeightPerSlot (0.4)

	if err := ValidateSpec(ms); err == nil {
		t.Fatalf("expected error for slot weight > MaxWeightPerSlot, got nil")
	}
}

func TestEvaluate_AllAccepted_Passes(t *testing.T) {
	ms := makeValidMissionSpec()

	input := BuildScenarioInput(ms, "all-accepted")
	res := Evaluate(ms, input)

	if !res.Valid {
		t.Fatalf("expected mission to be valid, got invalid. Reason: %s", res.Reason)
	}
	if res.Ratio < MinAcceptedRatio {
		t.Fatalf("expected ratio >= %.2f, got %.2f", MinAcceptedRatio, res.Ratio)
	}
}

func TestEvaluate_MissingProof_FailsRE0(t *testing.T) {
	ms := makeValidMissionSpec()

	input := BuildScenarioInput(ms, "missing-proof")
	res := Evaluate(ms, input)

	if res.Valid {
		t.Fatalf("expected mission to be invalid due to missing proof, but got valid")
	}
	if res.Reason == "" {
		t.Fatalf("expected a reason for failure, got empty string")
	}
}

func TestEvaluate_LowAcceptance_FailsRE1(t *testing.T) {
	ms := makeValidMissionSpec()

	input := BuildScenarioInput(ms, "low-acceptance")
	res := Evaluate(ms, input)

	if res.Valid {
		t.Fatalf("expected mission to be invalid due to low acceptance ratio, but got valid")
	}
	if res.Ratio >= MinAcceptedRatio {
		t.Fatalf("expected ratio < %.2f, got %.2f", MinAcceptedRatio, res.Ratio)
	}
}
