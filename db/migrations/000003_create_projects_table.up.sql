-- Create projects table
CREATE TABLE `projects` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `description` longtext,
  `created_by` bigint(20) unsigned DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_projects_deleted_at` (`deleted_at`),
  KEY `fk_workspaces_projects` (`workspace_id`),
  CONSTRAINT `fk_workspaces_projects` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create project_users table for model ProjectUser
CREATE TABLE `project_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `role_in_project` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_project_users_user` (`user_id`),
  KEY `fk_projects_members` (`project_id`),
  CONSTRAINT `fk_project_users_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_projects_members` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
