package svc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"wq_submitter/internal/model"
	"wq_submitter/internal/repo"
	"wq_submitter/internal/viewer"
)

var ideaRepo *repo.IdeaRepo
var mutex sync.Mutex

func init() {
	ideaRepo = repo.NewIdeaRepo()
	mutex = sync.Mutex{}
}

func UploadIdea(ideaViewer *viewer.Idea) (ideaId int64, err error) {

	ideaModel := model.Idea{
		ID:                0,
		IdeaAlphaTemplate: ideaViewer.IdeaAlphaTemplate,
		IdeaTitle:         ideaViewer.IdeaTitle,
		IdeaDesc:          ideaViewer.IdeaDesc,
		StartIdx:          ideaViewer.StartIdx,
		EndIdx:            ideaViewer.EndIdx,
		NextIdx:           ideaViewer.NextIdx,
		ConcurrencyNum:    ideaViewer.ConcurrencyNum,
		IsFinished:        0,
	}
	ideaId, err = ideaRepo.Add(context.Background(), &ideaModel)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	return ideaId, nil
}

func UpdateIdea(ideaViewer *viewer.Idea) (int64, error) {

	ideaModel := model.Idea{}
	ideaModel.ID = ideaViewer.ID
	ideaModel.ConcurrencyNum = ideaViewer.ConcurrencyNum

	ideaId, err := ideaRepo.Update(context.Background(), &ideaModel)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}
	return ideaId, nil

}
func UpdateIdeaNextIdx(ideaId int64, nextIdx int64) error {

	err := ideaRepo.UpdateFields(context.Background(), ideaId, map[string]interface{}{
		"next_idx": nextIdx,
	})
	if err != nil {
		log.Errorf("Error in UpdateIdeaNextIdx: %s", err.Error())
		return err
	}
	return nil
}
func UpdateIdeaIsFinished(ideaId int64, isFinished int64) error {

	err := ideaRepo.UpdateFields(context.Background(), ideaId, map[string]interface{}{
		"is_finished": isFinished,
	})
	if err != nil {
		log.Errorf("Error in UpdateIdeaIsFinished: %s", err.Error())
		return err
	}
	return nil
}
func AddIdeaSuccessNum(ideaId int64, successNum int64) error {
	idea, err := ideaRepo.FindById(context.Background(), ideaId)
	if err != nil {
		log.Errorf("Error in AddIdeaSuccessNum: %s", err.Error())
		return err
	}

	err = ideaRepo.UpdateFields(context.Background(), ideaId, map[string]interface{}{
		"success_num": idea.SuccessNum + successNum,
	})
	if err != nil {
		log.Errorf("Error in AddIdeaSuccessNum: %s", err.Error())
		return err
	}
	return nil
}
func AddIdeaFailNum(ideaId int64, successNum int64) error {
	idea, err := ideaRepo.FindById(context.Background(), ideaId)
	if err != nil {
		log.Errorf("Error in AddIdeaFailNum: %s", err.Error())
		return err
	}
	err = ideaRepo.UpdateFields(context.Background(), ideaId, map[string]interface{}{
		"fail_num": idea.FailNum + successNum,
	})
	if err != nil {
		log.Errorf("Error in UpdateIdeaFailedNum: %s", err.Error())
		return err
	}
	return nil
}
func IsLegalConcurrency(ideaId int64, concurrency int64) bool {
	runningIdeas, err := ideaRepo.FindNeedRun(context.Background())
	if err != nil {
		log.Error(err.Error())
		return false
	}
	var currencyNum int64
	mutex.Lock()
	defer mutex.Unlock()
	for _, idea := range runningIdeas {
		if idea.ID == ideaId {
			continue
		}
		currencyNum += idea.ConcurrencyNum
	}
	if currencyNum+concurrency > conf.AppConfig.Concurrency {
		return false
	}
	return true
}
func GetIdeaById(id int64) viewer.Idea {
	ideaModel, err := ideaRepo.FindById(context.Background(), id)
	if err != nil {
		log.Error(err.Error())
		return viewer.Idea{}
	}
	ideaViewer := viewer.Idea{}
	ideaViewer.ID = ideaModel.ID
	ideaViewer.IdeaAlphaTemplate = ideaModel.IdeaAlphaTemplate
	ideaViewer.IdeaTitle = ideaModel.IdeaTitle
	ideaViewer.IdeaDesc = ideaModel.IdeaDesc
	ideaViewer.StartIdx = ideaModel.StartIdx
	ideaViewer.EndIdx = ideaModel.EndIdx
	ideaViewer.NextIdx = ideaModel.NextIdx
	ideaViewer.ConcurrencyNum = ideaModel.ConcurrencyNum
	ideaViewer.IsFinished = ideaModel.IsFinished
	return ideaViewer
}

