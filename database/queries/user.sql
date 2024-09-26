-- name: CreateUser :one
INSERT INTO users (id, email, reg_no, password, role, round_qualified, score, name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, email, reg_no, password, role, round_qualified, score, name
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT id, email, reg_no, password, role, round_qualified, score, name
FROM users
WHERE name = $1;

-- name: GetUserById :one
SELECT id, email, reg_no, password, role, round_qualified, score, name
FROM users
WHERE id = $1;
-- name: GetAllUsers :many
SELECT id, email, reg_no, password, role, round_qualified, score, name
FROM users;
-- name: UpgradeUsersToRound :batchexec
UPDATE users
SET round_qualified = GREATEST(round_qualified, $2)
WHERE id = $1;
-- name: BanUser :exec
UPDATE users
SET is_banned = TRUE
WHERE id = $1;
-- name: UnbanUser :exec
UPDATE users
SET is_banned = FALSE
WHERE id = $1;