-- +goose Up
-- +goose StatementBegin
INSERT INTO payment_status (name) VALUES
    ('PAYMENT_PENDING_REFUND'),
    ('PAYMENT_REFUNDED'),
    ('PAYMENT_PENDING');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
