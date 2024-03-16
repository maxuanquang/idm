-- Create account table
CREATE TABLE IF NOT EXISTS `account` (
    `account_id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `account_name` VARCHAR(50) UNIQUE NOT NULL
);

-- Create account_password table
CREATE TABLE IF NOT EXISTS `account_password` (
    `of_account_id` BIGINT UNSIGNED PRIMARY KEY,
    `hashed` VARCHAR(128) NOT NULL,
    FOREIGN KEY (`of_account_id`) REFERENCES `account` (`account_id`)
);

-- Create download_task table
CREATE TABLE IF NOT EXISTS `download_task` (
    `download_task_id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `of_account_id` BIGINT UNSIGNED NOT NULL,
    `download_type` SMALLINT NOT NULL,
    `download_url` TEXT NOT NULL,
    `download_status` SMALLINT NOT NULL,
    `metadata` TEXT NOT NULL,
    FOREIGN KEY (`of_account_id`) REFERENCES `account` (`account_id`)
);

-- Create token_public_key table
CREATE TABLE IF NOT EXISTS `token_public_key` (
    `token_public_key_id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `token_public_key_value` VARBINARY(4096) NOT NULL
);
