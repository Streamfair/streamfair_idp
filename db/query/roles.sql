-- name: CreateRole :one
INSERT INTO "idp_svc"."Roles" (role_name, permissions) VALUES ($1, $2) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM "idp_svc"."Roles" WHERE id = $1 LIMIT 1;

-- name: GetRoleByName :one
SELECT * FROM "idp_svc"."Roles" WHERE role_name = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM "idp_svc"."Roles" ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateRole :one
UPDATE "idp_svc"."Roles" SET role_name = COALESCE($2, role_name), permissions = COALESCE($3, permissions) WHERE id = $1 RETURNING role_name, permissions;

-- name: DeleteRole :exec
DELETE FROM "idp_svc"."Roles" WHERE id = $1;
