-- Create account table
CREATE TABLE IF NOT EXISTS `account` (
    `account_id` bigint unsigned PRIMARY KEY,
    `account_name` varchar(50) NOT NULL
);

-- Create account_password table
CREATE TABLE IF NOT EXISTS `account_password` (
    `of_account_id` bigint unsigned PRIMARY KEY,
    `hashed` varchar(128) NOT NULL,
    FOREIGN KEY (`of_account_id`) REFERENCES `account` (`account_id`)
);

-- Create download_task table
CREATE TABLE IF NOT EXISTS `download_task` (
    `download_task_id` bigint unsigned PRIMARY KEY,
    `of_account_id` bigint unsigned NOT NULL,
    `download_type` smallint NOT NULL,
    `download_url` text NOT NULL,
    `download_status` smallint NOT NULL,
    `metadata` text NOT NULL,
    FOREIGN KEY (`of_account_id`) REFERENCES `account` (`account_id`)
);
