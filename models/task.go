package models

import "encoding/json"

const (
	MODESHELL   = "shell"
	MODEHTTP    = "http"
	MODEANSIBLE = "ansible"
	MODEHOOK    = "hook"
)

type Task struct {
	Id          string `json:"id" validate:"-"`
	Name        string `json:"name" validate:"required"`
	Mode        string `json:"mode" validate:"required"`
	Url         string `json:"url,omitempty" validate:"omitempty"`
	Method      string `json:"method,omitempty" validate:"omitempty"`
	Script      string `json:"script,omitempty"`
	Hosts       string `json:"hosts,omitempty"`
	Content     string `json:"content,omitempty" validate:"omitempty"`
	Description string `json:"desc" validate:"-"`
}

type (
	PipelineRuntimeData struct {
		PipelineId string
		StepPR     StepParamResult
	}
	StepParamResult map[int]ParamResult
	ParamResult     struct {
		Param  StringMap   `json:"param"`
		Result []StringMap `json:"result"`
	}
	StringMap map[string]string
)

func(p *PipelineRuntimeData) ToString() string {
	marshal, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(marshal)
}

//type LogicalBackupAction struct {
//	BackupPath    string    `json:"backup_path"`
//	BackupHostIP  string    `json:"backup_host_ip"`
//	Type          string    `json:"type"`
//	TarFilename   string    `json:"tar_filename"`
//	DBInfoPath    string    `json:"db_info_path"`
//	HostIp        string    `json:"host_ip"`
//	SuperUserInfo SuperUser `json:"super_user_info"`
//}
