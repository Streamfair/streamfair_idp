CREATE SCHEMA "idp_svc";

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

CREATE INDEX "idx_tokens_id" ON "idp_svc"."Tokens" ("id");

CREATE INDEX "idx_tokens_user_id" ON "idp_svc"."Tokens" ("user_id");

CREATE INDEX "idx_refresh_tokens_id" ON "idp_svc"."RefreshTokens" ("id");

CREATE INDEX "idx_refresh_tokens_user_id" ON "idp_svc"."RefreshTokens" ("user_id");

CREATE INDEX "idx_roles_id" ON "idp_svc"."Roles" ("id");

CREATE INDEX "idx_user_roles_user_id" ON "idp_svc"."UserRoles" ("user_id");

CREATE INDEX "idx_user_roles_role_id" ON "idp_svc"."UserRoles" ("role_id");

ALTER TABLE "idp_svc"."UserRoles" ADD FOREIGN KEY ("role_id") REFERENCES "idp_svc"."Roles" ("id");
