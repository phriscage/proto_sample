CREATE TABLE IF NOT EXISTS `person_phone_numbers` (
	`id` VARCHAR(36) NOT NULL PRIMARY KEY,
	`person_id` VARCHAR(36) NOT NULL,
	`phone_number_id` VARCHAR(36) NOT NULL,
	`created_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