func GetAllIdea() ([]viewer.Idea, error) {
	ideas, err := ideaRepo.FindAll(context.Background())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ideaViewers := make([]viewer.Idea, 0)

	for _, idea := range ideas {
		ideaViewer := viewer.Idea{}
		ideaViewer.ID = idea.ID
		ideaViewer.IdeaAlphaTemplate = idea.IdeaAlphaTemplate
		ideaViewer.IdeaTitle = idea.IdeaTitle
		ideaViewer.IdeaDesc = idea.IdeaDesc
		ideaViewer.StartIdx = idea.StartIdx
		ideaViewer.EndIdx = idea.EndIdx
		ideaViewer.NextIdx = idea.NextIdx
		ideaViewer.SuccessNum = idea.SuccessNum
		ideaViewer.FailNum = idea.FailNum
		ideaViewer.ConcurrencyNum = idea.ConcurrencyNum
		ideaViewer.IsFinished = idea.IsFinished
		ideaViewers = append(ideaViewers, ideaViewer)
	}
	return ideaViewers, nil

}

func GetUnfinishedIdea() ([]viewer.Idea, error) {
	ideas, err := ideaRepo.FindValid(context.Background())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ideaViewers := make([]viewer.Idea, 0)
	for _, idea := range ideas {
		ideaViewer := viewer.Idea{}
		ideaViewer.ID = idea.ID
		ideaViewer.IdeaAlphaTemplate = idea.IdeaAlphaTemplate
		ideaViewer.IdeaTitle = idea.IdeaTitle
		ideaViewer.IdeaDesc = idea.IdeaDesc
		ideaViewer.StartIdx = idea.StartIdx
		ideaViewer.EndIdx = idea.EndIdx
		ideaViewer.NextIdx = idea.NextIdx
		ideaViewer.SuccessNum = idea.SuccessNum
		ideaViewer.FailNum = idea.FailNum
		ideaViewer.ConcurrencyNum = idea.ConcurrencyNum
		ideaViewer.IsFinished = idea.IsFinished
		ideaViewers = append(ideaViewers, ideaViewer)
	}
	return ideaViewers, nil

}

func GetNeedRunIdea() ([]viewer.Idea, error) {
	ideas, err := ideaRepo.FindNeedRun(context.Background())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ideaViewers := make([]viewer.Idea, 0)
	for _, idea := range ideas {
		ideaViewer := viewer.Idea{}
		ideaViewer.ID = idea.ID
		ideaViewer.IdeaAlphaTemplate = idea.IdeaAlphaTemplate
		ideaViewer.IdeaTitle = idea.IdeaTitle
		ideaViewer.IdeaDesc = idea.IdeaDesc
		ideaViewer.StartIdx = idea.StartIdx
		ideaViewer.EndIdx = idea.EndIdx
		ideaViewer.NextIdx = idea.NextIdx
		ideaViewer.SuccessNum = idea.SuccessNum
		ideaViewer.FailNum = idea.FailNum
		ideaViewer.ConcurrencyNum = idea.ConcurrencyNum
		ideaViewer.IsFinished = idea.IsFinished
		ideaViewers = append(ideaViewers, ideaViewer)
	}
	return ideaViewers, nil

}

func DeleteIdea(id int64) (int64, error) {

	byId, err := ideaRepo.DeleteById(context.Background(), id)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return byId, nil
}

func DeleteIdeaWithTx(id int64) (viewer.Idea, error) {

	db := repo.GetDbCli()
	tx := db.Begin()
	if tx.Error != nil {
		log.Error(tx.Error.Error())
		return viewer.Idea{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	byId, err := ideaRepo.DeleteByIdWithTx(context.Background(), tx, id)
	if err != nil {
		log.Error("Delete Idea Error")
		return viewer.Idea{}, err
	}
	ideaModel, err := ideaRepo.FindById(context.Background(), byId)
	if err != nil {
		log.Error("Find Deleted Idea Error")
		return viewer.Idea{}, err
	}
	ideaViewer := viewer.Idea{}
	ideaViewer.ID = ideaModel.ID
	ideaViewer.IdeaAlphaTemplate = ideaModel.IdeaAlphaTemplate
	ideaViewer.IdeaTitle = ideaModel.IdeaTitle
	ideaViewer.IdeaDesc = ideaModel.IdeaDesc
	ideaViewer.StartIdx = ideaModel.StartIdx
	ideaViewer.EndIdx = ideaModel.EndIdx
	ideaViewer.NextIdx = ideaModel.NextIdx
	ideaViewer.SuccessNum = ideaModel.SuccessNum
	ideaViewer.FailNum = ideaModel.FailNum
	ideaViewer.ConcurrencyNum = ideaModel.ConcurrencyNum
	ideaViewer.IsFinished = ideaModel.IsFinished

	if !DeleteAlphaListByIdeaIdWithTx(tx, id) {
		return viewer.Idea{}, err
	}

	return ideaViewer, nil
}
