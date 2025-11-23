-- name: CreateUser :one
INSERT INTO users (
    name, email, total_trip_distance, total_trip_duration, total_trip_elevation_gain, geom, first_name, last_name, has_set_location
) VALUES (
    $1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8, $9, $10
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;