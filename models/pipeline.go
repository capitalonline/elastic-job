package models

import "time"
import "github.com/gorhill/cronexpr"

const (
	CREATE = 1
	UPDATE = 2
)

type Pipeline struct {
	Id           string               `json:"id" validate:"-" `
	ClusterId    string               `json:"clusterid" validate:"required"`
	UserId       string               `json:"user_id" validate:"required"`
	Name         string               `json:"name" validate:"required"`
	Description  string               `json:"desc,omitempty" validate:"-"`
	ScheduleType string               `json:"sched_type" validate:"omitempty"` // cron | delay
	ExecTime     int64                `json:"exec_time,omitempty"`
	Spec         string               `json:"spec,omitempty"`
	Status       int                  `json:"status,omitempty" validate:"numeric"`
	CreatedAt    int64                `json:"c_at" validate:"-"`
	UpdatedAt    int64                `json:"u_at" validate:"-"`
	Nodes        string               `json:"node"`
	Steps        []*PipelineTaskPivot `json:"steps"`
	Expression   *cronexpr.Expression `json:"-"`
	NextTime     time.Time            `json:"-"`
	FinishedTask *Task                `json:"finished_t,omitempty"`
	FailedTask   *Task                `json:"failed_t,omitempty"`
}
type PipelineVo struct {
	PipelineId           string               `json:"pipeline_id" validate:"required" `
	Spec                 string               `json:"spec" validate:"required"`
}
