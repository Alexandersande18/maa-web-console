package dao

import (
	"context"
	"fmt"
)

// CheckCronJobSchema prevents GoFrame from silently filtering out task data when
// an existing database has not been upgraded to the current schema.
func CheckCronJobSchema(ctx context.Context) error {
	fields, err := CronJob.DB().TableFields(ctx, CronJob.Table())
	if err != nil {
		return fmt.Errorf("check cron_job schema: %w", err)
	}
	if _, ok := fields[CronJob.Columns().TaskRoguelike]; !ok {
		return fmt.Errorf("cron_job.task_roguelike is missing: run manifest/sql/migrations/20261107_add_task_roguelike.sql and restart the backend")
	}
	columns := CronJob.Columns()
	for _, name := range []string{columns.TaskActivity, columns.TaskReclamation} {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("cron_job.%s is missing: run manifest/sql/migrations/20261118_add_toolbox_tasks.sql and restart the backend", name)
		}
	}
	return nil
}
