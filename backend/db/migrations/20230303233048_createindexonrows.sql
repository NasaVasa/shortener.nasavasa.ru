-- +goose Up
-- +goose StatementBegin
-- create indexes on table url_storage for columns long_url, short_url, secret_key
CREATE INDEX url_storage_long_url_idx ON url_storage (long_url);
CREATE INDEX url_storage_short_url_idx ON url_storage (short_url);
CREATE INDEX url_storage_secret_key_idx ON url_storage (secret_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop indexes on table url_storage for columns long_url, short_url, secret_key
DROP INDEX url_storage_long_url_idx;
DROP INDEX url_storage_short_url_idx;
DROP INDEX url_storage_secret_key_idx;
-- +goose StatementEnd
