create table if not exists permissions(
    id serial primary key,
    user_type varchar check ("user_type" in('superadmin','user')) not null,
    resource varchar not null,
    action varchar not null,
    unique(user_type,resource,action)
);

--- Admin permissions
--Categories
insert into permissions(user_type,resource,action)
values('superadmin','categories','create');
insert into permissions(user_type,resource,action)
values('superadmin','categories','update');
insert into permissions(user_type,resource,action)
values('superadmin','categories','delete');

--Users
insert into permissions(user_type,resource,action)
values('superadmin','users','create');
insert into permissions(user_type,resource,action)
values('superadmin','users','update');
insert into permissions(user_type,resource,action)
values('superadmin','users','delete');

--Comments
insert into permissions(user_type,resource,action)
values('superadmin','comments','delete');

--Posts
insert into permissions(user_type,resource,action)
values('superadmin','posts','delete');




--- User permissions

--Likes
insert into permissions(user_type,resource,action)
values('user','likes','create');
insert into permissions(user_type,resource,action)
values('user','likes','get');

--Comments
insert into permissions(user_type,resource,action)
values('user','comments','create');
insert into permissions(user_type,resource,action)
values('user','comments','update');
insert into permissions(user_type,resource,action)
values('user','comments','delete');

-- Posts
insert into permissions(user_type,resource,action)
values('user','posts','create');
insert into permissions(user_type,resource,action)
values('user','posts','update');
insert into permissions(user_type,resource,action)
values('user','posts','delete');

--Auth
insert into permissions(user_type,resource,action)
values('user','auth','update-password');
insert into permissions(user_type,resource,action)
values('user','users','update');
insert into permissions(user_type,resource,action)
values('user','users','delete');






