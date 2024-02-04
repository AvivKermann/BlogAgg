-- name: CreateFeedFollow :one
INSERT INTO feed_follow(id, feed_id, user_id,created_at,updated_at)
VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetFeedFollowByID :one
SELECT * FROM feed_follow
WHERE feed_follow.id = $1;

-- name: DeleteFeedFollowByID :exec
DELETE FROM feed_follow
WHERE id = $1;

-- name: GetFeedFollowByUserID :many
SELECT * FROM feed_follow
WHERE user_id = $1;