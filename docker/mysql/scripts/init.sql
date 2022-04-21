CREATE USER 'otus'@'%' IDENTIFIED BY 'otus';
CREATE DATABASE otus;
GRANT ALL PRIVILEGES ON *.* TO 'otus'@'%';
CREATE USER 'monitor'@'%' IDENTIFIED BY 'monitorpassword';
GRANT ALL PRIVILEGES ON *.* TO 'monitor'@'%';

USE otus;
CREATE TABLE messages (
   id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
   author_id int NOT NULL,
   chat_id VARCHAR(255) NOT NULL,
   message VARCHAR(4096),
   shard_factor VARCHAR(2),
   created_at TIMESTAMP DEFAULT (NOW()),
   is_read BOOLEAN DEFAULT (false )
);
CREATE INDEX messages_chat_id_idx ON messages (shard_factor, chat_id);