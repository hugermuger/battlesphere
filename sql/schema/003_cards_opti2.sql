-- +goose Up
CREATE INDEX idx_cards_oracle_lang_release
ON cards (oracle_id, lang, release_date DESC)
WHERE ('paper' = ANY(games));

CREATE INDEX idx_card_faces_card_id ON card_faces (card_id);

CREATE INDEX idx_rulings_oracle_id ON rulings (oracle_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cards_oracle_lang_release;
DROP INDEX IF EXISTS idx_card_faces_card_id;
DROP INDEX IF EXISTS idx_rulings_oracle_id;
