-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS doc_chunks (
  id SERIAL PRIMARY KEY,
  text TEXT NOT NULL,
  text_tsvector TSVECTOR
);

CREATE INDEX IF NOT EXISTS text_tsvector_idx
    ON doc_chunks USING GIN (text_tsvector);

CREATE OR REPLACE FUNCTION update_text_tsvector() RETURNS TRIGGER AS $$
BEGIN
            NEW.text_tsvector = to_tsvector('russian', NEW.text);
RETURN NEW;
END;
        $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS text_tsvector_trigger ON doc_chunks;

CREATE TRIGGER text_tsvector_trigger
    BEFORE INSERT OR UPDATE ON doc_chunks
                         FOR EACH ROW EXECUTE FUNCTION update_text_tsvector();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER text_tsvector_trigger;
DROP FUNCTION update_text_tsvector();
DROP INDEX text_tsvector_idx;
DROP TABLE doc_chunks;
-- +goose StatementEnd

