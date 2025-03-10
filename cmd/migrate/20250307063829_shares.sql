-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    shares (
        id UUID PRIMARY KEY,
        user_id UUID,
        user_name VARCHAR(255),
        expense_id UUID NOT NULL,
        amount_owed DECIMAL(10, 2) NOT NULL,
        amount_paid DECIMAL(10, 2) NOT NULL,
        created_by UUID NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (expense_id) REFERENCES expenses (id),
        FOREIGN KEY (created_by) REFERENCES users (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE shares;

-- +goose StatementEnd