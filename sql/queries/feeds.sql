-- name: CreateFeed :one
insert into feeds (id, user_id, created_at, updated_at, name, url)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- name: GetFeeds :many
select * from feeds;
