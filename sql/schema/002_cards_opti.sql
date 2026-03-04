-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_cards_name_trgm
ON cards USING gin (name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_cards_printed_name_trgm
ON cards USING gin (printed_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_cards_games_gin
ON cards USING gin (games);

CREATE INDEX IF NOT EXISTS idx_cards_lang_rarity
ON cards (lang, rarity);

-- +goose Down
DROP INDEX IF EXISTS idx_cards_lang_rarity;
DROP INDEX IF EXISTS idx_cards_games_gin;
DROP INDEX IF EXISTS idx_cards_printed_name_trgm;
DROP INDEX IF EXISTS idx_cards_name_trgm;
