-- +goose Up
-- +goose StatementBegin
-- add a column to url_storage table user_id uuid that references users table on cascade delete
ALTER TABLE url_storage ADD COLUMN user_id uuid REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop column user_id from url_storage table
ALTER TABLE url_storage DROP COLUMN user_id;
-- +goose StatementEnd
