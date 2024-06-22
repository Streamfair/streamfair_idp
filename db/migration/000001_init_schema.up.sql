-- Create the schema
CREATE SCHEMA IF NOT EXISTS "idp_svc";

-- Create Users table
CREATE TABLE IF NOT EXISTS "idp_svc"."Users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar UNIQUE NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_hash" varchar NOT NULL,
    "password_salt" varchar NOT NULL,
    "country_code" varchar NOT NULL,
    "role_id" bigint,
    "status" varchar,
    "last_login_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "username_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "email_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Create AccountTypes table
CREATE TABLE IF NOT EXISTS "idp_svc"."AccountTypes" (
    "id" serial PRIMARY KEY,
    "type" int NOT NULL,
    "permissions" varchar NOT NULL,
    "is_artist" boolean NOT NULL DEFAULT false,
    "is_producer" boolean NOT NULL DEFAULT false,
    "is_writer" boolean NOT NULL DEFAULT false,
    "is_label" boolean NOT NULL DEFAULT false,
    "is_user" boolean NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Create Accounts table
CREATE TABLE IF NOT EXISTS "idp_svc"."Accounts" (
    "id" bigserial PRIMARY KEY,
    "account_name" varchar NOT NULL,
    "account_type" int NOT NULL,
    "owner" varchar NOT NULL,
    "bio" varchar NOT NULL,
    "status" varchar NOT NULL,
    "plan" int NOT NULL,
    "avatar_uri" varchar,
    "plays" bigint NOT NULL DEFAULT 0,
    "likes" bigint NOT NULL DEFAULT 0,
    "follows" bigint NOT NULL DEFAULT 0,
    "shares" bigint NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_account_type FOREIGN KEY ("account_type") REFERENCES "idp_svc"."AccountTypes" ("id")
);

-- Create indices
CREATE INDEX IF NOT EXISTS "idx_user_id" ON "idp_svc"."Users" ("id");
CREATE INDEX IF NOT EXISTS "idx_user_username" ON "idp_svc"."Users" ("username");
CREATE INDEX IF NOT EXISTS "idx_users_email" ON "idp_svc"."Users" ("email");
CREATE INDEX IF NOT EXISTS "idx_acc_owner" ON "idp_svc"."Accounts" ("owner");
CREATE INDEX IF NOT EXISTS "idx_accType_id" ON "idp_svc"."AccountTypes" ("id");

-- Add unique constraint
ALTER TABLE "idp_svc"."Accounts" ADD CONSTRAINT "unique_account" UNIQUE ("owner", "account_type");

-- Create view for combined user, account, and account type details
CREATE OR REPLACE VIEW "idp_svc"."UserAccount_View" AS
SELECT
    u.id AS user_id,
    u.username,
    u.full_name,
    u.email,
    u.password_hash,
    u.password_salt,
    u.country_code,
    u.role_id,
    u.status AS user_status,
    u.last_login_at,
    u.username_changed_at,
    u.email_changed_at,
    u.password_changed_at,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    a.id AS account_id,
    a.account_name,
    a.account_type,
    a.owner,
    a.bio,
    a.status AS account_status,
    a.plan,
    a.avatar_uri,
    a.plays,
    a.likes,
    a.follows,
    a.shares,
    at.id AS account_type_id,
    at.type AS account_type_name,
    at.permissions,
    at.is_artist,
    at.is_producer,
    at.is_writer,
    at.is_label,
    at.is_user AS is_regular_user
FROM "idp_svc"."Users" u
JOIN "idp_svc"."Accounts" a ON u.username = a.account_name
JOIN "idp_svc"."AccountTypes" at ON a.account_type = at.id;
