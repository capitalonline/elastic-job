package plugins

import (
	"context"
	"github.com/mongodb-job/models"
	"testing"
	"time"
)

func TestHttp_Exec(t *testing.T) {
	now := time.Now()
	if time.Unix(1697297983, 0).Before(now) {
		println("fff")
	}
	h:= Http{
		Url:     "http://202.202.0.2:8080/inner/v1/updateCluster",
		Method:  "POST",
		Content: "{\"cluster_id\":\"142bfec6-c2f9-495b-8fb8-93bb84ea858c\", \"status\":\"Backuping\"}",
	}
	var result models.PipelineRuntimeData
	h.Exec(context.TODO(), 1, result)
	print("end")
}
