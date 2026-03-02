package scryfall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type CardJSON struct {
	ID              uuid.UUID    `json:"id"`
	OracleID        *uuid.UUID   `json:"oracle_id"`
	MtgoID          *int64       `json:"mtgo_id"`
	ArenaID         *int64       `json:"arena_id"`
	CardmarketID    *int64       `json:"cardmarket_id"`
	Name            string       `json:"name"`
	FlavorName      *string      `json:"flavor_name"`
	PrintedName     *string      `json:"printed_name"`
	Lang            string       `json:"lang"`
	ReleasedAt      string       `json:"released_at"`
	Layout          string       `json:"layout"`
	ImageUris       *ImageUris   `json:"image_uris"`
	ManaCost        *string      `json:"mana_cost"`
	Cmc             float64      `json:"cmc"`
	TypeLine        string       `json:"type_line"`
	PrintedTypeLine *string      `json:"printed_type_line"`
	OracleText      *string      `json:"oracle_text"`
	PrintedText     *string      `json:"printed_text"`
	Power           *string      `json:"power"`
	Toughness       *string      `json:"toughness"`
	Loyalty         *string      `json:"loyalty"`
	Colors          *[]string    `json:"colors"`
	ColorIdentity   []string     `json:"color_identity"`
	Defense         *string      `json:"defense"`
	Keywords        []string     `json:"keywords"`
	FlavorText      *string      `json:"flavor_text"`
	CardFaces       *[]CardFaces `json:"card_faces"`
	Legalities      Legalities   `json:"legalities"`
	Games           []string     `json:"games"`
	GameChanger     *bool        `json:"game_changer"`
	Finishes        []string     `json:"finishes"`
	Set             string       `json:"set"`
	SetName         string       `json:"set_name"`
	RulingsURI      string       `json:"rulings_uri"`
	CollectorNumber string       `json:"collector_number"`
	Rarity          string       `json:"rarity"`
	Artist          *string      `json:"artist"`
	EdhrecRank      *int         `json:"edhrec_rank"`
	Prices          Prices       `json:"prices"`
}

type CardFaces struct {
	Name            string     `json:"name"`
	PrintedName     *string    `json:"printed_name"`
	ManaCost        string     `json:"mana_cost"`
	Cmc             *float64   `json:"cmc"`
	TypeLine        *string    `json:"type_line"`
	PrintedTypeLine *string    `json:"printed_type_line"`
	OracleText      *string    `json:"oracle_text"`
	PrintedText     *string    `json:"printed_text"`
	Power           *string    `json:"power"`
	Toughness       *string    `json:"toughness"`
	Loyalty         *string    `json:"loyalty"`
	Colors          *[]string  `json:"colors"`
	Defense         *string    `json:"defense"`
	FlavorText      *string    `json:"flavor_text"`
	Artist          *string    `json:"artist"`
	Layout          *string    `json:"layout"`
	ImageUris       *ImageUris `json:"image_uris"`
}

type ImageUris struct {
	Small   *string `json:"small"`
	Normal  *string `json:"normal"`
	Large   *string `json:"large"`
	Png     *string `json:"png"`
	ArtCrop *string `json:"art_crop"`
}

type Legalities struct {
	Standard        *string `json:"standard"`
	Future          *string `json:"future"`
	Historic        *string `json:"historic"`
	Timeless        *string `json:"timeless"`
	Gladiator       *string `json:"gladiator"`
	Pioneer         *string `json:"pioneer"`
	Modern          *string `json:"modern"`
	Legacy          *string `json:"legacy"`
	Pauper          *string `json:"pauper"`
	Vintage         *string `json:"vintage"`
	Penny           *string `json:"penny"`
	Commander       *string `json:"commander"`
	Oathbreaker     *string `json:"oathbreaker"`
	Standardbrawl   *string `json:"standardbrawl"`
	Brawl           *string `json:"brawl"`
	Alchemy         *string `json:"alchemy"`
	Paupercommander *string `json:"paupercommander"`
	Duel            *string `json:"duel"`
	Oldschool       *string `json:"oldschool"`
	Premodern       *string `json:"premodern"`
	Predh           *string `json:"predh"`
}

type Prices struct {
	Usd       *string `json:"usd"`
	UsdFoil   *string `json:"usd_foil"`
	UsdEtched *string `json:"usd_etched"`
	Eur       *string `json:"eur"`
	EurFoil   *string `json:"eur_foil"`
}

func ImportSingleCard(url string) (CardJSON, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CardJSON{}, err
	}

	req.Header.Set("User-Agent", "BattlesphereApp/1.0")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return CardJSON{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return CardJSON{}, fmt.Errorf("scryfall error: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CardJSON{}, err
	}

	var card CardJSON
	err = json.Unmarshal(body, &card)
	if err != nil {
		return CardJSON{}, err
	}

	return card, nil
}
