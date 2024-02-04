// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPosts = `-- name: CreatePosts :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES($1,$2,$3,$4,$5,$6, $7, $8)
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostsParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePosts(ctx context.Context, arg CreatePostsParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPosts,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsByUserID = `-- name: GetPostsByUserID :many
SELECT posts.id, posts.created_at, posts.updated_at, title, posts.url, description, published_at, feed_id, feed.id, feed.created_at, feed.updated_at, feed.name, feed.url, user_id, last_fetched_at, users.id, users.created_at, users.updated_at, users.name, api_key FROM posts
INNER JOIN feed ON feed.id = posts.feed_id
INNER JOIN users ON users.id = feed.user_id
WHERE users.id = $1
ORDER BY posts.created_at DESC
LIMIT $2
`

type GetPostsByUserIDParams struct {
	ID    uuid.UUID
	Limit int32
}

type GetPostsByUserIDRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Url           string
	Description   string
	PublishedAt   time.Time
	FeedID        uuid.UUID
	ID_2          uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	Name          string
	Url_2         string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	ID_3          uuid.UUID
	CreatedAt_3   time.Time
	UpdatedAt_3   time.Time
	Name_2        string
	ApiKey        string
}

func (q *Queries) GetPostsByUserID(ctx context.Context, arg GetPostsByUserIDParams) ([]GetPostsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUserID, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByUserIDRow
	for rows.Next() {
		var i GetPostsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
			&i.Url_2,
			&i.UserID,
			&i.LastFetchedAt,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.Name_2,
			&i.ApiKey,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
