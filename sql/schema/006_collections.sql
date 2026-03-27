-- +goose Up
CREATE TABLE folders (
    id UUID NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    folder_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, folder_name)
);

CREATE TABLE collections (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    scryfall_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    oracle_id UUID NOT NULL,
    purchase_date TIME NOT NULL,
    purchase_price DOUBLE PRECISION NOT NULL DEFAULT '0.0',
    finish TEXT NOT NULL,
    condition TEXT NOT NULL,
    decks UUID[]
);

CREATE TABLE decks (
    id UUID NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    deck_name TEXT NOT NULL,
    format TEXT NOT NULL,
    commander UUID,
    commander_scryfall_id UUID,
    commander_collections_id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, deck_name)
);

CREATE TABLE deck_entries (
    id UUID PRIMARY KEY,
    deck_id UUID REFERENCES decks(id) ON DELETE CASCADE,
    oracle_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    scryfall_id UUID[],
    collections_id UUID[]
);

-- +goose Down
DROP TABLE deck_entries;
DROP TABLE decks;
DROP TABLE collections;
DROP TABLE folders;
