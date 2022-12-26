CREATE TABLE IF NOT EXISTS "users"(
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(30) NOT NULL,
    "last_name" VARCHAR(30) NOT NULL,
    "phone_number" VARCHAR(20) UNIQUE,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "gender" VARCHAR(10) CHECK ("gender" IN('male', 'female')) DEFAULT 'male',
    "password" VARCHAR NOT NULL,
    "username" VARCHAR(30) UNIQUE,
    "profile_image_url" VARCHAR,
    "type" VARCHAR(255) CHECK ("type" IN('superadmin', 'user')) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

insert into users(first_name,last_name,phone_number,email,gender,password,username,profile_image_url,type)
values('Samandar','Tukhtayev','+998 77 777 77 77','ukan265@gmail.com','male','$2a$10$E.9Xm//sp2WhqwEVyT.zEuzPGbHwagSV4LuHEHFhe5SBVlrlv8yA2','AdminUkan','samandaradmin.png','superadmin');