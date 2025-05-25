-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING *
)
SELECT 
    ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
    u.name as user_name,
    f.name as feed_name
FROM inserted_feed_follow ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, 
       u.name as user_name, 
       f.name as feed_name
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;

-- name: GetFeedFollowByUserAndFeed :one
SELECT * FROM feed_follows 
WHERE user_id = $1 AND feed_id = $2;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows 
WHERE id = $1;

-- name: DeleteFeedFollowByUserAndFeedURL :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_id = (
    SELECT id FROM feeds WHERE url = $2
);