USE mydatabase;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE -- Added 'is_active' column
);

-- Insert a few sample users for testing
INSERT INTO users (username, email, is_active) VALUES 
('alice_tester', 'alice@example.com', TRUE),
('bob_tester', 'bob@example.com', FALSE);