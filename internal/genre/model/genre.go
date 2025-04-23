package model

import (
	"gorm.io/gorm"
)

// Genre merepresentasikan entitas Genre pada database
type Genre struct {
	gorm.Model
	Name string `gorm:"size:255;unique;not null;index"`
}
