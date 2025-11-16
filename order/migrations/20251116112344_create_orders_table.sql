-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id                  UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    user_uuid           UUID                        NOT NULL,
    total_price         NUMERIC(10, 2)              NOT NULL,
    transaction_uuid    UUID,
    payment_method      TEXT                        NOT NULL,
    status              TEXT                        NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP WITH TIME ZONE
);

CREATE TABLE orders_parts (
    order_id    UUID    NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    part_id     UUID    NOT NULL,
    quantity    INT     DEFAULT 1,
    PRIMARY KEY (order_id, part_id)
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders (user_uuid);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders_parts;
drop table orders;
drop index if exists idx_orders_user_uuid;
drop index if exists idx_orders_status;
-- +goose StatementEnd
