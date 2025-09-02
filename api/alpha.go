package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"wq_submitter/internal/svc"
	"wq_submitter/internal/viewer"
)

func Hello(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "Hello")
}

func UploadAlphaListWithIdea(ctx *gin.Context) {
	//todo 上传Idea及其AlphaList
	alphaIdeaRequest := UploadAlphaListWithIdeaReq{}
	err := ctx.ShouldBind(&alphaIdeaRequest)
	if err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, UploadAlphaListWithIdeaResp{
			Message:   err.Error(),
			UploadNum: 0,
		})
		return
	}

	ideaReq := alphaIdeaRequest.Idea
	alphaRequestList := alphaIdeaRequest.AlphaList

	// Check Idea
	if ideaReq.IdeaAlphaTemplate == "" {
		ctx.JSON(http.StatusBadRequest, UploadAlphaListWithIdeaResp{
			Message:   "IdeaAlphaTemplate is Empty",
			UploadNum: 0,
		})
		return
	}

	//检查并发是否合法
	if !svc.IsLegalConcurrency(ideaReq.Id, ideaReq.ConcurrencyNum) {
		ctx.JSON(http.StatusBadRequest, UploadAlphaListWithIdeaResp{
			Message:   "Concurrency Num is Over Limit.",
			UploadNum: 0,
		})
		return
	}

	// 获取LastAlphaId
	lastAlphaId := svc.FindLastAlphaId()

	// Convert Idea AND Store
	ideaViewer := viewer.Idea{}
	ideaViewer.IdeaAlphaTemplate = ideaReq.IdeaAlphaTemplate
	ideaViewer.IdeaTitle = ideaReq.IdeaTitle
	ideaViewer.IdeaDesc = ideaReq.IdeaDesc
	ideaViewer.StartIdx = lastAlphaId + 1
	ideaViewer.EndIdx = ideaViewer.StartIdx + int64(len(alphaRequestList)) - 1
	ideaViewer.NextIdx = ideaViewer.StartIdx
	ideaViewer.ConcurrencyNum = ideaReq.ConcurrencyNum

	ideaId, err := svc.UploadIdea(&ideaViewer)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if ideaId <= 0 {
		log.Errorf("Return InValid IdeaId {%d}", ideaId)
		return
	}

	log.Infof("UpLoad Idea success, ideaId: {%d}", ideaId)

	// Convert alphaList
	alphaModelList := make([]viewer.Alpha, 0)
	for _, alphaRequest := range alphaRequestList {

		// Marshal
		settingByte, err := json.Marshal(alphaRequest.Settings)
		if err != nil {
			log.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, UploadAlphaListWithIdeaResp{
				Message:   "Settings Marshal Failed",
				UploadNum: 0,
			})
			return
		}
		alphaDataByte, err := json.Marshal(alphaRequest)
		if err != nil {
			log.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, UploadAlphaListWithIdeaResp{
				Message:   "Request Marshal Failed",
				UploadNum: 0,
			})
			return
		}

		alphaViewer := viewer.Alpha{
			SimulationEnv:  settingByte,
			Alpha:          alphaRequest.Regular,
			SimulationData: alphaDataByte,
			IdeaID:         ideaId,
			TestPeriod:     alphaRequest.Settings.TestPeriod,
		}
		alphaModelList = append(alphaModelList, alphaViewer)

	}
	// Upload alphaList
	uploadNum := svc.UploadAlphaList(alphaModelList)
	if uploadNum == -1 {
		ctx.JSON(http.StatusBadGateway, UploadAlphaListWithIdeaResp{
			Message:   "UpLoad AlphaList Failed",
			UploadNum: uploadNum,
			IdeaId:    ideaId,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, UploadAlphaListWithIdeaResp{
			Message:   "UpLoad AlphaList Success",
			UploadNum: uploadNum,
			IdeaId:    ideaId,
		})
		return
	}

}

func GetAlphaListByIdea(ctx *gin.Context) {

	ideaIdStr := ctx.Query("ideaId")
	if ideaIdStr == "" {
		ctx.JSON(http.StatusBadRequest, GetAlphaListResp{
			Message:   "Need Params Idea's Id",
			IdeaId:    0,
			AlphaList: nil,
		})
		return
	}

	ideaId, err := strconv.Atoi(ideaIdStr)
	if err != nil || ideaId <= 0 {
		ctx.JSON(http.StatusBadRequest, GetAlphaListResp{
			Message:   "Bad Params Idea's Id",
			IdeaId:    0,
			AlphaList: nil,
		})
		return
	}

	alphaList, err := svc.FindAlphaListByIdeaId(int64(ideaId))
	if err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadGateway, GetAlphaListResp{
			Message:   "Server Error",
			IdeaId:    0,
			AlphaList: nil,
		})
		return

	}

	log.Info("API|| GetAlphaList Success")
	ctx.JSON(http.StatusOK, GetAlphaListResp{
		Message:   "Success",
		IdeaId:    int64(ideaId),
		AlphaList: alphaList,
	})

}
