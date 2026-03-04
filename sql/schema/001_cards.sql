-- +goose Up
CREATE TABLE cards (
    id                UUID PRIMARY KEY,
    arena_id          BIGINT,
    mtgo_id           BIGINT,
    cardmarket_id     BIGINT,
    oracle_id         UUID,
    release_date      DATE NOT NULL,

    lang              TEXT NOT NULL,
    layout            TEXT NOT NULL,
    edhrec_rank       INTEGER,
    game_changer      BOOLEAN,
    multifaced        BOOLEAN NOT NULL,

    cmc               DOUBLE PRECISION NOT NULL,
    color_identity    TEXT[],
    colors            TEXT[],
    defense           TEXT,
    keywords          TEXT[] NOT NULL,
    loyalty           TEXT,
    mana_cost         TEXT,
    name              TEXT NOT NULL,
    oracle_text       TEXT,
    power             TEXT,
    toughness         TEXT,
    type_line         TEXT NOT NULL,

    artist            TEXT,
    collector_number  TEXT NOT NULL,
    finishes          TEXT[] NOT NULL,
    flavor_name       TEXT,
    flavor_text       TEXT,
    games             TEXT[] NOT NULL,
    image             TEXT,
    image_png         TEXT,
    image_large       TEXT,
    image_small       TEXT,
    image_crop        TEXT,
    price_usd         DOUBLE PRECISION,
    price_eur         DOUBLE PRECISION,
    price_foil_usd    DOUBLE PRECISION,
    price_foil_eur    DOUBLE PRECISION,
    price_etched_usd  DOUBLE PRECISION,
    printed_name      TEXT,
    printed_text      TEXT,
    printed_type_line TEXT,
    rarity            TEXT NOT NULL,
    set_name          TEXT NOT NULL,
    set_code          TEXT NOT NULL,

    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL
);

CREATE TABLE card_faces (
    card_id           UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    face_index        INTEGER NOT NULL,

    artist            TEXT,
    cmc               DOUBLE PRECISION,
    colors            TEXT[],
    defense           TEXT,
    flavor_text       TEXT,
    image             TEXT,
    image_png         TEXT,
    image_large       TEXT,
    image_small       TEXT,
    image_crop        TEXT,
    layout            TEXT,
    loyalty           TEXT,
    mana_cost         TEXT NOT NULL,
    name              TEXT NOT NULL,
    oracle_text       TEXT,
    power             TEXT,
    printed_name      TEXT,
    printed_text      TEXT,
    printed_type_line TEXT,
    toughness         TEXT,
    type_line         TEXT,
    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL,

    PRIMARY KEY (card_id, face_index)
);

CREATE TABLE legalities (
    card_id           UUID PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    standard          TEXT,
    pauper            TEXT,
    vintage           TEXT,
    pioneer           TEXT,
    modern            TEXT,
    legacy            TEXT,
    commander         TEXT,
	future            TEXT,
	historic          TEXT,
	timeless          TEXT,
	gladiator         TEXT,
	penny             TEXT,
	oathbreaker       TEXT,
	standardbrawl     TEXT,
	brawl             TEXT,
	alchemy           TEXT,
	paupercommander   TEXT,
	duel              TEXT,
	oldschool         TEXT,
	premodern         TEXT,
	predh             TEXT,
    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL
);

CREATE TABLE rulings (
    oracle_id         UUID NOT NULL,
    source            TEXT,
    published_at      DATE NOT NULL,
    comment           TEXT NOT NULL,
    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL,

    PRIMARY KEY (oracle_id, published_at, comment)
);

CREATE TABLE sync_state (
    key               TEXT PRIMARY KEY,
    last_sync         TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE card_faces;
DROP TABLE legalities;
DROP TABLE cards;
DROP TABLE rulings;
DROP TABLE sync_state;
