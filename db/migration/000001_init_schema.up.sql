CREATE SCHEMA "idp_svc";

CREATE TABLE "idp_svc"."Users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "full_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "password_salt" VARCHAR(255) NOT NULL,
  "country_code" VARCHAR(10) NOT NULL,
  "role_id" BIGINT,
  "status" VARCHAR(50),
  "last_login_at" TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00Z',
  "username_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "email_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "idp_svc"."Tokens" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "token" VARCHAR(255) NOT NULL,
  "expires_at" TIMESTAMPTZ NOT NULL,
  "revoked" BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE "idp_svc"."RefreshTokens" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "token" VARCHAR(255) NOT NULL,
  "expires_at" TIMESTAMPTZ NOT NULL,
  "revoked" BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE "idp_svc"."Roles" (
  "id" BIGSERIAL PRIMARY KEY,
  "role_name" VARCHAR(50) UNIQUE NOT NULL,
  "permissions" JSONB NOT NULL
);

CREATE TABLE "idp_svc"."UserRoles" (
  "user_id" BIGINT NOT NULL,
  "role_id" BIGINT NOT NULL,
  PRIMARY KEY ("user_id", "role_id")
);

CREATE INDEX "idx_users_id" ON "idp_svc"."Users" ("id");

CREATE INDEX "idx_users_username" ON "idp_svc"."Users" ("username");

CREATE INDEX "idx_users_email" ON "idp_svc"."Users" ("email");

CREATE INDEX "idx_tokens_id" ON "idp_svc"."Tokens" ("id");

CREATE INDEX "idx_tokens_user_id" ON "idp_svc"."Tokens" ("user_id");

CREATE INDEX "idx_refresh_tokens_id" ON "idp_svc"."RefreshTokens" ("id");

CREATE INDEX "idx_refresh_tokens_user_id" ON "idp_svc"."RefreshTokens" ("user_id");

CREATE INDEX "idx_roles_id" ON "idp_svc"."Roles" ("id");

CREATE INDEX "idx_user_roles_user_id" ON "idp_svc"."UserRoles" ("user_id");

CREATE INDEX "idx_user_roles_role_id" ON "idp_svc"."UserRoles" ("role_id");

ALTER TABLE "idp_svc"."Tokens" ADD FOREIGN KEY ("user_id") REFERENCES "idp_svc"."Users" ("id");

ALTER TABLE "idp_svc"."RefreshTokens" ADD FOREIGN KEY ("user_id") REFERENCES "idp_svc"."Users" ("id");

ALTER TABLE "idp_svc"."UserRoles" ADD FOREIGN KEY ("user_id") REFERENCES "idp_svc"."Users" ("id");

ALTER TABLE "idp_svc"."UserRoles" ADD FOREIGN KEY ("role_id") REFERENCES "idp_svc"."Roles" ("id");
