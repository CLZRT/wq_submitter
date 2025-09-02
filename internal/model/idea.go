package model

import (
	"gorm.io/gorm"
	"time"
)

// Idea represents the idea table
// corresponds to the `idea` table in the database.
type Idea struct {
	ID                int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	IdeaAlphaTemplate string `gorm:"column:idea_alpha_template;type:varchar(255);not null" json:"idea_alpha_template"`
	IdeaTitle         string `gorm:"column:idea_title;type:varchar(255)" json:"idea_title"`
	IdeaDesc          string `gorm:"column:idea_desc;type:varchar(255)" json:"idea_desc"`
	StartIdx          int64  `gorm:"column:start_idx;default:-1" json:"start_idx"`
	EndIdx            int64  `gorm:"column:end_idx;default:-1" json:"end_idx"`
	NextIdx           int64  `gorm:"column:next_idx;default:-1" json:"next_idx"`
	SuccessNum        int64  `gorm:"column:success_num;default:0" json:"success_num"`
	FailNum           int64  `gorm:"column:fail_num;default:0" json:"fail_num"`
	ConcurrencyNum    int64  `gorm:"column:concurrency_num;default:0" json:"concurrency_num"`
	IsFinished        int64  `gorm:"column:is_finished;type:tinyint unsigned;default:0;not null" json:"is_finished"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (i *Idea) TableName() string {
	return "idea"
}
