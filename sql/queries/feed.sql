-- name: CreateFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetFeedById :one
SELECT * FROM feed 
WHERE id = $1;

-- name: GetAllFeeds :many
SELECT * FROM feed;