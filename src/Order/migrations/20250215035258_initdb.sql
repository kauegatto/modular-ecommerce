-- +goose Up

-- +goose StatementBegin
CREATE TABLE status (
    id INT PRIMARY KEY,
    status_name VARCHAR(32) NOT NULL
);
INSERT INTO status (id, status_name) VALUES
(1, 'ORDER_PLACED'),
(2, 'ORDER_CONFIRMED'),
(3, 'ORDER_CANCELLED'),
(4, 'ORDER_PENDING');
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE items (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    price BIGINT NOT NULL
);
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    customer_id VARCHAR(32) NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    total_price BIGINT NOT NULL,
    discount DECIMAL(5, 2),
    CONSTRAINT fk_status FOREIGN KEY (status_id) REFERENCES status(id)
);

CREATE TABLE order_items (
    order_id UUID NOT NULL,
    item_id UUID NOT NULL,
    PRIMARY KEY (order_id, item_id),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE status
DROP TABLE items
DROP TABLE orders
DROP TABLE order_items
-- +goose StatementEnd






