CREATE TABLE IF NOT EXISTS `persons` (
	`id` VARCHAR(36) NOT NULL PRIMARY KEY,
  	`name` VARCHAR(255) NOT NULL,
  	`email` VARCHAR(255),
  	`created_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  	`updated_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
