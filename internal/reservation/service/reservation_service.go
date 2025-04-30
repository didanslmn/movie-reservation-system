package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/repository"
	seat "github.com/didanslmn/movie-reservation-system.git/internal/seat/repository"
	showtime "github.com/didanslmn/movie-reservation-system.git/internal/showtime/repository"
	users "github.com/didanslmn/movie-reservation-system.git/internal/users/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, userID uint, req *request.CreateReservationRequest) (*response.ReservationResponse, error)
	GetReservationByID(ctx context.Context, id uint) (*response.ReservationResponse, error)
	GetAllReservations(ctx context.Context) ([]response.ReservationResponse, error)
	DeleteReservation(ctx context.Context, id uint) error
	ValidateSeatsAvailability(seatIDs []uint, showtimeID uint) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	userRepo        users.UserRepository
	showtimeRepo    showtime.ShowtimeRepository
	seatRepo        seat.SeatRepository
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	userRepo users.UserRepository,
	showtimeRepo showtime.ShowtimeRepository,
	seatRepo seat.SeatRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		userRepo:        userRepo,
		showtimeRepo:    showtimeRepo,
		seatRepo:        seatRepo,
	}
}

func (s *reservationService) CreateReservation(ctx context.Context, userID uint, req *request.CreateReservationRequest) (*response.ReservationResponse, error) {
	if exist, err := s.userRepo.ExistsByID(ctx, userID); err != nil || !exist {
		utils.ErrorLogger.Printf("User not found (ID: %d): %v", userID, err)
		return nil, fmt.Errorf("user not found")
	}
	err := s.ValidateSeatsAvailability(req.SeatIDs, req.ShowtimeID)
	if err != nil {
		return nil, err
	}

	showtime, err := s.showtimeRepo.GetByID(ctx, req.ShowtimeID)
	if err != nil {
		utils.ErrorLogger.Printf("Showtime not found (ID: %d): %v", req.ShowtimeID, err)
		return nil, fmt.Errorf("showtime not found")
	}

	seats, err := s.seatRepo.GetByIDs(ctx, req.SeatIDs)
	if err != nil || len(seats) != len(req.SeatIDs) {
		utils.ErrorLogger.Printf("Some seats not found or invalid: %v", err)
		return nil, fmt.Errorf("one or more seats not found")
	}

	reservation := &model.Reservation{
		UserID:     userID,
		ShowtimeID: req.ShowtimeID,
		Status:     "pending",
		ExpiredAt:  showtime.EndTime,
		Seats:      seats,
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		utils.ErrorLogger.Printf("Failed to create reservation: %v", err)
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}
	for _, seatID := range req.SeatIDs {
		if err := s.seatRepo.UpdateStatus(ctx, seatID, "booked"); err != nil {
			utils.ErrorLogger.Printf("Failed to update seat status (ID: %d): %v", seatID, err)
			return nil, fmt.Errorf("failed to update seat status")
		}
	}

	// Fetch lengkap untuk response
	createdReservation, err := s.reservationRepo.GetByID(ctx, reservation.ID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch created reservation: %v", err)
		return nil, fmt.Errorf("failed to fetch created reservation")
	}
	for _, seatID := range req.SeatIDs {
		err := s.seatRepo.UpdateStatus(ctx, seatID, "booked")
		if err != nil {
			utils.ErrorLogger.Printf("Failed to update seat status (ID: %d): %v", seatID, err)
			return nil, fmt.Errorf("failed to update seat status")
		}
	}

	utils.InfoLogger.Printf("Reservation created: %+v", createdReservation)
	return mapper.ToReservationResponse(createdReservation), nil
}

func (s *reservationService) GetReservationByID(ctx context.Context, id uint) (*response.ReservationResponse, error) {
	res, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get reservation by ID %d: %v", id, err)
		return nil, err
	}
	if res == nil {
		utils.ErrorLogger.Printf("Reservation not found (id: %d)", id)
		return nil, fmt.Errorf("reservation not found")
	}
	utils.InfoLogger.Printf("Fetched reservation by ID %d", id)
	return mapper.ToReservationResponse(res), nil
}

func (s *reservationService) GetAllReservations(ctx context.Context) ([]response.ReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAll(ctx)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get all reservations: %v", err)
		return nil, err
	}
	utils.InfoLogger.Printf("Fetched all reservations, total: %d", len(reservations))
	return mapper.ToReservationResponseList(reservations), nil
}

func (s *reservationService) DeleteReservation(ctx context.Context, id uint) error {
	if err := s.reservationRepo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("Failed to delete reservation (ID: %d): %v", id, err)
		return err
	}
	utils.InfoLogger.Printf("Reservation deleted (ID: %d)", id)
	return nil
}

func (s *reservationService) ValidateSeatsAvailability(seatIDs []uint, showtimeID uint) error {
	type SeatCheckResult struct {
		SeatID uint
		Err    error
	}

	resultChan := make(chan SeatCheckResult, len(seatIDs))
	var wg sync.WaitGroup

	for _, seatID := range seatIDs {
		wg.Add(1)

		go func(seatID uint) {
			defer wg.Done()

			available, err := s.seatRepo.IsSeatAvailable(seatID, showtimeID)
			if err != nil {
				resultChan <- SeatCheckResult{SeatID: seatID, Err: err}
				return
			}
			if !available {
				resultChan <- SeatCheckResult{SeatID: seatID, Err: errors.New("seat not available")}
				return
			}

			resultChan <- SeatCheckResult{SeatID: seatID, Err: nil}
		}(seatID)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		if result.Err != nil {
			return errors.New("seat ID " + fmt.Sprint(result.SeatID) + " is not available or has an error: " + result.Err.Error())
		}
	}

	return nil
}
