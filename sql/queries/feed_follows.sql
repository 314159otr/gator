-- name: CreateFeedFollow :one
with inserted as (
	insert into feed_follows(id, user_id, feed_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5)
	returning *
)
select inserted.*, feeds.name as feed_name, users.name as user_name
from inserted
join users on users.id = inserted.user_id
join feeds on feeds.id = inserted.feed_id;

-- name: GetFeedFollowsForUser :many
select
	feed_follows.*,
	feeds.name as feed_name,
	users.name as user_name
from feed_follows
join feeds on feeds.id = feed_follows.feed_id
join users on users.id = feed_follows.user_id
where users.id = $1;

