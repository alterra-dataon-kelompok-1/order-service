package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	updatedAt time.Time      `gorm:"autoUpdateTime"`
	createdAt time.Time      `gorm:"autoCreateTime"`
	deletedAt gorm.DeletedAt `gorm:"index"`
}
