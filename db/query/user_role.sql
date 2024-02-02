-- name: AssignRoleToUser :one
INSERT INTO "idp_svc"."UserRoles" (user_id, role_id) VALUES ($1, $2) RETURNING *;

-- name: GetUserRoles :many
SELECT * FROM "idp_svc"."UserRoles" WHERE user_id = $1;

-- name: CheckIfUserHasRole :one
SELECT * FROM "idp_svc"."UserRoles" WHERE user_id = $1 AND role_id = $2;

-- name: RemoveRoleFromUser :exec
DELETE FROM "idp_svc"."UserRoles" WHERE user_id = $1 AND role_id = $2;
