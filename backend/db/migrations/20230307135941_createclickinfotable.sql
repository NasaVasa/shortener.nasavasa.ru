-- +goose Up
-- +goose StatementBegin
-- create table with id, foreign key to url_storage id, and a timestamp
CREATE TABLE IF NOT EXISTS click_info (
  id SERIAL PRIMARY KEY,
  url_id UUID NOT NULL REFERENCES url_storage(id) ON DELETE CASCADE,
  click_time TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS click_info;
-- +goose StatementEnd
