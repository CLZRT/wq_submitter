package svc

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	
	"wq_submitter/internal/model"
	"wq_submitter/internal/repo"
	"wq_submitter/internal/viewer"
)

var alphaRepo *repo.AlphaRepo

func init() {
	alphaRepo = repo.NewAlphaRepo()
}

func UploadAlphaList(alphaViewerList []viewer.Alpha) int64 {

	alphaModelList := make([]*model.Alpha, 0)

	for _, alphaViewer := range alphaViewerList {

		alphaModel := model.Alpha{}
		alphaModel.Alpha = alphaViewer.Alpha
		alphaModel.IdeaID = alphaViewer.IdeaID
		alphaModel.TestPeriod = alphaViewer.TestPeriod
		alphaModel.SimulationData = datatypes.JSON(alphaViewer.SimulationData)
		alphaModel.SimulationEnv = datatypes.JSON(alphaViewer.SimulationEnv)
		alphaModelList = append(alphaModelList, &alphaModel)
	}
	uploadNum, err := alphaRepo.AddList(context.Background(), alphaModelList)
	if err != nil {
		log.Warnf("UploadAlphaList Failed reason: %s", err.Error())
		return -1
	}
	log.Infof("UploadAlphaList Success uploadNum: %d", uploadNum)
	return uploadNum
}
func TxUploadAlphaList(alphaViewerList []viewer.Alpha, ideaId int64, tx *gorm.DB) int64 {

	alphaModelList := make([]*model.Alpha, 0)

	for _, alphaViewer := range alphaViewerList {

		alphaModel := model.Alpha{}
		alphaModel.Alpha = alphaViewer.Alpha
		alphaModel.IdeaID = ideaId
		alphaModel.TestPeriod = alphaViewer.TestPeriod
		alphaModel.SimulationData = datatypes.JSON(alphaViewer.SimulationData)
		alphaModel.SimulationEnv = datatypes.JSON(alphaViewer.SimulationEnv)
		alphaModelList = append(alphaModelList, &alphaModel)
	}
	uploadNum, err := alphaRepo.AddListTx(context.Background(), alphaModelList, tx)
	if err != nil {
		log.Warnf("UploadAlphaList Failed reason: %s", err.Error())
		return -1
	}
	log.Infof("UploadAlphaList Success uploadNum: %d", uploadNum)
	return uploadNum
}

//func FindLastAlphaId() int {
//	repo.
//}

func TxUploadAlphaListWithIdea(alphaViewerList []viewer.Alpha, idea *viewer.Idea) (ideaId int64, alphaNum int64) {

	db := repo.GetDbCli()
	tx := db.Begin()
	tx.Begin()
	ideaId, err := TxUploadIdea(idea, tx)
	if err != nil {
		tx.Rollback()
		return 0, 0
	}

	alphaListNum := TxUploadAlphaList(alphaViewerList, ideaId, tx)
	if alphaListNum <= 0 {
		tx.Rollback()
		return 0, 0
	}
	tx.Commit()
	return ideaId, alphaListNum

}
func FindAlphaListByIdeaId(id int64) ([]viewer.Alpha, error) {
	alphaModels, err := alphaRepo.FindByIdeaId(context.Background(), id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	alphaViewers := make([]viewer.Alpha, 0)

	for _, alphaModel := range alphaModels {

		alphaViewer := viewer.Alpha{}
		alphaViewer.ID = alphaModel.ID
		alphaViewer.SimulationEnv = json.RawMessage(alphaModel.SimulationEnv)
		alphaViewer.IdeaID = alphaModel.IdeaID
		alphaViewer.Alpha = alphaModel.Alpha
		alphaViewer.TestPeriod = alphaModel.TestPeriod
		alphaViewer.SimulationData = json.RawMessage(alphaModel.SimulationData)
		alphaViewer.IsSubmitted = alphaModel.IsSubmitted
		alphaViewer.IsDeleted = alphaModel.IsDeleted

		alphaViewers = append(alphaViewers, alphaViewer)
	}
	return alphaViewers, nil

}

func FindLastAlphaId() int64 {
	lastAlpha, err := alphaRepo.FindLast(context.Background())
	if err != nil {
		log.Error(err)
		return 0
	}
	if lastAlpha == nil {
		return 0
	}
	return lastAlpha.ID

}

func FindAlphaById(alphaId int64) *viewer.Alpha {
	alphaModel, err := alphaRepo.FindById(context.Background(), alphaId)
	if err != nil {
		return nil
	}
	alphaViewer := viewer.Alpha{}
	alphaViewer.ID = alphaModel.ID
	alphaViewer.SimulationEnv = json.RawMessage(alphaModel.SimulationEnv)
	alphaViewer.IdeaID = alphaModel.IdeaID
	alphaViewer.Alpha = alphaModel.Alpha
	alphaViewer.TestPeriod = alphaModel.TestPeriod
	alphaViewer.SimulationData = json.RawMessage(alphaModel.SimulationData)

	return &alphaViewer
}
func FindValidAlphaById(alphaId int64) *viewer.Alpha {
	alphaModel, err := alphaRepo.FindById(context.Background(), alphaId)
	//查找失败
	if err != nil {
		log.Errorf("FindValidAlphaById Error|| %s ", err.Error())
		return nil
	}

	//已经删除或者已经提交
	if alphaModel.IsSubmitted != 0 || alphaModel.DeletedAt.Valid {
		log.Warnf("FindValidAlphaById Error|| this alpha is deleted or submitted, alphaId: %d", alphaId)
		return nil
	}
	alphaViewer := viewer.Alpha{}
	alphaViewer.ID = alphaModel.ID
	alphaViewer.SimulationEnv = json.RawMessage(alphaModel.SimulationEnv)
	alphaViewer.IdeaID = alphaModel.IdeaID
	alphaViewer.Alpha = alphaModel.Alpha
	alphaViewer.TestPeriod = alphaModel.TestPeriod

	alphaViewer.SimulationData = json.RawMessage(alphaModel.SimulationData)

	return &alphaViewer
}

//func FindAlphaListByIdeaIdWithLen(id int64, lenLimit int64) ([]viewer.BrainServiceAlpha, error) {
//	FindAlphaListByIdeaId()
//}

func DeleteAlphaListByIdeaId(id int64) bool {
	isDeleted, err := alphaRepo.DeleteByIdeaId(context.Background(), id)
	if err != nil {
		log.Error(err.Error())
	}
	return isDeleted
}

func DeleteAlphaListByIdeaIdWithTx(tx *gorm.DB, id int64) bool {
	isDeleted, err := alphaRepo.DeleteByIdeaIdWithTx(context.Background(), tx, id)
	if err != nil {
		log.Error(err.Error())
	}
	return isDeleted
}

func UpdateAlphaStatusByID(id int64, alphaStatus int64) error {
	err := alphaRepo.UpdateFields(context.Background(), id, map[string]interface{}{
		"is_submitted": alphaStatus,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
