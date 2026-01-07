-- Create task_images table
CREATE TABLE `task_images` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint(20) unsigned DEFAULT NULL,
  `url` longtext,
  `uploaded_by` bigint(20) unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_tasks_images` (`task_id`),
  CONSTRAINT `fk_tasks_images` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
