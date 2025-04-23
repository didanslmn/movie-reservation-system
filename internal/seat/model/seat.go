package model

import "gorm.io/gorm"

type Seat struct {
	gorm.Model
	CinemaHallID uint   `gorm:"not null;index"`
	SeatNumber   string `gorm:"not null"`
	Row          string `gorm:"not null"`
	Status       string `gorm:"default:'available'"`
}
