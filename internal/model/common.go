package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:mili"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
