package main

import (
	//"fmt"
	"log"
	"os"

	"github.com/didanslmn/movie-reservation-system.git/config"

	//"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/didanslmn/movie-reservation-system.git/router"
	//"golang.org/x/crypto/bcrypt"
)

func main() {
	config.LoadEnv()
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// //perintah untuk membuat user dengan role admin secara manual tanpa melalui endpoint
	// newUser := model.User{
	// 	Name:     "admin123",
	// 	Email:    "admin123@yahoo.com",
	// 	Password: "11111111",
	// 	Role:     model.RoleAdmin,
	// }
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	// newUser.Password = string(hashedPassword)
	// // insert data ke database
	// result := db.Create(&newUser)
	// if result.Error != nil {
	// 	fmt.Println("Gagal menambahkan user:", result.Error)
	// 	return
	// }
	// fmt.Println("Berhasil menambahkan user dengan ID:", newUser.ID)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable is not set")
	}
	r := router.SetupRouter(db, jwtSecret)
	if r == nil {
		log.Fatalf("Failed to setup router")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.Printf("Server running at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

/*
{
    "title":"jumbo",
    "description":"film indo",
    "duration":140,
    "release_date":"2025-04-18T14:57:00Z",
    "image_url":"https://i0.wp.com/plopdo.com/wp-content/uploads/2021/11/feature-pic.jpg?fit=537%2C322&ssl=1",
    "rating":4,
    "genre_ids":[2,3]
}*/
