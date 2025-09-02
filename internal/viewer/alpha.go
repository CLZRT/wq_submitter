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
