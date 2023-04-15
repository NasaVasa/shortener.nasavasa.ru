-- +goose Up
-- +goose StatementBegin
-- in postgresql
-- create table url_storage with columns uuid id with default gen_random_uuid(),
-- string long_url,
-- string short_url,
-- uuid secret_key with default gen_random_uuid(),
-- int url_clicks with default 0,
-- timestamp url_created_at with default now() with time zone,
-- timestamp url_live_time with time zone,
-- boolean is_vip with default false,

CREATE TABLE url_storage (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    long_url text NOT NULL,
    short_url text NOT NULL,
    secret_key uuid NOT NULL DEFAULT gen_random_uuid(),
    url_clicks integer NOT NULL DEFAULT 0,
    url_created_at timestamp with time zone NOT NULL DEFAULT now(),
    url_will_delete timestamp with time zone,
    is_vip boolean NOT NULL DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop table url_storage;
DROP TABLE url_storage;
-- +goose StatementEnd
