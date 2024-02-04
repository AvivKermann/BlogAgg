-- name: CreatePosts :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES($1,$2,$3,$4,$5,$6, $7, $8)
RETURNING *;

-- name: GetPostsByUserID :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id FROM posts
INNER JOIN feed ON feed.id = posts.feed_id
INNER JOIN users ON users.id = feed.user_id
WHERE users.id = $1
ORDER BY posts.created_at DESC
LIMIT $2;