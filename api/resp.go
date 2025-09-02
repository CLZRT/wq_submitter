package api

import (
	"wq_submitter/internal/viewer"
)

type UploadAlphaListWithIdeaResp struct {
	Message   string `json:"message"`
	UploadNum int64  `json:"uploadNum"`
	IdeaId    int64  `json:"ideaId"`
}

type GetAlphaListResp struct {
	Message   string         `json:"message"`
	IdeaId    int64          `json:"ideaId"`
	AlphaList []viewer.Alpha `json:"alphaList"`
}

type UpdateIdeaResp struct {
	Message        string `json:"message"`
	IdeaId         int64  `json:"ideaId"`
	ConcurrencyNum int64  `json:"concurrencyNum"`
}

type GetAllIdeaResp struct {
	Message string        `json:"message"`
	Ideas   []viewer.Idea `json:"ideas"`
}

type GetUnFinishedIdeaResp struct {
	Message string        `json:"message"`
	Ideas   []viewer.Idea `json:"ideas"`
}

type GetRunningIdeaResp struct {
	Message string        `json:"message"`
	Ideas   []viewer.Idea `json:"ideas"`
}

type DeleteIdeaResp struct {
	Message string      `json:"message"`
	Idea    viewer.Idea `json:"idea"`
}
