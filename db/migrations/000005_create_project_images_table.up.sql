-- Create project_images table
CREATE TABLE `project_images` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned DEFAULT NULL,
  `url` longtext,
  `uploaded_by` bigint(20) unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_project_images_deleted_at` (`deleted_at`),
  KEY `fk_projects_images` (`project_id`),
  CONSTRAINT `fk_projects_images` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
