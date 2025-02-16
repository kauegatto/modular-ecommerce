-- +goose Up
-- +goose StatementBegin
CREATE TABLE Payment(
    id UUID PRIMARY KEY,
    orderID UUID NOT NULL,
    totalAmount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    integratorExternalID VARCHAR(42)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Payment
-- +goose StatementEnd
