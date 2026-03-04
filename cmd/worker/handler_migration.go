package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hugermuger/battlesphere/internal/scryfall"
)

func (cfg *apiConfig) handler_initialMigration() error {
	_, err := cfg.db.GetSyncState(context.Background(), "rulings")
	if err == sql.ErrNoRows {
		err = cfg.handler_bulkImportRulings()
		if err != nil {
			return fmt.Errorf("Error bulk importing rulings: %v", err)
		}
	} else if err != nil {
		return err
	} else {
		fmt.Println("Ruling Bulk Data exists!")
	}
	_, err = cfg.db.GetSyncState(context.Background(), "all_cards")
	if err == sql.ErrNoRows {
		err = cfg.handler_bulkImportCards()
		if err != nil {
			return fmt.Errorf("Error bulk importing cards: %v", err)
		}
	} else if err != nil {
		return err
	} else {
		fmt.Println("Card Bulk Data exists!")
	}

	return nil
}

func (cfg *apiConfig) handler_bulkImportCards() error {
	const batchSize = 500

	url, err := scryfall.GetBulkURL(cfg.bulkURL, "all_cards")
	if err != nil {
		return err
	}

	log.Printf("Starting download from %v...", url)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %v", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	var batch []scryfall.CardJSON
	count := 0

	for decoder.More() {
		var card scryfall.CardJSON
		err := decoder.Decode(&card)
		if err != nil {
			log.Printf("Decode error: %v", err)
			continue
		}

		batch = append(batch, card)

		if len(batch) >= batchSize {
			err = cfg.handler_batchImportCards(batch)
			if err != nil {
				return err
			}
			count += len(batch)
			log.Printf("Imported %v cards...", count)
			batch = make([]scryfall.CardJSON, 0, batchSize)
		}
	}

	_, err = decoder.Token()
	if err != nil {
		log.Printf("Warning: failed to read closing bracket: %v", err)
	}

	if len(batch) > 0 {
		err = cfg.handler_batchImportCards(batch)
		if err != nil {
			return err
		}
		count += len(batch)
		log.Printf("Imported %v cards completly", count)
	}

	err = cfg.db.InsertSyncState(context.Background(), "all_cards")
	if err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) handler_batchImportCards(cards []scryfall.CardJSON) error {
	tx, err := cfg.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := cfg.db.WithTx(tx)

	for _, card := range cards {
		err = scryfall.SingleCardImport(context.Background(), qtx, card)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (cfg *apiConfig) handler_bulkImportRulings() error {
	const batchSize = 500

	url, err := scryfall.GetBulkURL(cfg.bulkURL, "rulings")
	if err != nil {
		return err
	}

	log.Printf("Starting download from %v...", url)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %v", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	var batch []scryfall.Rulings
	count := 0

	for decoder.More() {
		var rule scryfall.Rulings
		err := decoder.Decode(&rule)
		if err != nil {
			log.Printf("Decode error: %v", err)
			continue
		}

		batch = append(batch, rule)

		if len(batch) >= batchSize {
			err = cfg.handler_batchImportRulings(batch)
			if err != nil {
				return err
			}
			count += len(batch)
			log.Printf("Imported %v rules...", count)
			batch = make([]scryfall.Rulings, 0, batchSize)
		}
	}

	_, err = decoder.Token()
	if err != nil {
		log.Printf("Warning: failed to read closing bracket: %v", err)
	}

	if len(batch) > 0 {
		err = cfg.handler_batchImportRulings(batch)
		if err != nil {
			return err
		}
		count += len(batch)
		log.Printf("Imported %v rules completly", count)
	}

	err = cfg.db.InsertSyncState(context.Background(), "rulings")
	if err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) handler_batchImportRulings(rules []scryfall.Rulings) error {
	tx, err := cfg.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := cfg.db.WithTx(tx)

	for _, rule := range rules {
		err = scryfall.SingleRuleImport(context.Background(), qtx, rule)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
