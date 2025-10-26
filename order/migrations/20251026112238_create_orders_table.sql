-- +goose Up
-- +goose StatementBegin
create table orders (
    order_uuid uuid primary key default gen_random_uuid(),
    user_uuid uuid not null,
    part_uuids uuid[] not null,
    total_price numeric(10, 2) not null,
    transaction_uuid varchar default '',
    payment_method varchar not null,
    status varchar not null,
    created_at timestamp not null default now(),
    updated_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order;
-- +goose StatementEnd
