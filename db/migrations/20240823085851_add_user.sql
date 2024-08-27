-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users"
(
    "username"            varchar PRIMARY KEY,
    "password"            varchar        NOT NULL,
    "full_name"           varchar        NOT NULL,
    "email"               varchar UNIQUE NOT NULL,
    "password_changed_at" timestamptz              DEFAULT '0001-01-01 00:00:00z',
    "created_at"          timestamptz    NOT NULL DEFAULT (now())
);
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP table IF EXISTS users;
-- +goose StatementEnd
