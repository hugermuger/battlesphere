package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/auth"
	"github.com/hugermuger/battlesphere/internal/database"
)

type Missing struct {
	Missing [][]string `json:"missing"`
}

func (cfg *apiConfig) handlerImportCollection(c *gin.Context) {
	type response struct {
		CueID uuid.UUID `json:"cue_id"`
	}

	formatName := c.Query("format")

	token, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find JWT"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't validate JWT"})
		return
	}

	if formatName != "dragon_shield" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fromat not supported"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read upload"})
		return
	}

	cueID := uuid.New()

	go cfg.importCollection(fileBytes, userID, cueID, formatName)

	c.JSON(http.StatusAccepted, response{
		CueID: cueID,
	})
}

func (cfg *apiConfig) handlerImportStatus(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong UUID format"})
		return
	}

	if cfg.importCues[id].Code == 302 || cfg.importCues[id].Code == 201 {
		c.JSON(cfg.importCues[id].Code, cfg.importCues[id])
	} else {
		c.JSON(cfg.importCues[id].Code, gin.H{"error": cfg.importCues[id].Message})
	}
}

func (cfg *apiConfig) importCollection(f []byte, userID, cueID uuid.UUID, formatName string) {

	importCue := importCue{
		Message:  "In Progress",
		Code:     http.StatusFound,
		Progress: 0,
	}

	cfg.importCues[cueID] = importCue

	tx, err := cfg.dbConn.Begin()
	if err != nil {
		importCue.Code = http.StatusInternalServerError
		importCue.Message = "Couldn't connect to DB"
		cfg.importCues[cueID] = importCue
		return
	}
	defer tx.Rollback()

	reader := csv.NewReader(bytes.NewReader(f))

	if formatName == "dragon_shield" {
		if _, err := reader.Read(); err != nil {
			importCue.Code = http.StatusBadRequest
			importCue.Message = "Couldn't read csv"
			cfg.importCues[cueID] = importCue
			return
		}
		if _, err := reader.Read(); err != nil {
			importCue.Code = http.StatusBadRequest
			importCue.Message = "Couldn't read csv"
			cfg.importCues[cueID] = importCue
			return
		}
	}

	qtx := cfg.db.WithTx(tx)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			importCue.Code = http.StatusBadRequest
			importCue.Message = "Couldn't read csv"
			cfg.importCues[cueID] = importCue
			return
		}

		params := database.AddCardToCollectionsParams{}

		if formatName == "dragon_shield" {
			params, err = mapDragonShield(line, cfg, context.Background())
			if err != nil {
				importCue.Missing = append(importCue.Missing, line)
				cfg.importCues[cueID] = importCue
				continue
			}
		}

		params.UserID = userID

		err = qtx.AddCardToCollections(context.Background(), params)
		if err != nil {
			importCue.Code = http.StatusInternalServerError
			importCue.Message = "Couldn't add card to DB"
			cfg.importCues[cueID] = importCue
			return
		}

		importCue.Progress++
		cfg.importCues[cueID] = importCue
	}

	err = tx.Commit()
	if err != nil {
		importCue.Code = http.StatusInternalServerError
		importCue.Message = "Couldn't add collection to DB"
		cfg.importCues[cueID] = importCue
		return
	}

	importCue.Code = http.StatusCreated
	importCue.Message = "Finished successful"
	cfg.importCues[cueID] = importCue
}
