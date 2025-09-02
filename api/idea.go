package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"wq_submitter/internal/svc"
	"wq_submitter/internal/viewer"
)

func GetAllIdea(ctx *gin.Context) {
	ideas, err := svc.GetAllIdea()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, GetAllIdeaResp{
			Message: "GetAllIdea Failed || Service Error",
			Ideas:   nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, GetAllIdeaResp{
		Message: "GetAllIdea Success",
		Ideas:   ideas,
	})

}

func GetUnfinishedIdea(ctx *gin.Context) {
	ideas, err := svc.GetUnfinishedIdea()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, GetUnFinishedIdeaResp{
			Message: "GetUnfinishedIdea Failed || Service Error",
			Ideas:   nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, GetUnFinishedIdeaResp{
		Message: "GetUnfinishedIdea Success",
		Ideas:   ideas,
	})
}

func GetRunningIdea(ctx *gin.Context) {
	ideas, err := svc.GetNeedRunIdea()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, GetRunningIdeaResp{
			Message: "GetNeedRunIdea Failed || Service Error",
			Ideas:   nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, GetRunningIdeaResp{
		Message: "GetNeedRunIdea Success",
		Ideas:   ideas,
	})

}

func UpdateIdeaConcurrency(ctx *gin.Context) {
	var ideaRequest UpdateIdeaReq
	err := ctx.ShouldBind(&ideaRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, UpdateIdeaResp{
			Message:        "Param Error",
			IdeaId:         -1,
			ConcurrencyNum: -1,
		})
		return
	}

	//判断是否合法
	if !svc.IsLegalConcurrency(ideaRequest.Id, ideaRequest.ConcurrencyNum) {
		ctx.JSON(http.StatusBadRequest, UpdateIdeaResp{
			Message:        "Concurrency Num is Over Limit.",
			IdeaId:         -1,
			ConcurrencyNum: -1,
		})
		return
	}

	ideaViewer := &viewer.Idea{}
	ideaViewer.ID = ideaRequest.Id
	ideaViewer.ConcurrencyNum = ideaRequest.ConcurrencyNum
	ideaViewerByte, err := json.Marshal(ideaViewer)

	if err != nil {
		log.Error(err)
	}
	log.Infof("Change Concurrency %s ", string(ideaViewerByte))

	ideaId, err := svc.UpdateIdea(ideaViewer)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadGateway, UpdateIdeaResp{
			Message:        "Server Error",
			IdeaId:         -1,
			ConcurrencyNum: -1,
		})
		return
	}

	ctx.JSON(http.StatusOK, UpdateIdeaResp{
		Message:        "Success",
		IdeaId:         ideaId,
		ConcurrencyNum: ideaRequest.ConcurrencyNum,
	})

}

// 根据Idea ID删除Idea,并且关闭其运行池，删除其alphaList
func DeleteIdea(ctx *gin.Context) {
	req := DeleteIdeaReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, DeleteIdeaResp{
			Message: "Bad Param",
			Idea:    viewer.Idea{},
		})
		return
	}

	ideaViewer, err := svc.DeleteIdeaWithTx(req.Id)
	if err != nil || ideaViewer.ID == 0 {
		ctx.JSON(http.StatusBadGateway, DeleteIdeaResp{
			Message: "Service Error || DeleteIdea",
			Idea:    viewer.Idea{},
		})
		return
	}

	ctx.JSON(http.StatusOK, DeleteIdeaResp{
		Message: "Success",
		Idea:    ideaViewer,
	})

}
