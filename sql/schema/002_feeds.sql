-- +goose Up
create table feeds(
	id         uuid primary key,
	user_id    uuid not null references users(id) on delete cascade,
	created_at timestamp not null,
	updated_at timestamp not null,
	name       text not null,
	url        text not null unique
);

-- +goose Down
drop table feeds;
