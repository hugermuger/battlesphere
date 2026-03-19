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
	db     database.Querier
	dbConn *sql.DB
}

func main() {
	const port = "8080"

	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env file not found: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		db:     dbQueries,
		dbConn: dbConn,
	}

	router := gin.Default()

	router.Use(handlerError())

	router.POST("/users", cfg.handlerUsersCreate)
	router.POST("/login", cfg.handlerLogin)

	router.GET("/cards/search", cfg.handlerSearchCards)
	router.GET("/cards/oracle/:id", cfg.handlerCardsByOracleID)
	router.GET("/cards/:id", cfg.handlerCardByID)
	router.GET("/rulings/:id", cfg.handlerRulings)
	router.Run(":" + port)
}
