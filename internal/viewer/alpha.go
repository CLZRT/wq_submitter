package viewer

import "encoding/json"

type Alpha struct {
	ID             int64
	SimulationEnv  json.RawMessage
	Alpha          string
	IdeaID         int64
	SimulationData json.RawMessage
	TestPeriod     string
	IsSubmitted    int64
	IsDeleted      int64
}

type UploadAlphaList struct {
	Type     string   `json:"type"`
	Regular  string   `json:"regular"`
	Settings AlphaEnv `json:"settings"`
}
type AlphaEnv struct {
	InstrumentType string  `json:"instrumentType"`
	Region         string  `json:"region"`
	Universe       string  `json:"universe"`
	Delay          int64   `json:"delay"`
	Decay          int64   `json:"decay"`
	Neutralization string  `json:"neutralization"`
	Truncation     float64 `json:"truncation"`
	Pasteurization string  `json:"pasteurization"`
	NanHandling    string  `json:"nanHandling"`
	Language       string  `json:"language"`
	TestPeriod     string  `json:"testPeriod,omitempty"`
	Visualization  bool    `json:"visualization"`
	UnitHandling   string  `json:"unitHandling"`
}
