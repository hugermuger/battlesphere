package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/scryfall"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	dbConn         *sql.DB
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
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

	apiCfg := apiConfig{
		db:     dbQueries,
		dbConn: dbConn,
	}

	card, err := scryfall.ImportSingleCard("https://api.scryfall.com/cards/named?exact=invasion-of-zendikar")
	if err != nil {
		fmt.Println(err)
	}

	err = apiCfg.handler_importSingleCardToDB(card)
	if err != nil {
		fmt.Println(err)
	}
}
