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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    name = $2,
    email = $3,
    first_name = $4,
    last_name = $5,
    description = $6,
    locality = $7,
    administrative_area = $8,
    country_code = $9,
    postal_code = $10,
    geom = $11,
    has_set_location = $12,
    highlighted_photo_id = $13,
    locale = $14
WHERE id = $1
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: UpdateRoute :exec
UPDATE routes SET
    name = $2,
    description = $3,
    highlighted_photo_id = $4,
    distance = $5,
    duration = $6,
    elevation_gain = $7,
    elevation_loss = $8,
    path_geom = $9,
    bbox = $10,
    first_point = $11,
    last_point = $12,
    visibility = $13
WHERE id = $1;

-- name: GetRouteByID :one
SELECT * FROM routes WHERE id = $1;

-- name: GetRoutesByUserID :many
SELECT * FROM routes WHERE user_id = $1;

-- name: DeleteRoute :exec
DELETE FROM routes WHERE id = $1;

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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: GetCoursePointsByRouteID :many
SELECT * FROM course_points WHERE route_id = $1 ORDER BY step_order ASC;

-- name: DeleteCoursePoint :exec
DELETE FROM course_points WHERE id = $1;

-- name: CreateWaypoint :exec
INSERT INTO waypoints (
    id,
    route_id,
    location
) VALUES (
    $1, $2, $3
);

-- name: GetWaypointsByRouteID :many
SELECT * FROM waypoints WHERE route_id = $1 ORDER BY id ASC;

-- name: DeleteWaypoint :exec
DELETE FROM waypoints WHERE id = $1;


