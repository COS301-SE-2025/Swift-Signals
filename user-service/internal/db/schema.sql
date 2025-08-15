DROP TABLE IF EXISTS user_intersections;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    uuid UUID PRIMARY KEY, 
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_intersections (
    user_id UUID REFERENCES users(uuid) ON DELETE CASCADE,
    intersection_id UUID NOT NULL,
    PRIMARY KEY (user_id, intersection_id)
);

INSERT INTO users (uuid, name, email, password, is_admin)
VALUES
    ('9b9b1c5c-2e57-4e18-a15c-e3219be9dc01', 'Alice Smith', 'alice@example.com', 'password123', false),
    ('2f1a9b99-bdc2-44ce-9f0c-d3903f7b9eb1', 'Bob Johnson', 'bob@example.com', 'securepass', false),
    ('30d6cbb9-0f3f-4f9b-bb2a-52c9d5e7231e', 'Charlie Lee', 'charlie@example.com', 'hunter2', false);

