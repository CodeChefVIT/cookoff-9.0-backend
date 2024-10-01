-- name: CreateUser :one
INSERT INTO users (id, email, reg_no, password, role, round_qualified, score, name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE name = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: UpgradeUsersToRound :exec
UPDATE users
SET round_qualified = $2
WHERE id::TEXT = ANY($1::TEXT[]);

-- name: BanUser :exec
UPDATE users
SET is_banned = TRUE
WHERE id = $1;
-- name: UnbanUser :exec
UPDATE users
SET is_banned = FALSE
WHERE id = $1;

-- name: GetLeaderboard :many
select id, name, score from users
order by score;

-- name: UpdateProfile :exec
UPDATE users SET reg_no = $1, password = $2, name = $3
WHERE id = $4;

