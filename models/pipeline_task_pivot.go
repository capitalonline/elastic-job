package models

type PipelineTaskPivot struct {
	Id          string `json:"id"`
	PipelineId  string `json:"pipeline_id" validate:"required,uuid4"`
	Step        int    `json:"step" validate:"numeric"`
	Timeout     int    `json:"timeout,omitempty" validate:"numeric"`
	Interval    int    `json:"interval,omitempty" validate:"numeric"`
	Retries     int    `json:"retries,omitempty" validate:"numeric"`
	Directory   string `json:"dir,omitempty" validate:"omitempty"`
	User        string `json:"user,omitempty" validate:"omitempty"`
	Environment string `json:"env,omitempty" validate:"omitempty"`
	Dependence  string `json:"dependence,omitempty" validate:"required"`
	CreatedAt   int64  `json:"c_at,omitempty" validate:"-"`
	UpdatedAt   int64  `json:"u_at,omitempty" validate:"-"`
	Task        *Task  `json:"task" validate:"-"`
}
