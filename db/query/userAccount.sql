-- name: GetAllUsers :many
SELECT *
FROM "idp_svc"."UserAccount_View";

-- name: GetUserById :one
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE email = $1;

-- name: GetUsersByAccountType :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE account_type = $1;

-- name: GetUsersByAccountStatus :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE account_status = $1;

-- name: GetUsersByRoleId :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE role_id = $1;

-- name: GetUsersCreatedAfter :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE user_created_at > $1;

-- name: GetUsersUpdatedAfter :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE user_updated_at > $1;

-- name: GetUsersByCountryCode :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE country_code = $1;

-- name: GetUsersByOwner :many
SELECT *
FROM "idp_svc"."UserAccount_View"
WHERE owner = $1;
