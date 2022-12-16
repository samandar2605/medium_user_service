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