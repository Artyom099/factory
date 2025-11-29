-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id                      UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    login                   TEXT                        NOT NULL,
    email                   TEXT                        NOT NULL,
    hash                    TEXT                        NOT NULL,
    notification_methods    TEXT
    created_at              TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT now(),
    updated_at              TIMESTAMP WITH TIME ZONE
);

CREATE TABLE notification_methods (
    id                  UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    provider_name       TEXT                        NOT NULL,
    target              TEXT                        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table notification_methods;
drop table users;
-- +goose StatementEnd
