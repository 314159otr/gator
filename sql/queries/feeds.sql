-- name: CreateFeed :one
insert into feeds (id, user_id, created_at, updated_at, name, url)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedByURL :one
select * from feeds
where url = $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $1,
	updated_at      = $2
where id = $3;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at asc nulls first
limit 1;
