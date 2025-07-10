-- Drop existing tables if they exist (for reruns during development)
DROP TABLE IF EXISTS user_intersections;
DROP TABLE IF EXISTS users;

-- Create users table with integer primary key
CREATE TABLE users (
    uuid SERIAL PRIMARY KEY,                    -- Auto-incrementing integer ID
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create user_intersections table
CREATE TABLE user_intersections (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    intersection_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, intersection_id)
);

-- Insert test data into users
INSERT INTO users (name, email, password, is_admin)
VALUES
    ('Alice Smith', 'alice@example.com', 'password123', false),
    ('Bob Johnson', 'bob@example.com', 'securepass', true),
    ('Charlie Lee', 'charlie@example.com', 'hunter2', false);
