CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age SMALLINT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    location geography(Point, 4326) NOT NULL,
    interests TEXT[] NOT NULL,
    preferences JSONB,
    last_active TIMESTAMP 
);  

CREATE INDEX idx_users_location ON users USING GIST(location);

CREATE INDEX idx_users_age ON users(age);
CREATE INDEX idx_users_gender ON users(gender);