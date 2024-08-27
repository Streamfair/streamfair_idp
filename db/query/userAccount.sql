-- name: CreateUserAccount :one
INSERT INTO "idp_svc"."UserAccounts" (
    username,
    full_name,
    email,
    password_hash,
    password_salt,
    country_code,
    role_id,
    status,
    last_login_at,
    username_changed_at,
    email_changed_at,
    password_changed_at,
    user_created_at,
    user_updated_at,
    account_name,
    account_type,
    owner,
    bio,
    account_status,
    plan,
    avatar_uri,
    plays,
    likes,
    follows,
    shares,
    account_created_at,
    account_updated_at,
    type,
    permissions,
    is_artist,
    is_producer,
    is_writer,
    is_label,
    is_user,
    account_type_created_at,
    account_type_updated_at,
    uuid,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    session_expires_at,
    session_created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43
)
RETURNING *;

-- name: ListUserAccounts :many
SELECT
    id,
    username,
    full_name,
    email,
    country_code,
    role_id,
    status,
    last_login_at,
    user_created_at,
    user_updated_at,
    account_name,
    account_type,
    owner,
    bio,
    account_status,
    plan,
    avatar_uri,
    plays,
    likes,
    follows,
    shares,
    is_artist,
    is_producer,
    is_writer,
    is_label,
    is_user
FROM "idp_svc"."UserAccounts"
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUserAccount :one
UPDATE "idp_svc"."UserAccounts"
SET 
    username = COALESCE($1, username),
    full_name = COALESCE($2, full_name),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash),
    password_salt = COALESCE($5, password_salt),
    country_code = COALESCE($6, country_code),
    role_id = COALESCE($7, role_id),
    status = COALESCE($8, status),
    last_login_at = COALESCE($9, last_login_at),
    username_changed_at = COALESCE($10, username_changed_at),
    email_changed_at = COALESCE($11, email_changed_at),
    password_changed_at = COALESCE($12, password_changed_at),
    user_created_at = COALESCE($13, user_created_at),
    user_updated_at = NOW(),
    account_name = COALESCE($14, account_name),
    account_type = COALESCE($15, account_type),
    owner = COALESCE($16, owner),
    bio = COALESCE($17, bio),
    account_status = COALESCE($18, account_status),
    plan = COALESCE($19, plan),
    avatar_uri = COALESCE($20, avatar_uri),
    plays = COALESCE($21, plays),
    likes = COALESCE($22, likes),
    follows = COALESCE($23, follows),
    shares = COALESCE($24, shares),
    is_artist = COALESCE($25, is_artist),
    is_producer = COALESCE($26, is_producer),
    is_writer = COALESCE($27, is_writer),
    is_label = COALESCE($28, is_label),
    is_user = COALESCE($29, is_user)
WHERE id = $30
RETURNING *;

-- name: DeleteUserAccountById :exec
DELETE FROM "idp_svc"."UserAccounts"
WHERE id = $1;

-- name: DeleteUserAccountByValue :exec
DELETE FROM "idp_svc"."UserAccounts"
WHERE username = $1 OR email = $1;

-- name: GetAllUserAccounts :many
SELECT * FROM "idp_svc"."UserAccounts";

-- name: GetUserAccountById :one
SELECT * FROM "idp_svc"."UserAccounts"
WHERE id = $1;

-- name: GetUserAccountByUserAccountname :one
SELECT * FROM "idp_svc"."UserAccounts"
WHERE username = $1;

-- name: GetUserAccountByEmail :one
SELECT * FROM "idp_svc"."UserAccounts"
WHERE email = $1;

-- name: GetUserAccountsByAccountType :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE account_type = $1;

-- name: GetUserAccountsByAccountStatus :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE account_status = $1;

-- name: GetUserAccountsByRoleId :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE role_id = $1;

-- name: GetUserAccountsCreatedAfter :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE user_created_at > $1;

-- name: GetUserAccountsUpdatedAfter :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE user_updated_at > $1;

-- name: GetUserAccountsByCountryCode :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE country_code = $1;

-- name: GetUserAccountsByOwner :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE owner = $1;

-- name: GetUserAccountWithActiveSessions :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE session_is_blocked = false;

-- name: GetUserAccountWithBlockedSessions :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE session_is_blocked = true;

-- name: GetUserAccountWithPermissions :many
SELECT * FROM "idp_svc"."UserAccounts"
WHERE permissions = $1;

-- name: OrderResultsByLastLoginTime :many
SELECT * FROM "idp_svc"."UserAccounts"
ORDER BY last_login_at DESC;

-- name: CountTotalNumberOfSessionsPerUser :many
SELECT
    username,
    COUNT(session_uuid) AS total_sessions
FROM "idp_svc"."UserAccounts"
GROUP BY username;

-- name: FetchDataForSpecificDateRange :many
SELECT *
FROM "idp_svc"."UserAccounts"
WHERE user_created_at BETWEEN $1 AND $2;
