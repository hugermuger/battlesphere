package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"strings"
	"time"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong UUID format"})
		return
	}

	cue, ok := cfg.importCues[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "cue not found"})
		return
	}

	c.JSON(http.StatusOK, cue)
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
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	if formatName == "dragon_shield" {
		line, err := reader.Read()
		if err != nil {
			importCue.Code = http.StatusBadRequest
			importCue.Message = "Couldn't read csv"
			cfg.importCues[cueID] = importCue
			return
		}

		if len(line) == 1 && strings.HasPrefix(line[0], "sep=") {
			sep := strings.TrimPrefix(line[0], "sep=")
			if len(sep) > 0 {
				reader.Comma = rune(sep[0])
			}
			if _, err := reader.Read(); err != nil {
				importCue.Code = http.StatusBadRequest
				importCue.Message = "Couldn't read header"
				cfg.importCues[cueID] = importCue
				return
			}
		}
	}

	qtx := cfg.db.WithTx(tx)

	folderCache := make(map[string]uuid.UUID)

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
		quantity := 1
		folderName := ""

		if formatName == "dragon_shield" {
			params, quantity, folderName, err = mapDragonShield(line, cfg, context.Background())
			if err != nil {
				line[9] = "English"
				params, quantity, folderName, err = mapDragonShield(line, cfg, context.Background())
				if err != nil {
					importCue.Missing = append(importCue.Missing, line)
					cfg.importCues[cueID] = importCue
					continue
				}
			}
		}

		if folderName != "" {
			if folderID, ok := folderCache[folderName]; ok {
				params.FolderID = uuid.NullUUID{UUID: folderID, Valid: true}
			} else {
				dbFolder, err := qtx.GetFolderByUserAndName(context.Background(), database.GetFolderByUserAndNameParams{
					UserID:     uuid.NullUUID{UUID: userID, Valid: true},
					FolderName: folderName,
				})
				if err != nil {
					dbFolder, err = qtx.CreateFolder(context.Background(), database.CreateFolderParams{
						UserID:     uuid.NullUUID{UUID: userID, Valid: true},
						FolderName: folderName,
					})
					if err != nil {
						importCue.Code = http.StatusInternalServerError
						importCue.Message = "Couldn't create folder in DB"
						cfg.importCues[cueID] = importCue
						return
					}
				}
				folderCache[folderName] = dbFolder.ID
				params.FolderID = uuid.NullUUID{UUID: dbFolder.ID, Valid: true}
			}
		}

		params.UserID = userID

		for i := 0; i < quantity; i++ {
			err = qtx.AddCardToCollections(context.Background(), params)
			if err != nil {
				importCue.Code = http.StatusInternalServerError
				importCue.Message = "Couldn't add card to DB"
				cfg.importCues[cueID] = importCue
				return
			}
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
	importCue.Message = "Finished successful!"
	cfg.importCues[cueID] = importCue
	go cfg.deleteCue(cueID)
}

func (cfg *apiConfig) deleteCue(cueID uuid.UUID) {
	time.Sleep(time.Minute * 10)
	delete(cfg.importCues, cueID)
}
