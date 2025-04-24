package router

import (
	cinemahallHandler "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/handler"
	cinemahallRepository "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/repository"
	cinemahallRouter "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/router"
	cinemahallService "github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/service"
	genreHandler "github.com/didanslmn/movie-reservation-system.git/internal/genre/handler"
	genreRepository "github.com/didanslmn/movie-reservation-system.git/internal/genre/repository"
	genreRouter "github.com/didanslmn/movie-reservation-system.git/internal/genre/router" // Perbaikan di sini
	genreService "github.com/didanslmn/movie-reservation-system.git/internal/genre/service"
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	movieHandler "github.com/didanslmn/movie-reservation-system.git/internal/movie/handler"
	movieRepository "github.com/didanslmn/movie-reservation-system.git/internal/movie/repository"
	movieRouter "github.com/didanslmn/movie-reservation-system.git/internal/movie/router" // Perbaikan di sini
	movieService "github.com/didanslmn/movie-reservation-system.git/internal/movie/service"
	seatHandler "github.com/didanslmn/movie-reservation-system.git/internal/seat/handler"
	seatRepository "github.com/didanslmn/movie-reservation-system.git/internal/seat/repository"
	seatRouter "github.com/didanslmn/movie-reservation-system.git/internal/seat/router"
	seatService "github.com/didanslmn/movie-reservation-system.git/internal/seat/service"
	userHandler "github.com/didanslmn/movie-reservation-system.git/internal/users/handler"
	userRepository "github.com/didanslmn/movie-reservation-system.git/internal/users/repository"
	userRouter "github.com/didanslmn/movie-reservation-system.git/internal/users/router" // Perbaikan di sini
	userService "github.com/didanslmn/movie-reservation-system.git/internal/users/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// === User Setup ===
	userRepo := userRepository.NewUserRepository(db)
	userSvc := userService.NewUserService(userRepo, jwtSecret)
	userHdl := userHandler.NewUserHandler(userSvc)
	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	// === Genre Setup ===
	genreRepo := genreRepository.NewRepositoryGenre(db)
	genreSvc := genreService.NewGenreService(genreRepo)
	genreHdl := genreHandler.NewGenreHandler(genreSvc)

	// === Movie Setup ===
	movieRepo := movieRepository.NewRepositoryMovie(db)
	movieSvc := movieService.NewMovieService(movieRepo, genreRepo)
	movieHdl := movieHandler.NewMovieHandler(movieSvc)

	// === Cinema Hall Setup ===
	cinemahallRepo := cinemahallRepository.NewCinemaHallRepository(db)
	cinemahallSvc := cinemahallService.NewCinemaHallService(cinemahallRepo)
	cinemahallHdl := cinemahallHandler.NewCinemaHallHandler(cinemahallSvc)

	// === Seat Setup ===
	seatRepo := seatRepository.NewSeatRepository(db)
	seatSvc := seatService.NewSeatService(seatRepo)
	seatHdl := seatHandler.NewSeatHandler(seatSvc)

	// Register routes
	api := r.Group("/api/v1")
	userRouter.AuthRoutes(api, userHdl)
	Protected := api.Group("/")
	Protected.Use(authMiddleware)
	userRouter.UserRoutes(Protected, userHdl, jwtSecret)
	genreRouter.GenreRoutes(Protected, genreHdl, jwtSecret)
	movieRouter.MovieRoutes(Protected, movieHdl, jwtSecret)
	cinemahallRouter.SeatRouts(Protected, cinemahallHdl, jwtSecret)
	seatRouter.SeatRouts(Protected, seatHdl, jwtSecret)

	return r
}
