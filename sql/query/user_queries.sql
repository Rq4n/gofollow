-- name: CreateNewUser :exec
INSERT INTO users (email, password)
VALUES ($1, $2);

-- name: GetUserByName :one
SELECT id, email, password FROM users
WHERE email = $1;
