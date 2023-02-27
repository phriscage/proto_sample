CREATE TABLE IF NOT EXISTS `books` (
	`id` VARCHAR(36) NOT NULL PRIMARY KEY,
	`name` VARCHAR(255) NOT NULL,
	`title` VARCHAR(255),
	`created_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	`updated_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
