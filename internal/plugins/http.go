package plugins

import (
	"context"
	"github.com/mongodb-job/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type (
	Http struct {
		Url     string
		Method  string
		Content string
	}
)

func (actuator *Http) Exec(ctx context.Context, step int, pr models.PipelineRuntimeData) *models.TaskRecords {
	actuator.Content = Render(actuator.Content, pr)
	client := &http.Client{}
	record := &models.TaskRecords{}
	acc := models.ParamResult{Param: models.StringMap{"content": actuator.Content}}
	req, err := http.NewRequest(actuator.Method, actuator.Url, strings.NewReader(actuator.Content))
	if err != nil {
		record.Result = err.Error()
		record.Status = "failed"
		return record
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		record.Result = err.Error()
		record.Status = "failed"
		return record
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		record.Result = err.Error()
		record.Status = "failed"
		return record
	}

	record.Result = string(body)
	record.Status = "finished"
	acc.Result = []models.StringMap{{"result": record.Result}}
	pr.StepPR[step] = acc
	log.Printf("Run HTTP Record status: %s", record.Status)
	return record
}

func RunHttp(){
	counter := 0
	for {
		if counter > 120 {
			return
		}

		switch "resp.Data.Status" {
		case "FINISH":
			return
		case "DOING":
			break
		case "ERROR":
			return
		}
		time.Sleep(30 * time.Second)
		counter = counter + 1
	}
}