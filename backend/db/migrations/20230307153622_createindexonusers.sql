-- +goose Up
-- +goose StatementBegin
-- create index on table users on row login
CREATE INDEX users_login_idx ON users(token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_login_idx;
-- +goose StatementEnd
