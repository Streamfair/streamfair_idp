-- name: CreateToken :one
INSERT INTO "idp_svc"."Tokens" (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTokenByID :one
SELECT * FROM "idp_svc"."Tokens" WHERE id = $1 LIMIT 1;

-- name: GetTokenByValue :one
SELECT * FROM "idp_svc"."Tokens" WHERE token = $1 LIMIT 1;

-- name: ListTokens :many
SELECT * FROM "idp_svc"."Tokens" ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateToken :one
UPDATE "idp_svc"."Tokens" SET token = COALESCE($2, token), expires_at = COALESCE($3, expires_at) WHERE id = $1 RETURNING token, expires_at;

-- name: DeleteToken :exec
DELETE FROM "idp_svc"."Tokens" WHERE id = $1;
