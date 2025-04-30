package service

import (
	"context"
	"errors"
	"fmt"

	hallRepository "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/repository"
	movieRepository "github.com/didanslmn/movie-reservation-system.git/internal/movie/repository"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
)

type ShowtimeService interface {
	CreateShowtime(ctx context.Context, req request.CreateShowtimeRequest) (*response.ShowtimeResponse, error)
	GetShowtimeByID(ctx context.Context, id uint) (*response.ShowtimeResponse, error)
	GetAllShowtimes(ctx context.Context) ([]response.ShowtimeResponse, error)
	UpdateShowtime(ctx context.Context, id uint, req request.UpdateShowtimeRequest) (*response.ShowtimeResponse, error)
	DeleteShowtime(ctx context.Context, id uint) error
}

type showtimeService struct {
	showtimeRepo repository.ShowtimeRepository
	movieRepo    movieRepository.MovieRepository
	hallRepo     hallRepository.CinemaHallRepository
}

func NewShowtimeService(showtimeRepo repository.ShowtimeRepository, movieRepo movieRepository.MovieRepository, hallRepo hallRepository.CinemaHallRepository) ShowtimeService {
	return &showtimeService{
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		hallRepo:     hallRepo,
	}
}

// Create Showtime
func (s *showtimeService) CreateShowtime(ctx context.Context, req request.CreateShowtimeRequest) (*response.ShowtimeResponse, error) {
	// Validate Movie ID
	exist, err := s.movieRepo.ExistsByID(ctx, req.MovieID)
	if err != nil || !exist {
		utils.ErrorLogger.Printf("Invalid movie ID: %d", req.MovieID)
		return nil, errors.New("invalid movie ID")
	}

	// Validate Cinema Hall ID
	exist, err = s.hallRepo.ExistsByID(ctx, req.CinemaHallID)
	if err != nil || !exist {
		utils.ErrorLogger.Printf("Invalid cinema hall ID: %d", req.CinemaHallID)
		return nil, errors.New("invalid cinema hall ID")
	}

	showtime := model.Showtime{
		MovieID:      req.MovieID,
		CinemaHallID: req.CinemaHallID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	}

	if err := s.showtimeRepo.Create(ctx, &showtime); err != nil {
		utils.ErrorLogger.Printf("Failed to create showtime: %v", err)
		return nil, fmt.Errorf("failed to create showtime: %w", err)
	}

	utils.InfoLogger.Printf("Successfully created showtime (ID: %d)", showtime.ID)
	return mapper.ToShowtimeResponse(&showtime), nil
}

// Get Showtime by ID
func (s *showtimeService) GetShowtimeByID(ctx context.Context, id uint) (*response.ShowtimeResponse, error) {
	showtime, err := s.showtimeRepo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Showtime not found (ID: %d): %v", id, err)
		return nil, fmt.Errorf("showtime not found: %w", err)
	}
	return mapper.ToShowtimeResponse(showtime), nil
}

// Get All Showtimes
func (s *showtimeService) GetAllShowtimes(ctx context.Context) ([]response.ShowtimeResponse, error) {
	showtimes, err := s.showtimeRepo.GetAll(ctx)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get showtimes: %v", err)
		return nil, fmt.Errorf("failed to get showtimes: %w", err)
	}

	var result []response.ShowtimeResponse
	for _, s := range showtimes {
		result = append(result, *mapper.ToShowtimeResponse(&s))
	}
	return result, nil
}

// Update Showtime
func (s *showtimeService) UpdateShowtime(ctx context.Context, id uint, req request.UpdateShowtimeRequest) (*response.ShowtimeResponse, error) {
	showtime, err := s.showtimeRepo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Showtime not found (ID: %d): %v", id, err)
		return nil, fmt.Errorf("showtime not found: %w", err)
	}

	if req.MovieID != 0 {
		exist, err := s.movieRepo.ExistsByID(ctx, req.MovieID)
		if err != nil || !exist {
			utils.ErrorLogger.Printf("Invalid movie ID: %d", req.MovieID)
			return nil, errors.New("invalid movie ID")
		}
		showtime.MovieID = req.MovieID
	}

	if req.CinemaHallID != 0 {
		exist, err := s.hallRepo.ExistsByID(ctx, req.CinemaHallID)
		if err != nil || !exist {
			utils.ErrorLogger.Printf("Invalid cinema hall ID: %d", req.CinemaHallID)
			return nil, errors.New("invalid cinema hall ID")
		}
		showtime.CinemaHallID = req.CinemaHallID
	}

	if !req.StartTime.IsZero() {
		showtime.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		showtime.EndTime = req.EndTime
	}

	if err := s.showtimeRepo.Update(ctx, showtime); err != nil {
		utils.ErrorLogger.Printf("Failed to update showtime (ID: %d): %v", id, err)
		return nil, fmt.Errorf("failed to update showtime: %w", err)
	}

	utils.InfoLogger.Printf("Successfully updated showtime (ID: %d)", showtime.ID)
	return mapper.ToShowtimeResponse(showtime), nil
}

// Delete Showtime
func (s *showtimeService) DeleteShowtime(ctx context.Context, id uint) error {
	if err := s.showtimeRepo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("Failed to delete showtime (ID: %d): %v", id, err)
		return fmt.Errorf("failed to delete showtime: %w", err)
	}
	utils.InfoLogger.Printf("Successfully deleted showtime (ID: %d)", id)
	return nil
}
