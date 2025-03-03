package maa

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	log "github.com/sirupsen/logrus"
)

type Hook interface {
	Do(message int, details string, arg interface{})
}

type sampleHook struct {
}

func (sampleHook) Do(message int, details string, arg interface{}) {
	log.Infoln(message, details, arg)
}

type TaskCallback struct {
}

func (TaskCallback) Do(message int, details string, arg interface{}) {
	g.Log().Info(context.TODO(), "CallBack--->>>", message, details, arg)
}
