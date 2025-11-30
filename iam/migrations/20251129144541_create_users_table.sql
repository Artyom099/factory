-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id          UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    login       TEXT                        NOT NULL,
    email       TEXT                        NOT NULL,
    hash        TEXT                        NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_users_login ON users(login);

CREATE TABLE notification_methods (
    id            UUID PRIMARY KEY  DEFAULT uuid_generate_v4(),
    user_id       UUID              NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_name TEXT              NOT NULL,
    target        TEXT              NOT NULL
);

CREATE INDEX idx_notification_methods_user_id ON notification_methods(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table notification_methods;
drop table users;
-- +goose StatementEnd
