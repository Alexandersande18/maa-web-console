-- Run against the configured database before restarting the backend.
-- Preserve existing jobs by adding only the missing nullable column.
SET @roguelike_migration = IF(
  EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'cron_job' AND column_name = 'task_roguelike'
  ),
  'SELECT 1',
  'ALTER TABLE cron_job ADD COLUMN task_roguelike LONGTEXT NULL AFTER task_fight'
);
PREPARE roguelike_migration FROM @roguelike_migration;
EXECUTE roguelike_migration;
DEALLOCATE PREPARE roguelike_migration;
