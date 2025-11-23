-- +migrate Up
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    user_id_token VARCHAR(255) NOT NULL COMMENT 'Cognito User Sub',
    auth_type ENUM('cognito', 'none') NOT NULL DEFAULT 'cognito',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS users;
