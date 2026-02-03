ALTER TABLE `task_files`
ADD COLUMN `file_data` LONGBLOB NOT NULL AFTER `filename`,
DROP COLUMN `url`;