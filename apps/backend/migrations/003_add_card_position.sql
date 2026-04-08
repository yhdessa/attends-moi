ALTER TABLE cards ADD COLUMN position DOUBLE PRECISION;

UPDATE cards
SET position = EXTRACT(EPOCH FROM created_at) * 1000
WHERE position IS NULL;

ALTER TABLE cards ALTER COLUMN position SET NOT NULL;
ALTER TABLE cards ALTER COLUMN position SET DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000);

CREATE INDEX idx_cards_board_status_position ON cards(board_id, status, position);
