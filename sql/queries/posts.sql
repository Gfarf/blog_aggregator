-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts 
WHERE posts.feed_id IN 
(SELECT feed_follows.feed_id FROM feed_follows WHERE feed_follows.user_id=$1)
ORDER BY posts.published_at DESC 
LIMIT $2;

-- name: GetPostsUrls :many
SELECT url FROM posts;

-- name: GetPosts :many
SELECT * FROM posts LIMIT 10; 