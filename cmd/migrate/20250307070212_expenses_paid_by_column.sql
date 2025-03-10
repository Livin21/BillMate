-- +goose Up
-- +goose StatementBegin
ALTER TABLE expenses
ADD COLUMN paid_by UUID REFERENCES users (id);

ALTER TABLE expenses
ADD COLUMN paid_at TIMESTAMP NOT NULL DEFAULT NOW ();

ALTER TABLE expenses
ADD COLUMN paid_by_name VARCHAR(255);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE expenses
DROP COLUMN paid_by;

ALTER TABLE expenses
DROP COLUMN paid_at;

ALTER TABLE expenses
DROP COLUMN paid_by_name;

-- +goose StatementEnd