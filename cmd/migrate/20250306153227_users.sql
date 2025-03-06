-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        password bytea NOT NULL,
        phone VARCHAR(15),
        email CITEXT UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

DROP EXTENSION IF EXISTS citext;

-- +goose StatementEnd