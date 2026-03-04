-- name: InsertCard :exec
INSERT INTO cards (
    id, arena_id, mtgo_id, cardmarket_id, oracle_id, release_date,
    lang, layout, edhrec_rank, game_changer, multifaced,
    cmc, color_identity, colors, defense, keywords, loyalty,
    mana_cost, name, oracle_text, power, toughness, type_line,
    artist, collector_number, finishes, flavor_name, flavor_text,
    games, image, image_png, image_large, image_small, image_crop,
    price_usd, price_eur, price_foil_usd, price_foil_eur, price_etched_usd,
    printed_name, printed_text, printed_type_line,
    rarity, set_name, set_code,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17,
    $18, $19, $20, $21, $22, $23,
    $24, $25, $26, $27, $28,
    $29, $30, $31, $32,
    $33, $34, $35, $36,
    $37, $38, $39,
    $40, $41, $42, $43, $44, $45,
    NOW(), NOW()
)
ON CONFLICT (id) DO UPDATE SET
    arena_id = EXCLUDED.arena_id,
    mtgo_id = EXCLUDED.mtgo_id,
    cardmarket_id = EXCLUDED.cardmarket_id,
    oracle_id = EXCLUDED.oracle_id,
    release_date = EXCLUDED.release_date,
    lang = EXCLUDED.lang,
    layout = EXCLUDED.layout,
    edhrec_rank = EXCLUDED.edhrec_rank,
    game_changer = EXCLUDED.game_changer,
    multifaced = EXCLUDED.multifaced,
    cmc = EXCLUDED.cmc,
    color_identity = EXCLUDED.color_identity,
    colors = EXCLUDED.colors,
    defense = EXCLUDED.defense,
    keywords = EXCLUDED.keywords,
    loyalty = EXCLUDED.loyalty,
    mana_cost = EXCLUDED.mana_cost,
    name = EXCLUDED.name,
    oracle_text = EXCLUDED.oracle_text,
    power = EXCLUDED.power,
    toughness = EXCLUDED.toughness,
    type_line = EXCLUDED.type_line,
    artist = EXCLUDED.artist,
    collector_number = EXCLUDED.collector_number,
    finishes = EXCLUDED.finishes,
    flavor_name = EXCLUDED.flavor_name,
    flavor_text = EXCLUDED.flavor_text,
    games = EXCLUDED.games,
    image = EXCLUDED.image,
    image_png = EXCLUDED.image_png,
    image_large = EXCLUDED.image_large,
    image_small = EXCLUDED.image_small,
    image_crop = EXCLUDED.image_crop,
    price_usd = EXCLUDED.price_usd,
    price_eur = EXCLUDED.price_eur,
    price_foil_usd = EXCLUDED.price_foil_usd,
    price_foil_eur = EXCLUDED.price_foil_eur,
    price_etched_usd = EXCLUDED.price_etched_usd,
    printed_name = EXCLUDED.printed_name,
    printed_text = EXCLUDED.printed_text,
    printed_type_line = EXCLUDED.printed_type_line,
    rarity = EXCLUDED.rarity,
    set_name = EXCLUDED.set_name,
    set_code = EXCLUDED.set_code,
    updated_at = NOW();

-- name: InsertCardFace :exec
INSERT INTO card_faces (
    card_id, face_index,
    artist, cmc, colors, defense,
    flavor_text, image, image_png, image_large, image_small, image_crop,
    layout, loyalty, mana_cost, name, oracle_text,
    power, printed_name, printed_text, printed_type_line,
    toughness, type_line, created_at, updated_at
) VALUES (
    $1, $2, $3,
    $4, $5, $6, $7, $8,
    $9, $10, $11, $12,
    $13, $14, $15, $16, $17,
    $18, $19, $20, $21,
    $22, $23, Now(), Now()
)
ON CONFLICT (card_id, face_index) DO UPDATE SET
    artist = EXCLUDED.artist,
    cmc = EXCLUDED.cmc,
    colors = EXCLUDED.colors,
    defense = EXCLUDED.defense,
    flavor_text = EXCLUDED.flavor_text,
    image = EXCLUDED.image,
    image_png = EXCLUDED.image_png,
    image_large = EXCLUDED.image_large,
    image_small = EXCLUDED.image_small,
    image_crop = EXCLUDED.image_crop,
    layout = EXCLUDED.layout,
    loyalty = EXCLUDED.loyalty,
    mana_cost = EXCLUDED.mana_cost,
    name = EXCLUDED.name,
    oracle_text = EXCLUDED.oracle_text,
    power = EXCLUDED.power,
    printed_name = EXCLUDED.printed_name,
    printed_text = EXCLUDED.printed_text,
    printed_type_line = EXCLUDED.printed_type_line,
    toughness = EXCLUDED.toughness,
    type_line = EXCLUDED.type_line,
    updated_at = Now();

-- name: InsertLegality :exec
INSERT INTO legalities (
    card_id, standard, pauper, vintage, pioneer, modern, legacy, commander, future,
	historic, timeless, gladiator, penny, oathbreaker,
	standardbrawl, brawl, alchemy, paupercommander, duel,
	oldschool, premodern, predh, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, Now(), Now()
)
ON CONFLICT (card_id) DO UPDATE SET
    standard = EXCLUDED.standard,
    pauper = EXCLUDED.pauper,
    vintage = EXCLUDED.vintage,
    pioneer = EXCLUDED.pioneer,
    modern = EXCLUDED.modern,
    legacy = EXCLUDED.legacy,
    commander = EXCLUDED.commander,
	future = EXCLUDED.future,
	historic = EXCLUDED.historic,
	timeless = EXCLUDED.timeless,
	gladiator = EXCLUDED.gladiator,
	penny = EXCLUDED.penny,
	oathbreaker = EXCLUDED.oathbreaker,
	standardbrawl = EXCLUDED.standardbrawl,
	brawl = EXCLUDED.brawl,
	alchemy = EXCLUDED.alchemy,
	paupercommander = EXCLUDED.paupercommander,
	duel = EXCLUDED.duel,
	oldschool = EXCLUDED.oldschool,
	premodern = EXCLUDED.premodern,
	predh = EXCLUDED.predh,
    updated_at = Now();

-- name: InsertRulings :exec
INSERT INTO rulings (
    oracle_id, source, published_at, comment, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, Now(), Now()
)
ON CONFLICT (oracle_id, published_at, comment) DO NOTHING;

-- name: InsertSyncState :exec
INSERT INTO sync_state (
key, last_sync
) VALUES (
    $1, Now()
)
ON CONFLICT (key) DO UPDATE SET
    last_sync = Now();

-- name: GetSyncState :one
SELECT last_sync FROM sync_state WHERE key = $1;
