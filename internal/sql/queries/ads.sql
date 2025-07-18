-- name: CreateAdvertisement :one
INSERT INTO advertisements(id, title, description, image_address, price, created_at, updated_at, user_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) 
RETURNING *;