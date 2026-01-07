-- Create error_logs table
CREATE TABLE `error_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `action` longtext,
  `table_name` longtext,
  `item_id` bigint(20) unsigned DEFAULT NULL,
  `error_msg` longtext,
  `stack_trace` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_error_logs_user` (`user_id`),
  CONSTRAINT `fk_error_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
