package dao_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"maa-backend/internal/dao"
	"maa-backend/internal/model/do"
	"maa-backend/internal/model/dto"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

// Exercise the real GoFrame field filtering and JSON column conversion. All test
// data is rolled back, and no scheduler or game task is started.
func TestRoguelikeTemplateDatabaseRoundTrip(t *testing.T) {
	if os.Getenv("MAA_DB_SMOKE") != "1" {
		t.Skip("set MAA_DB_SMOKE=1 to test against the configured database")
	}
	ctx := context.Background()
	if err := dao.CheckCronJobSchema(ctx); err != nil {
		t.Fatal(err)
	}
	rollback := errors.New("rollback smoke test")
	err := dao.CronJob.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(dao.CronJob.Table()).Data(do.CronJob{
			Name:       fmt.Sprintf("roguelike-smoke-%d", time.Now().UnixNano()),
			Username:   "schema-smoke-test",
			IsTemplate: 1,
			Hour:       0,
			Minute:     0,
			TaskRoguelike: &dto.TaskRoguelike{
				Enable: true, Theme: "BlackFlow", StartsCount: 8,
				BlackflowStartAnimals: []string{"swaddled_cat"},
			},
			TaskActivity: &dto.TaskActivity{Enable: true, ActivityID: "spa-regulations"},
			TaskReclamation: &dto.TaskReclamation{
				Enable: true, Theme: dto.ReclamationTales, Mode: 1,
				ToolsToCraft: []string{"荧光棒"}, NumCraftBatches: 3,
			},
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		var job dto.CronJob
		if err := tx.Model(dao.CronJob.Table()).WherePri(id).Scan(&job); err != nil {
			return err
		}
		if job.IsTemplate != 1 || job.TaskRoguelike == nil || !job.TaskRoguelike.Enable ||
			job.TaskRoguelike.Theme != "BlackFlow" || job.TaskRoguelike.StartsCount != 8 ||
			len(job.TaskRoguelike.BlackflowStartAnimals) != 1 || job.TaskRoguelike.BlackflowStartAnimals[0] != "swaddled_cat" {
			return fmt.Errorf("roguelike template configuration was not persisted: %+v", job.TaskRoguelike)
		}
		if job.TaskActivity == nil || !job.TaskActivity.Enable || job.TaskActivity.ActivityID != "spa-regulations" {
			return fmt.Errorf("activity configuration was not persisted: %+v", job.TaskActivity)
		}
		if job.TaskReclamation == nil || !job.TaskReclamation.Enable || job.TaskReclamation.Mode != 1 ||
			job.TaskReclamation.NumCraftBatches != 3 || len(job.TaskReclamation.ToolsToCraft) != 1 ||
			job.TaskReclamation.ToolsToCraft[0] != "荧光棒" {
			return fmt.Errorf("reclamation configuration was not persisted: %+v", job.TaskReclamation)
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("template round trip failed: %v", err)
	}
}
