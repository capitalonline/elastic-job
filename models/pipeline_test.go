package models

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestNewPipeline(t *testing.T) {
	id := uuid.NewV4().String()
	p := Pipeline{
		Id:           id,
		Name:         "测试集群ansible定时备份",
		Description:  "测试集群ansible定时备份",
		ScheduleType: "cron",
		Spec:         "*/15 * * * *",
		Status:       1,
		Steps: []*PipelineTaskPivot{
			{
				Id:          "0",
				PipelineId:  id,
				Step:        0,
				Task: &Task{
					Id:          "0",
					Name:        "backup",
					Mode:        MODEANSIBLE,
					Script:      "physics_backup.yaml",
					Hosts:       "202.202.0.18 ansible_ssh_pass=V!1DFkHUb vm_role=hidden private_ip=10.241.20.2 replica_set_name=Rpl-HJsQUr",
					Content:     `{"backup_path":"/data/data-backup/CustomerID/strings/physicsbackup","type":"PhysicsBackup",backup_host_ip":"0.0.0.0","tar_filename":"{{nowtime}}","host_ip":"100.131.0.45 mongo-pre.ae327e0452c545349e4b6b41478d72b5.oss-cnbj01.cdsgss.com","super_user_info":{"username":"cds_root","passwd":"S18NAMXPkDis"}}`,
					Description: "ansible backup",
				},
			},
			{
				Id:          "1",
				PipelineId:  id,
				Step:        1,
				Task: &Task{
					Id:          "1",
					Name:        "va param",
					Mode:        MODESHELL,
					Content:     `echo {{s3_remote . 0 "Rpl-HJsQUr"}}`,
					Description: "ansible backup",
				},
			},
		},
		FinishedTask: &Task{
			Id:          "0",
			Name:        "echo ok",
			Mode:        MODESHELL,
			Content:     "echo 'hello ok'",
			Description: "echo",
		},
		FailedTask: &Task{
			Id:          "0",
			Name:        "echo faild",
			Mode:        MODESHELL,
			Content:     "echo 'hello faild'",
			Description: "echo",
		},
	}

	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(p)
	fmt.Println(string(marshal))
}

//Rpl-HJsQUr
//V!1DFkHUb
//202.202.0.18
//10.241.20.2



//info.VMIP, info.VMPwd, map[string]string{"vm_role": info.VMRole, "privte_ip": info.PrivateIP, "replica_set_name": se.Shard[0].ReplSetName})/
