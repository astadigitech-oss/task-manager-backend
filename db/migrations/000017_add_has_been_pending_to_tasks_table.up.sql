-- Add has_been_pending column to tasks table
ALTER TABLE `tasks` ADD COLUMN `has_been_pending` BOOLEAN DEFAULT FALSE;
