package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	db         *database.Queries
	dbConn     *sql.DB
	jwtSecret  string
	importCues map[uuid.UUID]importCue
}

type importCue struct {
	Message  string     `json:"message"`
	Code     int        `json:"code"`
	Progress int        `json:"progress"`
	Missing  [][]string `json:"missing"`
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
	importCues := make(map[uuid.UUID]importCue)

	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		db:         dbQueries,
		dbConn:     dbConn,
		jwtSecret:  jwtSecret,
		importCues: importCues,
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

	router.GET("/cue/:id", cfg.handlerImportStatus)
	router.POST("/collections/import", cfg.handlerImportCollection)
	router.Run(":" + port)
}
