create table if not exists permissions(
    id serial primary key,
    user_type varchar check ("user_type" in('superadmin','user')) not null,
    resource varchar not null,
    action varchar not null,
    unique(user_type,resource,action)
);