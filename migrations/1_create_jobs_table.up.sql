CREATE TABLE IF NOT EXISTS jobs  (
    id UUID PRIMARY KEY,
    name VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);