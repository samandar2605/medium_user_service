create table if not exists permissions(
    id serial primary key,
    user_type varchar check ("user_type" in('superadmin','user')) not null,
    resource varchar not null,
    action varchar not null,
    unique(user_type,resource,action)
);

insert into permissions(user_type,resource,action)
values('superadmin','categories','create');

insert into permissions(user_type,resource,action)
values('superadmin','categories','update');

insert into permissions(user_type,resource,action)
values('superadmin','categories','delete');

insert into permissions(user_type,resource,action)
values('user','users','create');

insert into permissions(user_type,resource,action)
values('user','users','update');

insert into permissions(user_type,resource,action)
values('user','users','delete');