-- Create tasks table
CREATE TABLE `tasks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned DEFAULT NULL,
  `title` longtext,
  `description` longtext,
  `status` longtext,
  `priority` longtext,
  `start_date` datetime(3) DEFAULT NULL,
  `due_date` datetime(3) DEFAULT NULL,
  `finished_at` datetime(3) DEFAULT NULL,
  `notes` longtext,
  `overdue_duration` bigint(20) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tasks_deleted_at` (`deleted_at`),
  KEY `fk_projects_tasks` (`project_id`),
  CONSTRAINT `fk_projects_tasks` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create task_users table for model TaskUser
CREATE TABLE `task_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `role_in_task` longtext,
  `assigned_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_tasks_members` (`task_id`),
  KEY `fk_users_tasks` (`user_id`),
  CONSTRAINT `fk_tasks_members` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_users_tasks` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
