-- :name CreateUser :one
INSERT INTO "user" (id, submissions, email, "regNo", password, role, "roundQualified", "score", name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
-- name: GetUserByEmail :one
SELECT id, email, "regNo", password, role, "roundQualified", "score", name
FROM "user"
WHERE email = $1;
-- name: GetUserByUsername :one
SELECT id, email, "regNo", password, role, "roundQualified", "score", name
FROM "user"
WHERE name = $1;
-- name: GetUserById :one
SELECT id, email, "regNo", password, role, "roundQualified", "score", name
FROM "user"
WHERE id = $1;