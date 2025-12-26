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