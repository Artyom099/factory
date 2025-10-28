-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
    transaction_uuid VARCHAR DEFAULT '',
    payment_method VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE order_parts (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    part_id UUID NOT NULL,
    quantity INT DEFAULT 1,
    PRIMARY KEY (order_id, part_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_parts;
drop table orders;
-- +goose StatementEnd
