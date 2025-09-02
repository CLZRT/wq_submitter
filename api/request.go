package api

type UploadAlphaListWithIdeaReq struct {
	Idea      UploadIdeaReq        `json:"idea"`
	AlphaList []UploadAlphaListReq `json:"alphaList"`
}
type UploadAlphaListReq struct {
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

type DeleteIdeaReq struct {
	Id int64 `json:"id"`
}

type UploadIdeaReq struct {
	IdeaAlphaTemplate string `json:"ideaAlphaTemplate"`
	IdeaTitle         string `json:"ideaTitle"`
	IdeaDesc          string `json:"ideaDesc"`
	Id                int64  `json:"id"`
	ConcurrencyNum    int64  `json:"concurrencyNum"`
}

type UpdateIdeaReq struct {
	Id             int64 `json:"id"`
	ConcurrencyNum int64 `json:"num"`
}
