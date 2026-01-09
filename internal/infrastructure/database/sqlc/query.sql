-- name: CreateUser :one
INSERT INTO users (
    id,
    kratos_id,
    name,
    highlighted_photo_id,
    locale,
    description,
    locality,
    administrative_area,
    country_code,
    postal_code,
    geom,
    first_name,
    last_name,
    email,
    has_set_location
) VALUES (
    sqlc.arg(id), sqlc.arg(kratos_id), sqlc.arg(name), sqlc.arg(highlighted_photo_id), sqlc.arg(locale), sqlc.arg(description), sqlc.arg(locality), sqlc.arg(administrative_area), sqlc.arg(country_code), sqlc.arg(postal_code), ST_GeomFromEWKB(sqlc.arg(geom)), sqlc.arg(first_name), sqlc.arg(last_name), sqlc.arg(email), sqlc.arg(has_set_location)
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    name = sqlc.arg(name),
    email = sqlc.arg(email),
    first_name = sqlc.arg(first_name),
    last_name = sqlc.arg(last_name),
    description = sqlc.arg(description),
    locality = sqlc.arg(locality),
    administrative_area = sqlc.arg(administrative_area),
    country_code = sqlc.arg(country_code),
    postal_code = sqlc.arg(postal_code),
    geom = ST_GeomFromEWKB(sqlc.arg(geom)),
    has_set_location = sqlc.arg(has_set_location),
    highlighted_photo_id = sqlc.arg(highlighted_photo_id),
    locale = sqlc.arg(locale)
WHERE id = sqlc.arg(id)
RETURNING *; 

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByKratosID :one
SELECT * FROM users WHERE kratos_id = $1;

-- name: CreateRoute :exec
INSERT INTO routes (
    id,
    user_id,
    name,
    description,
    highlighted_photo_id,
    distance,
    duration,
    elevation_gain,
    elevation_loss,
    path_geom,
    bbox,
    first_point,
    last_point,
    visibility
) VALUES (
    sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(name), sqlc.arg(description), sqlc.arg(highlighted_photo_id), sqlc.arg(distance), sqlc.arg(duration), sqlc.arg(elevation_gain), sqlc.arg(elevation_loss), ST_GeomFromEWKB(sqlc.arg(path_geom)), ST_GeomFromEWKB(sqlc.arg(bbox)), ST_GeomFromEWKB(sqlc.arg(first_point)), ST_GeomFromEWKB(sqlc.arg(last_point)), sqlc.arg(visibility)
);

-- name: UpdateRoute :exec
UPDATE routes SET
    name = sqlc.arg(name),
    description = sqlc.arg(description),
    highlighted_photo_id = sqlc.arg(highlighted_photo_id),
    distance = sqlc.arg(distance),
    duration = sqlc.arg(duration),
    elevation_gain = sqlc.arg(elevation_gain),
    elevation_loss = sqlc.arg(elevation_loss),
    path_geom = ST_GeomFromEWKB(sqlc.arg(path_geom)),
    bbox = ST_GeomFromEWKB(sqlc.arg(bbox)),
    first_point = ST_GeomFromEWKB(sqlc.arg(first_point)),
    last_point = ST_GeomFromEWKB(sqlc.arg(last_point)),
    visibility = sqlc.arg(visibility)
WHERE id = sqlc.arg(id);

-- name: GetRouteByID :one
SELECT * FROM routes WHERE id = $1;

-- name: GetRoutesByUserID :many
SELECT * FROM routes WHERE user_id = $1;

-- name: CountRoutesByUserID :one
SELECT COUNT(*) FROM routes WHERE user_id = $1;

-- name: DeleteRoute :one
DELETE FROM routes WHERE id = $1 RETURNING id;

-- name: CreateCoursePoint :exec
INSERT INTO course_points (
    id,
    route_id,
    step_order,
    seg_dist_m,
    cum_dist_m,
    duration,
    instruction,
    road_name,
    maneuver_type,
    modifier,
    location,
    bearing_before,
    bearing_after
) VALUES (
    sqlc.arg(id), sqlc.arg(route_id), sqlc.arg(step_order), sqlc.arg(seg_dist_m), sqlc.arg(cum_dist_m), sqlc.arg(duration), sqlc.arg(instruction), sqlc.arg(road_name), sqlc.arg(maneuver_type), sqlc.arg(modifier), ST_GeomFromEWKB(sqlc.arg(location)), sqlc.arg(bearing_before), sqlc.arg(bearing_after)
);

-- name: GetCoursePointsByRouteID :many
SELECT * FROM course_points WHERE route_id = $1 ORDER BY step_order ASC;

-- name: DeleteCoursePoint :exec
DELETE FROM course_points WHERE id = $1;

-- name: DeleteCoursePointsByRouteID :exec
DELETE FROM course_points WHERE route_id = $1;

-- name: CreateWaypoint :exec
INSERT INTO waypoints (
    id,
    route_id,
    location
) VALUES (
    sqlc.arg(id), sqlc.arg(route_id), ST_GeomFromEWKB(sqlc.arg(location))
);

-- name: GetWaypointsByRouteID :many
SELECT * FROM waypoints WHERE route_id = $1 ORDER BY id ASC;

-- name: DeleteWaypoint :exec
DELETE FROM waypoints WHERE id = $1;

-- name: DeleteWaypointsByRouteID :exec
DELETE FROM waypoints WHERE route_id = $1;
