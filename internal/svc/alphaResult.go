package svc

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"wq_submitter/internal/model"
	"wq_submitter/internal/repo"
	"wq_submitter/internal/viewer"
)

var (
	alphaResultRepo *repo.AlphaResultRepo
)

func init() {
	alphaResultRepo = repo.NewAlphaResultRepo()
}

func StoreAlphaResult(alphaRequestViewer *viewer.AlphaResult) error {
	alphaResultModel := model.AlphaResult{}

	alphaResultModel.AlphaId = alphaRequestViewer.AlphaId
	alphaResultModel.IdeaId = alphaRequestViewer.IdeaId
	alphaResultModel.AlphaDetail = datatypes.JSON(alphaRequestViewer.AlphaDetail)
	alphaResultModel.AlphaCode = alphaRequestViewer.AlphaCode

	alphaResultModel.BasicResult = datatypes.JSON(alphaRequestViewer.BasicResult)
	alphaResultModel.CheckResult = datatypes.JSON(alphaRequestViewer.CheckResult)
	alphaResultModel.SelfCorrelation = datatypes.JSON(alphaRequestViewer.SelfCorrelation)
	alphaResultModel.ProdCorrelation = datatypes.JSON(alphaRequestViewer.ProdCorrelation)
	alphaResultModel.Turnover = datatypes.JSON(alphaRequestViewer.Turnover)
	alphaResultModel.Sharpe = datatypes.JSON(alphaRequestViewer.Sharpe)
	alphaResultModel.Pnl = datatypes.JSON(alphaRequestViewer.Pnl)
	alphaResultModel.DailyPnl = datatypes.JSON(alphaRequestViewer.DailyPnl)
	alphaResultModel.YearlyStats = datatypes.JSON(alphaRequestViewer.YearlyStats)

	rowsAffected, err := alphaResultRepo.Add(context.Background(), &alphaResultModel)
	if err != nil {
		log.Errorf("StoreAlphaResult Add Failed || %s", err.Error())
		return err
	}
	log.Infof("AlphaResult %d rows affected", rowsAffected)
	return nil

}
