CREATE TABLE IF NOT EXISTS auth_user (
    id bigint generated always as identity primary key,
    username varchar(128) not null unique,
    email varchar(255) not null unique,
    role varchar(128) not null,
    password varchar(255) not null,
    created_at timestamptz not null default current_timestamp,
    last_login timestamptz not null default current_timestamp,
    attr jsonb not null default '{}'::jsonb
);
CREATE TABLE IF NOT EXISTS auth_event (
    id bigint generated always as identity primary key,
    user_id bigint not null,
    action varchar(128) not null,
    app varchar(128) not null,
    detail jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default current_timestamp
);