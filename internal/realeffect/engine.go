package realeffect

import (
	"errors"
	"fmt"
)

// =============================
// ENGINE — REGRAS RE-0 e RE-1
// =============================

const (
	// Regra núcleo da linguagem: mínimo 80% de peso aceito.
	MinAcceptedRatio = 0.8

	// Limites para evitar specs desbalanceadas.
	MaxWeightPerSlot    = 0.4  // um único slot não pode ter mais que 40% do peso total
	MinRelevantWeight   = 0.10 // peso mínimo para ser considerado relevante
	MinRelevantEvidence = 3    // pelo menos 3 slots relevantes
)

// EvidenceStatus representa o estado de um slot de prova.
type EvidenceStatus string

const (
	StatusMissing   EvidenceStatus = "MISSING"
	StatusSubmitted EvidenceStatus = "SUBMITTED"
	StatusAccepted  EvidenceStatus = "ACCEPTED"
	StatusRejected  EvidenceStatus = "REJECTED"
)

// EvaluationInput representa as provas entregues pelos participantes.
// participantID -> slotID -> status
type EvaluationInput struct {
	Slots map[string]map[string]EvidenceStatus
}

// EvaluationResult é o resultado da avaliação RE-0/RE-1.
type EvaluationResult struct {
	Valid          bool
	Ratio          float64
	AcceptedWeight float64
	RejectedWeight float64
	Reason         string
}

// ValidateSpec aplica as regras estruturais da linguagem (peso, slots, etc).
func ValidateSpec(ms MissionSpec) error {
	if ms.MissionID == "" {
		return errors.New("mission_id is required")
	}
	if len(ms.EvidenceSlots) == 0 {
		return errors.New("at least one evidence slot is required")
	}

	var totalWeight float64
	for _, slot := range ms.EvidenceSlots {
		if slot.ID == "" {
			return errors.New("evidence slot id cannot be empty")
		}
		if slot.Weight <= 0 {
			return fmt.Errorf("evidence slot %q must have positive weight", slot.ID)
		}
		if slot.Weight > MaxWeightPerSlot {
			return fmt.Errorf(
				"evidence slot %q weight %.2f exceeds max allowed %.2f",
				slot.ID, slot.Weight, MaxWeightPerSlot,
			)
		}
		totalWeight += slot.Weight
	}

	if totalWeight <= 0 {
		return errors.New("total evidence weight must be > 0")
	}

	var relevantCount int
	for _, slot := range ms.EvidenceSlots {
		norm := slot.Weight / totalWeight
		if norm >= MinRelevantWeight {
			relevantCount++
		}
	}

	if relevantCount < MinRelevantEvidence {
		return fmt.Errorf(
			"at least %d evidence slots must have normalized weight >= %.2f, got %d",
			MinRelevantEvidence, MinRelevantWeight, relevantCount,
		)
	}

	return nil
}

// Evaluate aplica:
// RE-0: 100% de entrega (nenhum slot faltando)
// RE-1: 80% de peso aceito (norma global >= 0.8)
func Evaluate(ms MissionSpec, input EvaluationInput) EvaluationResult {
	// RE-0: checa entrega (ninguém pode ficar com MISSING)
	for participantID, slots := range input.Slots {
		for _, slot := range ms.EvidenceSlots {
			status, ok := slots[slot.ID]
			if !ok || status == StatusMissing {
				return EvaluationResult{
					Valid:  false,
					Reason: fmt.Sprintf("participant %s missing required evidence slot %s", participantID, slot.ID),
				}
			}
		}
	}

	var totalWeight float64
	for _, slot := range ms.EvidenceSlots {
		totalWeight += slot.Weight
	}

	if totalWeight <= 0 {
		return EvaluationResult{
			Valid:  false,
			Reason: "invalid spec: total weight <= 0",
		}
	}

	var acceptedWeightNorm float64
	var rejectedWeightNorm float64

	for _, slot := range ms.EvidenceSlots {
		norm := slot.Weight / totalWeight

		slotHasAccepted := false
		slotHasRejected := false

		for _, slots := range input.Slots {
			status := slots[slot.ID]
			if status == StatusAccepted {
				slotHasAccepted = true
			}
			if status == StatusRejected {
				slotHasRejected = true
			}
		}

		if slotHasRejected {
			rejectedWeightNorm += norm
		} else if slotHasAccepted {
			acceptedWeightNorm += norm
		}
	}

	ratio := acceptedWeightNorm

	if ratio >= MinAcceptedRatio {
		return EvaluationResult{
			Valid:          true,
			Ratio:          ratio,
			AcceptedWeight: acceptedWeightNorm,
			RejectedWeight: rejectedWeightNorm,
			Reason:         "mission meets RealEffect 80% acceptance rule",
		}
	}

	return EvaluationResult{
		Valid:          false,
		Ratio:          ratio,
		AcceptedWeight: acceptedWeightNorm,
		RejectedWeight: rejectedWeightNorm,
		Reason:         "mission failed RealEffect 80% acceptance rule",
	}
}

// BuildScenarioInput constrói diferentes estados de provas
// para demonstrar sucesso e falhas da linguagem.
func BuildScenarioInput(ms MissionSpec, scenario string) EvaluationInput {
	slots := map[string]map[string]EvidenceStatus{
		"participant_1": {},
	}

	switch scenario {
	case "missing-proof":
		// Cenário RE-0: pelo menos um slot fica como MISSING.
		first := true
		for _, slot := range ms.EvidenceSlots {
			if first {
				slots["participant_1"][slot.ID] = StatusMissing
				first = false
				continue
			}
			slots["participant_1"][slot.ID] = StatusAccepted
		}

	case "low-acceptance":
		// Cenário RE-1: todos os slots são entregues,
		// mas parte do peso é rejeitado para ficar < 80%.
		// Para o exemplo plant_100_trees (5 slots de 0.2),
		// aceitar 3 e rejeitar 2 => 60% de aceite.
		accepted := 0
		for _, slot := range ms.EvidenceSlots {
			if accepted < 3 {
				slots["participant_1"][slot.ID] = StatusAccepted
				accepted++
			} else {
				slots["participant_1"][slot.ID] = StatusRejected
			}
		}

	default:
		// Cenário padrão: tudo ACEITO (sucesso total).
		for _, slot := range ms.EvidenceSlots {
			slots["participant_1"][slot.ID] = StatusAccepted
		}
	}

	return EvaluationInput{Slots: slots}
}
