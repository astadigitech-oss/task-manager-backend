CREATE TABLE `attendance_images` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `attendance_id` BIGINT UNSIGNED NOT NULL,
    `url` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    CONSTRAINT `fk_attendance_images_attendance`
        FOREIGN KEY(`attendance_id`) 
        REFERENCES `attendances`(`id`)
        ON DELETE CASCADE
);
