package dto

import "fmt"

type TaskSelection struct {
	TaskRecruit     *TaskRecruit     `json:"taskRecruit" v:""`
	TaskInfrast     *TaskInfrast     `json:"taskInfrast" v:""`
	TaskFight       *TaskFight       `json:"taskFight" v:""`
	TaskRoguelike   *TaskRoguelike   `json:"taskRoguelike" v:""`
	TaskMall        *TaskMall        `json:"taskMall" v:""`
	TaskAward       *TaskAward       `json:"taskAward" v:""`
	TaskActivity    *TaskActivity    `json:"taskActivity" v:""`
	TaskReclamation *TaskReclamation `json:"taskReclamation" v:""`
}

func (s TaskSelection) Build(scheduled bool) ([]Tasker, error) {
	var tasks []Tasker
	if s.TaskRecruit != nil && s.TaskRecruit.Enable {
		tasks = append(tasks, s.TaskRecruit)
	}
	if s.TaskInfrast != nil && s.TaskInfrast.Enable {
		tasks = append(tasks, s.TaskInfrast)
	}
	if s.TaskFight != nil && s.TaskFight.Enable {
		tasks = append(tasks, s.TaskFight)
	}
	if s.TaskRoguelike != nil && s.TaskRoguelike.Enable {
		tasks = append(tasks, s.TaskRoguelike)
	}
	if s.TaskMall != nil && s.TaskMall.Enable {
		tasks = append(tasks, s.TaskMall)
	}
	if s.TaskAward != nil && s.TaskAward.Enable {
		tasks = append(tasks, s.TaskAward)
	}
	currentPage := false
	if s.TaskActivity != nil && s.TaskActivity.Enable {
		if s.TaskActivity.ActivityID == "" {
			return nil, fmt.Errorf("请选择活动")
		}
		tasks = append(tasks, s.TaskActivity)
		currentPage = s.TaskActivity.CurrentPage
	}
	if s.TaskReclamation != nil && s.TaskReclamation.Enable {
		if err := s.TaskReclamation.Validate(); err != nil {
			return nil, err
		}
		tasks = append(tasks, s.TaskReclamation)
		currentPage = currentPage || s.TaskReclamation.CurrentPage
	}
	if currentPage {
		if scheduled || len(tasks) != 1 {
			return nil, fmt.Errorf("从当前游戏页面开始需要单独手动运行")
		}
		return tasks, nil
	}
	batch := []Tasker{&TaskStartUp{Enable: true, StartGameEnabled: true}}
	batch = append(batch, tasks...)
	return append(batch, &TaskCloseDown{Enable: true}), nil
}
