create table if not exists permissions(
    id serial primary key,
    user_type varchar check ("user_type" in('superadmin','user')) not null,
    resource varchar not null,
    action varchar not null,
    unique(user_type,resource,action)
);

-- Admin
insert into permissions(user_type,resource,action)
values('superadmin','categories','create');
insert into permissions(user_type,resource,action)
values('superadmin','categories','update');
insert into permissions(user_type,resource,action)
values('superadmin','categories','delete');
insert into permissions(user_type,resource,action)
values('superadmin','users','create');
insert into permissions(user_type,resource,action)
values('superadmin','users','update');
insert into permissions(user_type,resource,action)
values('superadmin','users','delete');
insert into permissions(user_type,resource,action)
values('superadmin','comments','delete');
insert into permissions(user_type,resource,action)
values('superadmin','posts','delete');

-- User
insert into permissions(user_type,resource,action)
values('user','categories','create');
insert into permissions(user_type,resource,action)
values('user','categories','update');
insert into permissions(user_type,resource,action)
values('user','categories','delete');
insert into permissions(user_type,resource,action)
values('user','likes','likes');
insert into permissions(user_type,resource,action)
values('user','likes','likes/user-post');
insert into permissions(user_type,resource,action)
values('user','comments','create');
insert into permissions(user_type,resource,action)
values('user','comments','update');
insert into permissions(user_type,resource,action)
values('user','comments','delete');
insert into permissions(user_type,resource,action)
values('user','posts','create');
insert into permissions(user_type,resource,action)
values('user','posts','update');
insert into permissions(user_type,resource,action)
values('user','posts','delete');
insert into permissions(user_type,resource,action)
values('user','auth','/auth/update-password');






