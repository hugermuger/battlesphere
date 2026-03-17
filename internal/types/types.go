package types

import (
	"fmt"
	"time"

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

type Rulings struct {
	OracleID    uuid.UUID `json:"oracle_id"`
	Source      *string   `json:"source"`
	PublishedAt string    `json:"published_at"`
	Comment     string    `json:"comment"`
}

type URLS struct {
	Data []struct {
		Type        string    `json:"type"`
		UpdatedAt   time.Time `json:"updated_at"`
		DownloadURI string    `json:"download_uri"`
	} `json:"data"`
}

type CardResponseSearchByName struct {
	OracleID *uuid.UUID `json:"oracle_id"`
	Name     string     `json:"name"`
	Layout   string     `json:"layout"`
	ManaCost *string    `json:"mana_cost"`
	TypeLine string     `json:"type_line"`
}

func (i CardResponseSearchByName) Title() string {
	if i.ManaCost != nil && *i.ManaCost != "" {
		return fmt.Sprintf("%v - %v", i.Name, *i.ManaCost)
	} else {
		return i.Name
	}
}
func (i CardResponseSearchByName) Description() string {
	if i.Layout != "normal" {
		return fmt.Sprintf("%v (%v)", i.TypeLine, i.Layout)
	} else {
		return i.TypeLine
	}
}
func (i CardResponseSearchByName) FilterValue() string { return i.Name }

type CardResponseSearchByOracleID struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	FlavorName      *string   `json:"flavor_name"`
	ReleasedAt      time.Time `json:"released_at"`
	Set             string    `json:"set"`
	SetName         string    `json:"set_name"`
	CollectorNumber string    `json:"collector_number"`
}

func (i CardResponseSearchByOracleID) Title() string {
	if i.FlavorName != nil {
		return fmt.Sprintf("%v (%v)", *i.FlavorName, i.Name)
	} else {
		return i.Name
	}
}
func (i CardResponseSearchByOracleID) Description() string {
	return fmt.Sprintf("%v %v (%v)", i.SetName, i.CollectorNumber, i.ReleasedAt.Format("2006"))
}
func (i CardResponseSearchByOracleID) FilterValue() string { return i.Name }

type ResponseByOracleID struct {
	Name          string                 `json:"name"`
	Layout        string                 `json:"layout"`
	Cmc           float64                `json:"cmc"`
	Colors        *[]string              `json:"colors"`
	ColorIdentity []string               `json:"color_identity"`
	ManaCost      *string                `json:"mana_cost"`
	TypeLine      string                 `json:"type_line"`
	OracleText    *string                `json:"oracle_text"`
	PrintedText   *string                `json:"printed_text"`
	Power         *string                `json:"power"`
	Toughness     *string                `json:"toughness"`
	Loyalty       *string                `json:"loyalty"`
	Defense       *string                `json:"defense"`
	CardFaces     *[]CardFacesByOracleID `json:"card_faces"`
	Multifaced    bool                   `json:"multifaced"`
	Legalities    Legalities             `json:"legalities"`
	Rulings       []ResponseRulings      `json:"rulings"`
	GameChanger   *bool                  `json:"game_changer"`
	EdhrecRank    *int32                 `json:"edhrec_rank"`
}

type CardFacesByOracleID struct {
	Name        string    `json:"name"`
	ManaCost    string    `json:"mana_cost"`
	Cmc         *float64  `json:"cmc"`
	Colors      *[]string `json:"colors"`
	TypeLine    *string   `json:"type_line"`
	OracleText  *string   `json:"oracle_text"`
	PrintedText *string   `json:"printed_text"`
	Power       *string   `json:"power"`
	Toughness   *string   `json:"toughness"`
	Loyalty     *string   `json:"loyalty"`
	Defense     *string   `json:"defense"`
}

type ResponseRulings struct {
	OracleID    uuid.UUID `json:"oracle_id"`
	Source      *string   `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Comment     string    `json:"comment"`
}

type ResponseCard struct {
	ID              uuid.UUID            `json:"id"`
	OracleID        *uuid.UUID           `json:"oracle_id"`
	Name            string               `json:"name"`
	FlavorName      *string              `json:"flavor_name"`
	Lang            string               `json:"lang"`
	ReleasedAt      time.Time            `json:"released_at"`
	Layout          string               `json:"layout"`
	ImageUris       ResponseImages       `json:"image_uris"`
	ManaCost        *string              `json:"mana_cost"`
	Cmc             float64              `json:"cmc"`
	TypeLine        string               `json:"type_line"`
	OracleText      *string              `json:"oracle_text"`
	PrintedText     *string              `json:"printed_text"`
	Power           *string              `json:"power"`
	Toughness       *string              `json:"toughness"`
	Loyalty         *string              `json:"loyalty"`
	Colors          *[]string            `json:"colors"`
	ColorIdentity   []string             `json:"color_identity"`
	Defense         *string              `json:"defense"`
	Keywords        []string             `json:"keywords"`
	FlavorText      *string              `json:"flavor_text"`
	CardFaces       *[]ResponseCardFaces `json:"card_faces"`
	Legalities      Legalities           `json:"legalities"`
	GameChanger     *bool                `json:"game_changer"`
	Finishes        []string             `json:"finishes"`
	Set             string               `json:"set"`
	SetName         string               `json:"set_name"`
	CollectorNumber string               `json:"collector_number"`
	Rarity          string               `json:"rarity"`
	Artist          *string              `json:"artist"`
	EdhrecRank      *int32               `json:"edhrec_rank"`
	Multifaced      bool                 `json:"multifaced"`
}

type ResponseImages struct {
	ImageNormal *string `json:"image_normal"`
	ImagePNG    *string `json:"image_png"`
	ImageLarge  *string `json:"image_large"`
	ImageSmall  *string `json:"image_small"`
	ImageCrop   *string `json:"image_crop"`
}

type ResponseCardFaces struct {
	Name        string         `json:"name"`
	ManaCost    string         `json:"mana_cost"`
	Cmc         *float64       `json:"cmc"`
	TypeLine    *string        `json:"type_line"`
	OracleText  *string        `json:"oracle_text"`
	PrintedText *string        `json:"printed_text"`
	Power       *string        `json:"power"`
	Toughness   *string        `json:"toughness"`
	Loyalty     *string        `json:"loyalty"`
	Colors      *[]string      `json:"colors"`
	Defense     *string        `json:"defense"`
	FlavorText  *string        `json:"flavor_text"`
	Artist      *string        `json:"artist"`
	Layout      *string        `json:"layout"`
	ImageUris   ResponseImages `json:"image_uris"`
}
