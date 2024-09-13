-- name: CreateUser :one
INSERT INTO users (username,
                   full_name,
                   email,
                   password)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
update users
set full_name = coalesce(sqlc.narg('full_name'), full_name),
    email     = coalesce(sqlc.narg('email'), email),
    password  = coalesce(sqlc.narg('password'), password)
    where username = sqlc.arg(username) RETURNING *;
