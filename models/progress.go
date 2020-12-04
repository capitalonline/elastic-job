package models

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

const (
	DISPATCH   = "dispatch"
	UPDATETIME = "update"
	Knowleage  = "know"
	BEGIN      = "begin"
	RUNNING    = "running"
	WAITNEXT   = "waitnext"
	FAILED     = "failed"
)

type (
	Progress struct {
		SchdulePipeline *Pipeline         `json:"schdule_pipeline"`
		Status          string            `json:"status"`
		MeteData        map[string]string `json:"mete_data"`
	}
)

func (p *Progress) ToString () string {
	bytes, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("Parse progress to json failed, error: %s", err.Error())
		return ""
	}
	return string(bytes)
}