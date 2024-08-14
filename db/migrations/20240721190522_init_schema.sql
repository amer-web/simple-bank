-- +goose Up
CREATE TABLE "accounts"
(
    "id"         bigserial PRIMARY KEY,
    "balance"    bigint  NOT NULL,
    "owner"      varchar NOT NULL,
    "currency"   varchar NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "entries"
(
    "id"         bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL ,
    "amount"     bigint NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transfers"
(
    "id"              bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id"   bigint NOT NULL,
    "amount"          bigint NOT NULL,
    "created_at"      timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

ALTER TABLE "entries"
    ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");


-- +goose Down
-- Drop foreign key constraints on 'transfers' table
ALTER TABLE entries DROP CONSTRAINT IF EXISTS entries_account_id_fkey;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS transfers_from_account_id_fkey;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS transfers_to_account_id_fkey;

DROP table if exists accounts;
DROP table if exists entries;
DROP table if exists transfers;
