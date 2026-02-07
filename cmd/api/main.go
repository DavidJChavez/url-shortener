package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/davidjchavez/url-shortener/internal/handler"
	"github.com/davidjchavez/url-shortener/internal/repository"
	"github.com/davidjchavez/url-shortener/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	port := getEnv("PORT", "8080")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging database: %v\n", err)
	}

	r := repository.NewURLRepository(db)
	s := service.NewURLService(r, baseURL)
	h := handler.NewURLHandler(s)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://davidjchavez.com", "http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{"Content-Type"},
	}))

	e.POST("/shortener", h.CreateShortURL)
	e.GET("/shortener/stats/:code", h.GetStats)
	e.GET("/shortener/:code", h.Redirect)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
