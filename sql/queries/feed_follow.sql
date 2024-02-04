-- name: CreateFeedFollow :one
INSERT INTO feed_follow(id, created_at, updated_at, user_id, feed_id)
VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetFeedFollowByID :one
SELECT * FROM feed_follow
WHERE feed_follow.id = $1;

-- name: DeleteFeedFollowByID :exec
DELETE FROM feed_follow
WHERE id = $1;