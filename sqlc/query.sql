-- name: CreateUser :one
INSERT INTO users (
    name, email, geom, first_name, last_name, has_set_location
) VALUES (
    $1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;