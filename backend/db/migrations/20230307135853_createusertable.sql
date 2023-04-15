-- +goose Up
-- +goose StatementBegin
-- create user table with UUID id default get_random_uuid,
-- name,
-- email,
-- hashed password,
-- UUID token,
-- timestamp token_will_delete
CREATE TABLE IF NOT EXISTS users(
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    login VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token UUID,
    token_created_at TIMESTAMP NOT NULL DEFAULT now(),
    user_created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
