-- +goose Up
-- +goose StatementBegin
-- Create payment_status table
CREATE TABLE payment_status (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Create payment_kind table
CREATE TABLE payment_kind (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Insert initial data
INSERT INTO payment_status (name) VALUES
    ('PAYMENT_CREATED'),
    ('PAYMENT_CONFIRMED'),
    ('PAYMENT_CANCELLED');

INSERT INTO payment_kind (name) VALUES
    ('debit'),
    ('credit');

-- Add columns as nullable
ALTER TABLE payments
    ADD COLUMN status_id INTEGER,
    ADD COLUMN kind_id INTEGER;

-- Set default values for existing rows (update sem where kk)
UPDATE payments 
    SET status_id = 1, -- PAYMENT_CREATED
    kind_id = 1;       -- credit

-- Add foreign key constraints
ALTER TABLE payments
    ADD CONSTRAINT fk_payment_status
        FOREIGN KEY (status_id)
        REFERENCES payment_status(id),
    ADD CONSTRAINT fk_payment_kind
        FOREIGN KEY (kind_id)
        REFERENCES payment_kind(id);

-- Make sure not null
ALTER TABLE payments
    ALTER COLUMN status_id SET NOT NULL,
    ALTER COLUMN kind_id SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove foreign key constraints first
ALTER TABLE payments
    DROP CONSTRAINT fk_payment_status,
    DROP CONSTRAINT fk_payment_kind;

-- Remove columns from payments table
ALTER TABLE payments
    DROP COLUMN status_id,
    DROP COLUMN kind_id;

-- Drop tables
DROP TABLE payment_kind;
DROP TABLE payment_status;
-- +goose StatementEnd