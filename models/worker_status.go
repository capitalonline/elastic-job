package models

import (
	"github.com/sirupsen/logrus"
	"time"
)

const (
	_updateWorkerPipelineStatus = "UPDATE worker_pipeline_status SET `status`=?,`step`=?,`task_name`=?,`maybe_data`=?,`updated_at`=? WHERE `pipeline_id`=? AND `node_id`=?"
	_updateWorkerPipelineOnlyStatus = "UPDATE worker_pipeline_status SET `status`=?,`updated_at`=? WHERE `pipeline_id`=? AND `node_id`=?"
)

type (
	WorkerPipelineCurrentStatus struct {
		Id         string `json:"id" xorm:"not null pk comment('ID') CHAR(36)"`
		PipelineId string `json:"p_id" xorm:"not null comment('流水线ID') index CHAR(36)"`
		NodeId     string `json:"n_id" xorm:"not null comment('节点ID') index CHAR(36)"`
		WorkerName string `json:"w_name" xorm:"not null comment('节点名称') VARCHAR(255)"`
		Status     string `json:"status" xorm:"not null comment('状态:已知，开始执行，运行中，等待下一次') VARCHAR(64)"`
		Step       int    `json:"step" xorm:"comment('流水线当前执行任务') INT(10)"`
		TaskName   string `json:"task_name" xorm:"comment('流水线执行任务名称') VARCHAR(255)"`
		MaybeData  string `json:"maybe_data" xorm:"comment('流水线执行任务的参数') VARCHAR(255)"`
		CreatedAt  int64  `json:"c_at" xorm:"not null comment('任务创建时间') INT(10)"`
		UpdatedAt  int64  `json:"u_at" xorm:"not null comment('任务更新时间') INT(10)"`
	}

	PipelineStub struct {
		Id         string `json:"id" xorm:"not null pk comment('ID') CHAR(36)"`
		PipelineId string `json:"p_id" xorm:"not null comment('流水线ID') index CHAR(36)"`
		NodeId     string `json:"n_id" xorm:"not null comment('节点ID') index CHAR(36)"`
		WorkerName string `json:"w_name" xorm:"not null comment('节点名称') VARCHAR(255)"`
		Stub       string `json:"stub,omitempty" xorm:"comment('存根数据') VARCHAR(255)"`
		CreatedAt  int64  `json:"c_at" xorm:"not null comment('存根创建时间') INT(10)"`
		UpdatedAt  int64  `json:"u_at" xorm:"not null comment('存根更新时间') INT(10)"`
	}
)

func (records *WorkerPipelineCurrentStatus) TableName() string {
	return "worker_pipeline_status"
}

// 保存流水线记录
func (records *WorkerPipelineCurrentStatus) Store() error {
	_, err := Engine.InsertOne(records)
	return err
}

func (records *PipelineStub) TableName() string {
	return "pipeline_stub"
}

// 保存流水线记录
func (records *PipelineStub) Store() error {
	_, err := Engine.InsertOne(records)
	return err
}

func ExecUpdateWorkerPipelineStatus(status string, step int, taskName, maybeData, pipelineId, nodeId string) {
	nowTime := time.Now().Unix()
	if step == 0 && taskName == "" && maybeData == ""{
		if _, err := Engine.Exec(_updateWorkerPipelineOnlyStatus, status, nowTime, pipelineId, nodeId); err != nil {
			logrus.WithFields(logrus.Fields{
				"SQL":         _updateWorkerPipelineOnlyStatus,
				"status":      status,
				"pipeline_id": pipelineId,
				"node_id":     nodeId,
			}).Error("db.worker_pipeline_status.update() Failed")
		}
	} else {
		if _, err := Engine.Exec(_updateWorkerPipelineStatus, status, step, taskName, maybeData, nowTime, pipelineId, nodeId); err != nil {
			logrus.WithFields(logrus.Fields{
				"SQL":         _updateWorkerPipelineStatus,
				"status":      status,
				"step":        step,
				"task_name":   taskName,
				"maybe_data":  maybeData,
				"pipeline_id": pipelineId,
				"node_id":     nodeId,
			}).Error("db.worker_pipeline_status.update() Failed")
		}
	}
}
