-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
   access_token VARCHAR(512) NOT NULL,
   refresh_token VARCHAR(512) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd