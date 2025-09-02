package repo

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wq_submitter/internal/model"
)

type AlphaRepo struct {
}

func NewAlphaRepo() *AlphaRepo {
	return &AlphaRepo{}
}

func (alphaRepo *AlphaRepo) Add(_ context.Context, alphaModel *model.Alpha) (int64, error) {
	result := db.Create(alphaModel)
	if result.Error != nil {
		log.Errorf("failed to add alpha: %v", result.Error)
		return -1, result.Error
	}
	return alphaModel.ID, nil

}
func (alphaRepo *AlphaRepo) AddList(_ context.Context, alphaList []*model.Alpha) (int64, error) {
	listLen := len(alphaList)
	if listLen == 0 {
		return 0, nil
	}

	result := db.CreateInBatches(alphaList, 100)
	if result.Error != nil {
		log.Errorf("failed to add alpha list: %v", result.Error)
		return 0, result.Error
	}
	if result.RowsAffected < int64(listLen) {
		err := fmt.Errorf("expected to insert %d records, but only inserted %d", listLen, result.RowsAffected)
		log.Error(err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

func (alphaRepo *AlphaRepo) DeleteById(_ context.Context, alphaId int64) (bool, error) {
	result := db.Delete(&model.Alpha{}, alphaId)
	if result.Error != nil {
		log.Errorf("failed to delete alpha by id %d: %v", alphaId, result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil

}
func (alphaRepo *AlphaRepo) DeleteByIdeaId(_ context.Context, ideaId int64) (bool, error) {
	result := db.Where("idea_id = ?", ideaId).Delete(&model.Alpha{})
	if result.Error != nil {
		log.Errorf("failed to delete alphas by idea_id %d: %v", ideaId, result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil

}
func (alphaRepo *AlphaRepo) DeleteByIdeaIdWithTx(_ context.Context, tx *gorm.DB, ideaId int64) (bool, error) {
	result := tx.Where("idea_id = ?", ideaId).Delete(&model.Alpha{})
	if result.Error != nil {
		log.Errorf("failed to delete alphas by idea_id %d: %v", ideaId, result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil

}

func (alphaRepo *AlphaRepo) Update(_ context.Context, alphaModel *model.Alpha) (bool, error) {
	result := db.Where("id = ?", alphaModel.ID).Updates(*alphaModel)
	if result.Error != nil {
		log.Errorf("failed to update alpha %d: %v", alphaModel.ID, result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil

}
func (alphaRepo *AlphaRepo) UpdateFields(_ context.Context, id int64, fields map[string]interface{}) error {
	result := db.Model(&model.Alpha{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		log.Errorf("UpdateFields for Alpha ID %d Error: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := fmt.Errorf("UpdateFields for Alpha ID %d had no effect or record not found", id)
		log.Error(err.Error())
		// 根据业务逻辑，这可能是一个错误，也可能不是
		return err
	}
	return nil
}
func (alphaRepo *AlphaRepo) FindById(_ context.Context, alphaId int64) (*model.Alpha, error) {
	var alpha model.Alpha
	result := db.First(&alpha, alphaId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err := fmt.Errorf("not found alpha by alphaId %d", alphaId)
			log.Error(err)
			return nil, err
		}
		log.Errorf("failed to find alpha by id %d: %v", alphaId, result.Error)
		return nil, result.Error
	}
	return &alpha, nil

}

func (alphaRepo *AlphaRepo) FindByIdeaId(_ context.Context, ideaId int64) ([]model.Alpha, error) {
	var alphas []model.Alpha
	result := db.Where("idea_id = ?", ideaId).Find(&alphas)
	if result.Error != nil {
		log.Errorf("failed to find alphas by idea_id %d: %v", ideaId, result.Error)
		return nil, result.Error
	}
	return alphas, nil
}
func (alphaRepo *AlphaRepo) FindLast(_ context.Context) (*model.Alpha, error) {
	var alpha model.Alpha
	// 使用Unscoped()方法查询包括已删除的所有记录
	result := db.Unscoped().Last(&alpha)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an application error
		}
		log.Errorf("failed to find last alpha  %s", result.Error.Error())
		return nil, result.Error
	}
	return &alpha, nil

}
