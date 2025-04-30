package model

import (
	"time"

	cinemaHallModel "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/model"
	movieModel "github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
	"gorm.io/gorm"
)

type Showtime struct {
	gorm.Model
	MovieID      uint                       `gorm:"not null"`
	Movie        movieModel.Movie           `gorm:"foreignKey:MovieID"`
	CinemaHallID uint                       `gorm:"not null"`
	CinemaHall   cinemaHallModel.CinemaHall `gorm:"foreignKey:CinemaHallID"`
	StartTime    time.Time                  `gorm:"not null"`
	EndTime      time.Time                  `gorm:"not null"`
}
