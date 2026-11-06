-- Empty schema only. Do not commit database dumps or user data.
-- Create a database named maa with utf8mb4 before importing this file.
CREATE TABLE IF NOT EXISTS cron_job (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(191) NOT NULL,
  username VARCHAR(191) NOT NULL,
  is_template TINYINT NOT NULL DEFAULT 0,
  hour TINYINT UNSIGNED NOT NULL DEFAULT 0,
  minute TINYINT UNSIGNED NOT NULL DEFAULT 0,
  task_recruit LONGTEXT NULL,
  task_infrast LONGTEXT NULL,
  task_fight LONGTEXT NULL,
  task_roguelike LONGTEXT NULL,
  task_mall LONGTEXT NULL,
  task_award LONGTEXT NULL,
  created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  last_executed_at DATETIME NULL,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_cron_job_name (name),
  KEY ix_cron_job_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
