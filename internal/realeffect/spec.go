package realeffect

// =========================
// MODELO DA MISSÃO (SPEC)
// =========================

// Context descreve o contexto básico da missão.
type Context struct {
	Title          string `yaml:"title"           json:"title"`
	Description    string `yaml:"description"     json:"description"`
	Location       string `yaml:"location"        json:"location"`
	TimeframeStart string `yaml:"timeframe_start" json:"timeframe_start"`
	TimeframeEnd   string `yaml:"timeframe_end"   json:"timeframe_end"`
	Owner          string `yaml:"owner"           json:"owner"`
}

// ParticipantRole define papéis mínimos/máximos.
type ParticipantRole struct {
	Role string `yaml:"role" json:"role"`
	Min  int    `yaml:"min"  json:"min"`
	Max  *int   `yaml:"max,omitempty" json:"max,omitempty"`
}

// EvidenceSlot define um tipo de prova exigida pela missão.
type EvidenceSlot struct {
	ID          string  `yaml:"id"          json:"id"`
	Description string  `yaml:"description" json:"description"`
	Category    string  `yaml:"category"    json:"category"`
	Weight      float64 `yaml:"weight"      json:"weight"`
	Required    bool    `yaml:"required"    json:"required"`
}

// MissionSpec é a especificação completa da missão.
type MissionSpec struct {
	SpecVersion   string            `yaml:"spec_version"   json:"spec_version"`
	MissionID     string            `yaml:"mission_id"     json:"mission_id"`
	Context       Context           `yaml:"context"        json:"context"`
	Participants  []ParticipantRole `yaml:"participants"   json:"participants"`
	EvidenceSlots []EvidenceSlot    `yaml:"evidence_slots" json:"evidence_slots"`
}
