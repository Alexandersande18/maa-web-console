package maa

import (
	"context"
	"maa-backend/internal/model/dto"

	"github.com/gogf/gf/v2/frame/g"
	log "github.com/sirupsen/logrus"
)

type Hook interface {
	Do(message int, details string, arg interface{})
}

type TaskCallback struct {
	Instance *Maa
}

// 回调
func (t *TaskCallback) Do(message int, details string, arg interface{}) {
	g.Log().Debug(context.TODO(), "CallBack--->>>", message, details, arg)
	if t.Instance != nil {
		t.Instance.recordTaskCallback(message, details)
	}
	if message == dto.TaskChainError || message == dto.SubTaskError || message == dto.InternalError || message == dto.InitFailed {
		fields := log.Fields{"code": message, "details": details}
		if t.Instance != nil && t.Instance.cfg != nil {
			fields["instance"] = t.Instance.cfg.ID
		}
		log.WithFields(fields).Error("MaaCore task callback error")
	}
	if dto.Filterd(message) {
		return
	}
	if t.Instance != nil {
		t.Instance.AppendLog(*dto.NewMaalogFromDetail(details, message))
	}

}
