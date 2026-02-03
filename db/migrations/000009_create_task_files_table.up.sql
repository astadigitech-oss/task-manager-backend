SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `task_files` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `task_id` bigint(20) unsigned NOT NULL,
    `filename` VARCHAR(255) NOT NULL,
    `file_data` LONGBLOB NOT NULL,
    `mime_type` VARCHAR(100) NOT NULL,
    `file_size` bigint(20) NOT NULL,
    `uploaded_by` bigint(20) unsigned NOT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    
    PRIMARY KEY (`id`),
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_uploaded_by` (`uploaded_by`),
    
    CONSTRAINT `fk_task_files_task` 
        FOREIGN KEY (`task_id`) 
        REFERENCES `tasks`(`id`) 
        ON DELETE CASCADE,
        
    CONSTRAINT `fk_task_files_user` 
        FOREIGN KEY (`uploaded_by`) 
        REFERENCES `users`(`id`) 
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;