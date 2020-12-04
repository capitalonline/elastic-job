package models

import "encoding/json"

type (
	// 流水线调度记录模型
	PipelineRecords struct {
		Id           string         `json:"id" xorm:"not null pk comment('ID') CHAR(36)"`
		PipelineId   string         `json:"p_id" xorm:"not null comment('流水线ID') index CHAR(36)"`
		NodeId       string         `json:"n_id" xorm:"not null comment('节点ID') index CHAR(36)"`
		WorkerName   string         `json:"w_name" xorm:"not null comment('节点名称') VARCHAR(255)"`
		ScheduleType string         `json:"sched_type" xorm:"comment('调度类型') CHAR(18)"`
		Spec         string         `json:"spec,omitempty" xorm:"comment('定时器') CHAR(64)"`
		ExecTime     int64          `json:"exec_time,omitempty" xorm:"not null comment('持续时间') INT(10)"`
		Status       int            `json:"status" xorm:"not null comment('状态') INT(10)"`
		Duration     int64          `json:"drt" xorm:"not null comment('持续时间') INT(10)"`
		BeginWith    int64          `json:"b_with" xorm:"not null comment('任务开始时间') INT(10)"`
		FinishWith   int64          `json:"f_with" xorm:"not null comment('任务完成时间') INT(10)"`
		CreatedAt    int64          `json:"c_at" xorm:"not null comment('任务创建时间') INT(10)"`
		UpdatedAt    int64          `json:"u_at" xorm:"not null comment('任务更新时间') INT(10)"`
		Steps        []*TaskRecords `json:"-" xorm:"-"`
	}
	// 流水线执行结果
	Result struct {
		Pipeline *PipelineRecords // 流水线执行记录
		Pr       *PipelineRuntimeData `json:"-"` // 流水线执行记录
	}

	TaskRecords struct {
		Id               int64  `json:"id" xorm:"pk autoincr comment('ID') BIGINT(20)"`
		PipelineRecordId string `json:"p_r_id" xorm:"not null comment('流水线记录ID') index CHAR(36)"`
		TaskId           string `json:"t_id" xorm:"not null comment('任务ID') index VARCHAR(24)"`
		NodeId           string `json:"n_id" xorm:"not null comment('节点ID') index index CHAR(36)"`
		TaskName         string `json:"t_n" xorm:"not null comment('任务名称') VARCHAR(255)"`
		WorkerName       string `json:"w_n" xorm:"not null comment('节点名称') VARCHAR(255)"`
		Mode             string `json:"mode" xorm:"not null comment('执行方式') VARCHAR(255)"`
		Url              string `json:"url,omitempty" xorm:"comment('http地址') VARCHAR(255)"`
		Method           string `json:"method,omitempty" xorm:"comment('http方法') CHAR(36)"`
		Script           string `json:"script,omitempty" xorm:"comment('ansible脚本') CHAR(36)"`
		Hosts            string `json:"hosts,omitempty" xorm:"comment('ansible目标主机') TEXT"`
		Content          string `json:"content" xorm:"not null comment('执行内容') TEXT"`
		Timeout          int    `json:"timeout" xorm:"not null default 0 comment('超时时间') INT(10)"`
		Retries          int    `json:"r" xorm:"not null default 0 comment('重试次数') INT(10)"`
		Status           string `json:"status" xorm:"not null default 'finished' comment('状态') VARCHAR(255)"`
		Result           string `json:"res" xorm:"not null comment('执行结果') TEXT"`
		Duration         int64  `json:"drt" xorm:"not null comment('持续时间') INT(10)"`
		BeginWith        int64  `json:"b_with" xorm:"not null comment('开始于') INT(10)"`
		FinishWith       int64  `json:"f_with" xorm:"not null comment('结束于') INT(10)"`
		CreatedAt        int64  `json:"c_at" xorm:"not null comment('创建于') INT(10)"`
	}
)

func (records *PipelineRecords) TableName() string {
	return "pipeline_records"
}

func (records *TaskRecords) TableName() string {
	return "task_records"
}

// 保存流水线记录
func (records *PipelineRecords) Store() error {
	_, err := Engine.InsertOne(records)
	return err
}

// 保存流水线记录
func (records *TaskRecords) Store() error {
	_, err := Engine.InsertOne(records)
	return err
}


func (records *PipelineRecords) ToString() string {
	result, _ := json.Marshal(records)
	return string(result)
}

func (records *TaskRecords) ToString() string {
	result, _ := json.Marshal(records)
	return string(result)
}