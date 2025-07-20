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

-- name: GetAdvertisements :many
SELECT 
  ads.title, 
  ads.description, 
  ads.image_address, 
  ads.price, 
  ads.user_id, 
  users.login AS author_login
FROM advertisements AS ads
JOIN users ON users.id = ads.user_id
WHERE 
  (sqlc.arg(min_price)::int IS NULL OR ads.price >= sqlc.arg(min_price))
  AND (sqlc.arg(max_price)::int IS NULL OR ads.price <= sqlc.arg(max_price))
ORDER BY
  CASE WHEN sqlc.arg(order_by) = 'price'      AND sqlc.arg(order_dir) = 'asc'  THEN ads.price     END ASC,
  CASE WHEN sqlc.arg(order_by) = 'price'      AND sqlc.arg(order_dir) = 'desc' THEN ads.price     END DESC,
  CASE WHEN sqlc.arg(order_by) = 'created_at' AND sqlc.arg(order_dir) = 'asc'  THEN ads.created_at END ASC,
  CASE WHEN sqlc.arg(order_by) = 'created_at' AND sqlc.arg(order_dir) = 'desc' THEN ads.created_at END DESC,
  ads.created_at DESC
LIMIT $1 OFFSET $2;