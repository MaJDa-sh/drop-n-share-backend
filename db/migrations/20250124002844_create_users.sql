-- +goose Up
-- +goose StatementBegin
create table users (
    id bigint primary key generated by default as identity,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    created_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
