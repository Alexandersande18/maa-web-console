package maa

import (
	"context"
	"fmt"
	dto "maa-backend/internal/model/dto"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sMaasst struct {
	instances map[string]*Maa
}

func init() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	idToInstance := make(map[string]*Maa)
	for _, ast := range cfg.Asst {
		maa, err := New(ast)
		if err != nil {
			panic(err)
		}
		idToInstance[ast.ID] = maa
	}
    service.RegisterMaasst(&sMaasst{instances: idToInstance})
}

func (s *sMaasst) getInstance(id string) *Maa {
	if s.instances[id] == nil {
		return nil
	}
	return s.instances[id]
}

func (s *sMaasst) HasInstance(id string) bool {
	return s.instances[id] != nil
}

func (s *sMaasst) TakeScreenshot(id string) ([]byte, error) {
	if !s.HasInstance(id) {
		return nil, fmt.Errorf("instance not found")
	}
    maa := s.getInstance(id)
    return maa.GetImage()
}

func (s *sMaasst) LinkStart(id string, tasks []dto.Tasker) error {
    if !s.HasInstance(id) {
        return fmt.Errorf("instance not found")
    }
    maa := s.getInstance(id)
    for _, task := range tasks {
        tskID, err := maa.AsstAppendTask(task.GetType(), task)
        if err != nil {
            return err
        }
        g.Log().Debugf(context.TODO(), "Task %s appended with ID %d", task.GetType(), tskID)
    }
    if ok, err := maa.AsstStart(); !ok{
		return err
	}
    return nil
}

func (s *sMaasst) StopTask(id string) error {
	if !s.HasInstance(id) {
		return fmt.Errorf("instance not found")
	}
	maa := s.getInstance(id)

	if ok, err := maa.AsstStop(); !ok{
		return err
	}
	return nil
}

func (s *sMaasst) NewSSEClient(mid string, cid string) (*dto.SSEClientConn, error) {
	if !s.HasInstance(mid) {
        return nil, fmt.Errorf("instance not found")
    }
	maa := s.getInstance(mid)
	maa.NewClient(cid)
	
	return s.GetSSEClient(mid, cid)
}

func (s *sMaasst) DeleteSSEClient(mid string, cid string) error {
	if !s.HasInstance(mid) {
		return fmt.Errorf("instance not found")
	}
	maa := s.getInstance(mid)
	maa.DeleteClient(cid)
	return nil
}

func (s *sMaasst) GetSSEClient(mid string, cid string) (*dto.SSEClientConn, error) {
	if !s.HasInstance(mid) {
		return nil, fmt.Errorf("instance not found")
	}
	maa := s.getInstance(mid)
	return maa.GetClient(cid), nil
}