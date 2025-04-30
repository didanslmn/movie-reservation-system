# 🎬 Movie API

Movie API adalah RESTful API yang dibuat dengan **Golang**, menggunakan arsitektur clean & modular, untuk menangani sistem reservasi film lengkap dengan fitur authentication, role-based authorization, movie management, genre, dan integrasi dengan database PostgreSQL.

## 💻 Tech Stack

- Language: **Go**
- Framework: **Gin**
- ORM: **GORM**
- Database: **PostgreSQL**
- JWT Auth: **github.com/golang-jwt/jwt/v5**
- Password Hashing: **bcrypt**
- Validation: **binding + validator**
- Migration: **golang-migrate**

## 🚀 Fitur Utama

### ✅ Autentikasi & Otorisasi
- Register dan login dengan hashing password (menggunakan `bcrypt`)
- JWT-based authentication
- Role-based authorization (`admin`, `user`)
- Middleware untuk validasi token dan hak akses

### 🎥 Modul Movie
- CRUD movie (Create, Read, Update, Delete)
- Relasi many-to-many dengan genre
- Filter movie berdasarkan genre
- Validasi input menggunakan DTO dan validator

### 📚 Modul Genre
- CRUD genre
- Relasi ke movie
- Akses terbatas untuk admin (create, update, delete), publik hanya bisa melihat

### 👤 Modul User
- Register user (default role: user)
- Login user dengan JWT
- Lihat & update profil
- Ganti password

### Modul Cinema Hall
- CRUD
- relasi dengan seat

### Modul Seat
- CRUD

### Modul Showtime
- CRUD
- relasi many-to-one dengan cinema hall dan movie

### Modul Reservation
- Create reservation
- reservasi invalid jika showtime sudah selesai atau film sudah selesai
- bisa reservasi lebih dari 1 seats

## 🧠 Pembelajaran & Konsep yang Diterapkan

- ✅ Clean architecture dan modular folder per fitur
- ✅ JWT Authentication + Middleware
- ✅ Role-based authorization dengan custom middleware
- ✅ Relasi many-to-many antara Movie dan Genre menggunakan GORM
- ✅ Custom validation untuk input request (dengan validator dan binding tag)
- ✅ Dependency injection pada repository dan service
- ✅ DTO untuk memisahkan struktur data internal dan external
- ✅ Environment configuration (`.env`)
- ✅ Error handling dan response yang terstruktur
- ✅ Testing API menggunakan Postman
- ✅ Validate reuest

## 📂 Endpoint Utama (Contoh)

### Auth
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Movie
- `GET /api/v1/movies/`
- `GET /api/v1/movies/:id`
- `POST /api/v1/movies/` (admin only)

### Reservaton
- `GET /api/v1/reservations/:id`
- `POST /api/v1/reservatons/` 

## ⚙️ Setup & Jalankan

1. Clone repo ini
2. Buat file `.env`:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=moviedb
   JWT_SECRET=your_jwt_secret
   PORT=8080
   ```
   
