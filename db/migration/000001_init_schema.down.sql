-- Drop the view if it exists
DROP VIEW IF EXISTS "idp_svc"."UserAccount_View";

-- Drop the unique constraint
ALTER TABLE "idp_svc"."Accounts" DROP CONSTRAINT IF EXISTS "unique_account";

-- Drop the indices
DROP INDEX IF EXISTS "idx_user_id";
DROP INDEX IF EXISTS "idx_user_username";
DROP INDEX IF EXISTS "idx_users_email";
DROP INDEX IF EXISTS "idx_acc_owner";
DROP INDEX IF EXISTS "idx_accType_id";

-- Drop the tables if they exist
DROP TABLE IF EXISTS "idp_svc"."Accounts";
DROP TABLE IF EXISTS "idp_svc"."AccountTypes";
DROP TABLE IF EXISTS "idp_svc"."Users";

-- Drop the schema if it exists
DROP SCHEMA IF EXISTS "idp_svc";
