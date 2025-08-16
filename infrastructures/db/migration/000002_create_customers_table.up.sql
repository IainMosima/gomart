CREATE TABLE IF NOT EXISTS customers (
    user_id UUID PRIMARY KEY,
    phone_number VARCHAR(13) NOT NULL,
    user_name VARCHAR(25) NOT NULL,
    email VARCHAR(50) NOT NULL
);