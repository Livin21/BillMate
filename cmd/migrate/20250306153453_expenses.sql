-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE TABLE
    expenses (
        id UUID PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        user_id UUID NOT NULL,
        amount DECIMAL(10, 2) NOT NULL,
        description TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';
DROP TABLE expenses;

-- +goose StatementEnd