package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	db        *database.Queries
	dbConn    *sql.DB
	platform  string
	jwtSecret string
	polkaKey  string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		db:     dbQueries,
		dbConn: dbConn,
	}

	router := gin.Default()

	router.Use(handlerError())

	router.GET("/cards/search", cfg.handlerSearchCards)
	router.Run(":" + port)
}
