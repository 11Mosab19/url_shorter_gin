--enums

CREATE TYPE "roles" AS ENUM ('admin','user');
CREATE TYPE "status" AS ENUM ('active','disabled','deleted');

--tables

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    full_name TEXT NOT NULL,
    role roles DEFAULT 'user'
);

CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'active',
    hashed_password TEXT,   
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_clicks INT DEFAULT 0 NOT NULL CHECK(total_clicks >= 0 ),
    user_id INT REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    device_type TEXT NOT NULL,
);


CREATE INDEX idx_url_id_analytics ON analytics(url_id);
CREATE INDEX idx_user_id_urls ON urls(user_id);
CREATE INDEX idx_original_url_urls ON urls(original_url);
