-- Create workspaces table
CREATE TABLE `workspaces` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `description` longtext,
  `created_by` bigint(20) unsigned DEFAULT NULL,
  `color` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workspaces_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create workspace_users table for many-to-many relationship
CREATE TABLE `workspace_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `role_in_workspace` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_workspace_users_user` (`user_id`),
  KEY `fk_workspaces_members` (`workspace_id`),
  CONSTRAINT `fk_workspace_users_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_members` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
