-- name: AddFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
JOIN users ON users.id = feeds.user_id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedByURL :one
SELECT
    users.id AS user_id, feeds.id AS feed_id
FROM
    feeds
JOIN users ON users.id = feeds.user_id
WHERE feeds.url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
JOIN users
    ON users.id = feed_follows.user_id
JOIN feeds
    ON feeds.id = feed_follows.feed_id
WHERE users.name = $1;

-- name: DeleteFollowFeed :exec
DELETE FROM feed_follows
WHERE
    user_id = $1
    AND
    feed_id = $2;

-- name: MarkFeedFetched :one
UPDATE feeds
SET
    last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id=$1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
