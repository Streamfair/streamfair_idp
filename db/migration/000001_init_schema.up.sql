CREATE SCHEMA "idp_svc";

CREATE TABLE "idp_svc"."UserAccounts" (
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
    "user_created_at" timestamptz NOT NULL DEFAULT now(),
    "user_updated_at" timestamptz NOT NULL DEFAULT now(),
    "account_name" varchar NOT NULL,
    "account_type" int NOT NULL,
    "owner" varchar NOT NULL,
    "bio" varchar NOT NULL,
    "account_status" varchar NOT NULL,
    "plan" int NOT NULL,
    "avatar_uri" varchar,
    "plays" bigint NOT NULL DEFAULT 0,
    "likes" bigint NOT NULL DEFAULT 0,
    "follows" bigint NOT NULL DEFAULT 0,
    "shares" bigint NOT NULL DEFAULT 0,
    "account_created_at" timestamptz NOT NULL DEFAULT now(),
    "account_updated_at" timestamptz NOT NULL DEFAULT now(),
    "type" int NOT NULL,
    "permissions" varchar NOT NULL,
    "is_artist" boolean NOT NULL DEFAULT false,
    "is_producer" boolean NOT NULL DEFAULT false,
    "is_writer" boolean NOT NULL DEFAULT false,
    "is_label" boolean NOT NULL DEFAULT false,
    "is_user" boolean NOT NULL DEFAULT false,
    "account_type_created_at" timestamptz NOT NULL DEFAULT now(),
    "account_type_updated_at" timestamptz NOT NULL DEFAULT now(),
    "uuid" uuid NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "session_expires_at" timestamptz NOT NULL,
    "session_created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create indices
CREATE INDEX "idx_user_id" ON "idp_svc"."UserAccounts" ("id");
CREATE INDEX "idx_user_username" ON "idp_svc"."UserAccounts" ("username");
CREATE INDEX "idx_users_email" ON "idp_svc"."UserAccounts" ("email");
CREATE INDEX "idx_acc_owner" ON "idp_svc"."UserAccounts" ("owner");

-- Add unique constraint
ALTER TABLE "idp_svc"."UserAccounts" ADD CONSTRAINT "unique_account" UNIQUE ("owner", "account_type");
-- ALTER TABLE "idp_svc"."UserAccounts" ADD CONSTRAINT "fk_session_uuid" UNIQUE ("uuid");
