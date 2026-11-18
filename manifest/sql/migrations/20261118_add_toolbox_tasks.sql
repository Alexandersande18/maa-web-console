SET @toolbox_migration = IF(
  EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'cron_job' AND column_name = 'task_activity'
  ),
  'SELECT 1',
  'ALTER TABLE cron_job ADD COLUMN task_activity LONGTEXT NULL'
);
PREPARE toolbox_migration FROM @toolbox_migration;
EXECUTE toolbox_migration;
DEALLOCATE PREPARE toolbox_migration;

SET @toolbox_migration = IF(
  EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'cron_job' AND column_name = 'task_reclamation'
  ),
  'SELECT 1',
  'ALTER TABLE cron_job ADD COLUMN task_reclamation LONGTEXT NULL'
);
PREPARE toolbox_migration FROM @toolbox_migration;
EXECUTE toolbox_migration;
DEALLOCATE PREPARE toolbox_migration;
