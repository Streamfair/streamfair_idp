-- name: CreateRefreshToken :one
INSERT INTO "idp_svc"."RefreshTokens" (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: GetRefreshTokenByID :one
SELECT * FROM "idp_svc"."RefreshTokens" WHERE id = $1 LIMIT 1;

-- name: GetRefreshTokenByValue :one
SELECT * FROM "idp_svc"."RefreshTokens" WHERE token = $1 LIMIT 1;

-- name: ListRefreshTokens :many
SELECT * FROM "idp_svc"."RefreshTokens" ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateRefreshToken :one
UPDATE "idp_svc"."RefreshTokens" SET token = COALESCE($2, token), expires_at = COALESCE($3, expires_at) WHERE id = $1 RETURNING token, expires_at;

-- name: DeleteRefreshToken :exec
DELETE FROM "idp_svc"."RefreshTokens" WHERE id = $1;
