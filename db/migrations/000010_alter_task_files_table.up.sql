ALTER TABLE `task_files` 
ADD COLUMN `url` VARCHAR(255) NOT NULL AFTER `filename`,
DROP COLUMN `file_data`;