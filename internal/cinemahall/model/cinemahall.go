package model

import (
	seatModel "github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
	"gorm.io/gorm"
)

type CinemaHall struct {
	gorm.Model
	Name     string           `gorm:"not null;uniqueIndex"`
	Capacity int              `gorm:"not null"`
	Seat     []seatModel.Seat `gorm:"foreignKey:CinemaHallID"`
}
