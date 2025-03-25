package maa

import (
	"context"
	"maa-backend/internal/model/dto"

	"github.com/gogf/gf/v2/frame/g"
)

type Hook interface {
	Do(message int, details string, arg interface{})
}

type TaskCallback struct {
	Instance *Maa
}

//回调
func (t *TaskCallback) Do(message int, details string, arg interface{}) {
	g.Log().Debug(context.TODO(), "CallBack--->>>", message, details, arg)
	if dto.Filterd(message) {
		return
	}
	t.Instance.AppendLog(*dto.NewMaalogFromDetail(details, message))

}
