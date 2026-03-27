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
	jwtSecret string
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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		db:        dbQueries,
		dbConn:    dbConn,
		jwtSecret: jwtSecret,
	}

	router := gin.Default()

	router.Use(handlerError())

	router.POST("/users", cfg.handlerUsersCreate)
	router.POST("/login", cfg.handlerLogin)
	router.POST("/refresh", cfg.handlerRefresh)
	router.POST("/revoke", cfg.handlerRevoke)

	router.GET("/cards/search", cfg.handlerSearchCards)
	router.GET("/cards/oracle/:id", cfg.handlerCardsByOracleID)
	router.GET("/cards/:id", cfg.handlerCardByID)
	router.GET("/rulings/:id", cfg.handlerRulings)

	router.POST("/collections/import", cfg.handlerImportCollection)
	router.Run(":" + port)
}
