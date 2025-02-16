-- +goose Up
-- +goose StatementBegin
ALTER TABLE payments
ALTER COLUMN orderID TYPE varchar (42);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE payments
ALTER COLUMN orderID TYPE UUID;
-- +goose StatementEnd
