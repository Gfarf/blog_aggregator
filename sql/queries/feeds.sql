-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsFromURL :one
SELECT * FROM feeds WHERE url=$1;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id=$1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at=$1, last_fetched_at=$2 WHERE id=$3;

-- name: GetNextFeedToFetch :one
SELECT id FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST;
