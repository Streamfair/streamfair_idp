CREATE SCHEMA "idp_svc";

CREATE TABLE "idp_svc"."Roles" (
  "id" BIGSERIAL PRIMARY KEY,
  "role_name" VARCHAR(50) UNIQUE NOT NULL,
  "permissions" JSONB NOT NULL
);

CREATE TABLE "idp_svc"."UserRoles" (
  "user_id" BIGINT NOT NULL,
  "role_id" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "role_id")
);

CREATE INDEX "idx_roles_id" ON "idp_svc"."Roles" ("id");

CREATE INDEX "idx_user_roles_user_id" ON "idp_svc"."UserRoles" ("user_id");

CREATE INDEX "idx_user_roles_role_id" ON "idp_svc"."UserRoles" ("role_id");

CREATE INDEX "idx_user_roles_created_at" ON "idp_svc"."UserRoles" ("created_at");

ALTER TABLE "idp_svc"."UserRoles" ADD FOREIGN KEY ("role_id") REFERENCES "idp_svc"."Roles" ("id");
