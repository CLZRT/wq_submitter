package repo

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wq_submitter/internal/model"
)

type IdeaRepo struct {
}

func NewIdeaRepo() *IdeaRepo {
	return &IdeaRepo{}
}

func (ideaRepo *IdeaRepo) Add(_ context.Context, idea *model.Idea) (int64, error) {

	result := db.Create(idea)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Errorf("Add Idea Error: %v", result.Error)
		return -1, result.Error
	}
	return idea.ID, nil
}

func (ideaRepo *IdeaRepo) DeleteById(_ context.Context, id int64) (int64, error) {

	idea := model.Idea{}
	result := db.Where(id).Delete(&idea)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Errorf("Delete Idea Error: %v", result.Error)
		return -1, result.Error

	}
	return idea.ID, nil
}

func (ideaRepo *IdeaRepo) DeleteByIdWithTx(_ context.Context, tx *gorm.DB, id int64) (int64, error) { // 更改函数签名
	idea := model.Idea{}

	result := tx.Where(id).Delete(&idea)
	if result.Error != nil {
		log.Errorf("Delete Idea Error in transaction: %v", result.Error)
		return -1, result.Error
	}

	if result.RowsAffected == 0 {
		log.Warnf("No Idea found with ID %d to delete in transaction", id)
	}

	return id, nil // 返回成功删除的ID
}

func (ideaRepo *IdeaRepo) Update(_ context.Context, idea *model.Idea) (int64, error) {

	result := db.Where("id = ?", idea.ID).Updates(*idea)
	if result.Error != nil || result.RowsAffected == 0 {
		return -1, result.Error
	}
	return idea.ID, nil
}

// UpdateFields updates specific fields of an idea using a map.
// This method allows updating fields to their zero values.
func (ideaRepo *IdeaRepo) UpdateFields(_ context.Context, id int64, fields map[string]interface{}) error {
	result := db.Model(&model.Idea{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		log.Errorf("UpdateFields for Idea ID %d Error: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := fmt.Errorf("UpdateFields for Idea ID %d had no effect or record not found", id)
		log.Error(err.Error())
		// 根据业务逻辑，这可能是一个错误，也可能不是
		return err
	}
	return nil
}

func (ideaRepo *IdeaRepo) UpdateList(_ context.Context, ideas []*model.Idea) (int64, error) {
	var totalRowAffected int64

	tx := db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, idea := range ideas {
		result := tx.Model(&model.Idea{}).Where("id = ?", idea.ID).Updates(*idea)
		if result.Error != nil {
			tx.Rollback()
			log.Errorf("Update Idea ID %d Error: %v", idea.ID, result.Error)
			return -1, result.Error
		}

		totalRowAffected += result.RowsAffected

	}

	if err := tx.Commit().Error; err != nil {
		log.Errorf("Commite transaction error: %v", err)
		return totalRowAffected, err
	}

	return totalRowAffected, nil

}

func (ideaRepo *IdeaRepo) FindAll(_ context.Context) ([]model.Idea, error) {

	var ideas []model.Idea
	result := db.Find(&ideas)
	if result.Error != nil {
		log.Errorf("GetAll Idea Error: %v", result.Error)
		return nil, result.Error
	}
	return ideas, nil

}
func (ideaRepo *IdeaRepo) IsValidById(_ context.Context, id int64) (bool, *model.Idea) {

	var idea model.Idea
	result := db.Where("id = ? AND next_idx <= end_idx AND deleted_at IS NULL", id).First(&idea)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Errorf("IsValidById Idea Error: %v", result.Error)
		return false, nil
	}
	return true, &idea

}
func (ideaRepo *IdeaRepo) IsRunningById(_ context.Context, id int64) (bool, *model.Idea) {

	var idea model.Idea
	result := db.Where("id = ? AND next_idx <= end_idx AND deleted_at IS NULL AND concurrency_num > ? ", id, 0).First(&idea)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Errorf("IsRunningById Idea Error: %v", result.Error)
		return false, nil
	}
	return true, &idea

}

func (ideaRepo *IdeaRepo) FindValid(_ context.Context) ([]model.Idea, error) {

	var ideas []model.Idea
	result := db.Where("next_idx <= end_idx AND deleted_at IS NULL").Find(&ideas)
	if result.Error != nil {
		log.Errorf("FindValid Idea Error: %v", result.Error)
		return nil, result.Error
	}
	return ideas, nil
}

func (ideaRepo *IdeaRepo) FindNeedRun(_ context.Context) ([]model.Idea, error) {
	var ideas []model.Idea
	result := db.Where("next_idx <= end_idx AND deleted_at IS NULL AND concurrency_num > ? ", 0).Find(&ideas)
	if result.Error != nil {
		log.Errorf("FindNeedRun Idea Error: %v", result.Error)
		return nil, result.Error
	}
	return ideas, nil
}

func (ideaRepo *IdeaRepo) FindById(_ context.Context, id int64) (*model.Idea, error) {
	var idea model.Idea
	result := db.Where("id = ? ", id).Find(&idea)
	if result.Error != nil {
		log.Errorf("Find Idea By Id Error: %v", result.Error)
		return nil, result.Error
	}
	return &idea, nil
}
