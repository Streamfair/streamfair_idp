-- Drop unique constraint for uuid
ALTER TABLE "idp_svc"."UserAccounts" DROP CONSTRAINT IF EXISTS "fk_session_uuid";

-- Drop foreign key constraint for owner
ALTER TABLE "idp_svc"."UserAccounts" DROP CONSTRAINT IF EXISTS "unique_account";

-- Drop foreign key constraint for account_type
ALTER TABLE "idp_svc"."UserAccounts" DROP CONSTRAINT IF EXISTS "fk_account_type";

-- Drop indices
DROP INDEX IF EXISTS "idp_svc"."unique_account";
DROP INDEX IF EXISTS "idp_svc"."fk_session_uuid";

-- Drop table
DROP TABLE IF EXISTS "idp_svc"."UserAccounts";

-- Drop schema if it exists and is empty
DROP SCHEMA IF EXISTS "idp_svc" CASCADE;
