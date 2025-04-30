package model

import (
	"time"

	seatModel "github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
	showtimeModel "github.com/didanslmn/movie-reservation-system.git/internal/showtime/model"
	userModel "github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID     uint
	User       userModel.User `gorm:"foreignKey:UserID"`
	ShowtimeID uint
	Showtime   showtimeModel.Showtime `gorm:"foreignKey:ShowtimeID"`
	Status     string                 `gorm:"type:varchar(20);default:'pending'"`
	ExpiredAt  time.Time

	Seats []seatModel.Seat `gorm:"many2many:reservation_seats;"`
}
type ReservationSeat struct {
	ID            uint `gorm:"primaryKey"`
	ReservationID uint
	SeatID        uint
}
