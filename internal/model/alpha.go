package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// Alpha corresponds to the `alpha` table in the database.
type Alpha struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	SimulationEnv  datatypes.JSON `gorm:"column:simulation_env;type:json;not null" json:"simulation_env"`
	Alpha          string         `gorm:"column:alpha;type:longtext;not null" json:"alpha"`
	IdeaID         int64          `gorm:"column:idea_id" json:"idea_id"`
	SimulationData datatypes.JSON `gorm:"column:simulation_data;type:json;not null" json:"simulation_data"`
	TestPeriod     string         `gorm:"column:test_period" json:"test_period"`
	IsSubmitted    int64          `gorm:"column:is_submitted;type:tinyint unsigned;default:0;not null" json:"is_submitted"`
	IsDeleted      int64          `gorm:"column:is_deleted;type:tinyint unsigned;default:0;not null" json:"is_deleted"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (a *Alpha) TableName() string {
	return "alpha"
}
