package croner

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"maa-backend/internal/dao"
	"maa-backend/internal/model/do"
	"maa-backend/internal/model/dto"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

type sCroner struct {
	mu    sync.Mutex
	ctx   gctx.Ctx
	gcron *gcron.Cron
}

func init() {
	sc := &sCroner{ctx: gctx.New(), gcron: gcron.New()}
	service.RegisterCroner(sc)
}

// ExecCronJob 执行具体任务
func (s *sCroner) ExecCronJob(job *dto.CronJob) error {
	uid := job.UserName

	// 根据任务类型执行具体操作
	g.Log().Infof(s.ctx, "Executing job: %+v, User: %+v", job.Name, uid)

	taskList, err := job.TaskSelection.Build(job.IsTemplate == 0)
	if err != nil {
		return err
	}

	if err := service.Maasst().LinkStart(uid, taskList); err != nil {
		return err
	}
	// Rejected duplicate jobs must not overwrite the last accepted execution.
	if _, err := dao.CronJob.Ctx(s.ctx).Where("id", gconv.Int(job.Id)).Update(do.CronJob{
		LastExecutedAt: gconv.GTime(time.Now()),
	}); err != nil {
		g.Log().Errorf(s.ctx, "update last_executed_at failed: %v", err)
	}
	return nil
}

// SetupCronJobs 初始化定时任务
func (s *sCroner) SetupCronJobs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 清除所有现有任务
	s.gcron.Stop()
	s.gcron = gcron.New()

	// 查询有效任务（保持不变）
	cronJobList := make([]*dto.CronJob, 0)
	if err := dao.CronJob.Ctx(s.ctx).
		Where("deleted_at IS NULL").
		Where("is_template", 0).
		Scan(&cronJobList); err != nil {
		g.Log().Error(s.ctx, "SetupCronJobs query failed:", err)
		return
	}

	// 设置定时任务
	for _, cronJob := range cronJobList {
		// 解析基础时间参数
		baseHour := gconv.Int(cronJob.Hour)
		baseMinute := gconv.Int(cronJob.Minute)

		// 生成随机时间偏移（5分钟以内）
		minuteOffset := rand.Intn(11) - 5 // -5 ~ +5 分钟
		adjustedMinute := baseMinute + minuteOffset
		adjustedHour := baseHour

		// 处理分钟溢出
		if adjustedMinute < 0 {
			adjustedHour--
			adjustedMinute += 60
		} else if adjustedMinute >= 60 {
			adjustedHour++
			adjustedMinute -= 60
		}

		// 处理小时溢出
		adjustedHour = (adjustedHour + 24) % 24

		// 生成随机秒（0-59）
		randSecond := rand.Intn(60)

		// 构造完整cron表达式
		/*
			cron表达式格式说明：
			[秒] [分] [时] [日] [月] [周]
			示例：每天10点30分15秒执行 -> "15 30 10 * * *"
		*/
		cronPattern := fmt.Sprintf("%d %d %d * * *",
			randSecond,     // 秒（0-59）
			adjustedMinute, // 分（0-59）
			adjustedHour,   // 时（0-23）
		)

		// 克隆job对象用于闭包
		currentJob := cronJob

		// 添加任务到cron
		entry, err := s.gcron.Add(s.ctx, cronPattern, func(ctx context.Context) {
			if err := s.ExecCronJob(currentJob); err != nil {
				g.Log().Errorf(ctx, "执行定时任务失败 | 名称:%s | 错误:%v", currentJob.Name, err)
			}
		}, cronJob.Name)

		if err != nil {
			g.Log().Errorf(s.ctx, "添加定时任务失败 | 名称:%s | 表达式:%s | 错误:%v",
				cronJob.Name, cronPattern, err)
			continue
		}

		// 记录调试日志
		g.Log().Infof(s.ctx, `[定时任务已设置]用户：%s, 名称:%s, 原始时间:%02d:%02d, 实际时间:%02d:%02d:%02d, 表达式:%s, 任务ID : %s`,
			cronJob.UserName,
			cronJob.Name,
			baseHour, baseMinute,
			adjustedHour, adjustedMinute, randSecond,
			cronPattern,
			entry.Name)
	}
}
