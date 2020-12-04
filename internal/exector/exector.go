package exector

import (
	"context"
	"github.com/mongodb-job/internal/plugins"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
)

func RunPipeline(ctx context.Context, pipeline models.Pipeline, resChan chan *models.Result) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic in run pipeline with name: %s, id: %s, error: %v", pipeline.Name, pipeline.Id, r)
			return
		}
	}()
	if len(pipeline.Steps) > 0 {
		record := &models.PipelineRecords{
			Id:           uuid.NewV4().String(),
			PipelineId:   pipeline.Id,
			NodeId:       service.Runtime.Id,
			WorkerName:   service.Runtime.Name,
			ScheduleType: pipeline.ScheduleType,
			Spec:         pipeline.Spec,
			ExecTime:     pipeline.ExecTime,
			Status:       1,
		}

		beginWith := time.Now()
		record.BeginWith = beginWith.Unix()
		result := &models.Result{}
		runtimeData := models.PipelineRuntimeData{
			PipelineId: pipeline.Id,
			StepPR:     make(models.StepParamResult),
		}
		models.ExecUpdateWorkerPipelineStatus(models.BEGIN, 0, "", "", pipeline.Id, service.Runtime.Id)
		for _, step := range pipeline.Steps {
			// 如果设置了超时时间
			var pctx context.Context
			var cancelFunc func()
			if step.Timeout == 0 {
				pctx, cancelFunc = context.WithCancel(ctx)
			} else {
				pctx, cancelFunc = context.WithTimeout(ctx, time.Duration(step.Timeout)*time.Second)
			}
			logrus.Infof("Run pipeline: %s step %s", pipeline.Id, step.Id)
			taskRecord := RunStep(pctx, step, runtimeData)
			taskRecord.Timeout = step.Timeout
			taskRecord.Retries = step.Retries
			taskRecord.PipelineRecordId = record.Id
			taskRecord.CreatedAt = time.Now().Unix()

			record.Steps = append(record.Steps, taskRecord)

			if taskRecord.Status == "failed" {
				logrus.Errorf("Failed run pipeline: %s step %s error: %s", pipeline.Id, step.Id, taskRecord.Result)
				record.Status = 0
				goto END
			}

			models.ExecUpdateWorkerPipelineStatus(models.RUNNING, step.Step, step.Task.Name, runtimeData.ToString(), pipeline.Id, service.Runtime.Id)

			select {
			//TODO 增加强杀任务功能
			case <-ctx.Done():
				cancelFunc()
				break
			default:
				continue
			}

		}

	END:
		finishWith := time.Now()
		record.Duration = int64(finishWith.Sub(beginWith).Seconds())

		record.FinishWith = finishWith.Unix()
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()
		log.Printf("RunPipeline end : %d", record.Status)
		if record.Status == 1 {
			models.ExecUpdateWorkerPipelineStatus(models.WAITNEXT, 0, "", "", pipeline.Id, service.Runtime.Id)
			if pipeline.FinishedTask != nil {
				switch pipeline.FinishedTask.Mode {
				case models.MODESHELL:
					break
				case models.MODEHTTP:
					break
				case models.MODEHOOK:
					break
				}
			}
		}

		if record.Status == 0 {
			models.ExecUpdateWorkerPipelineStatus(models.FAILED, 0, "", "", pipeline.Id, service.Runtime.Id)
			if pipeline.FailedTask != nil {
				switch pipeline.FailedTask.Mode {
				case models.MODESHELL:
					break
				case models.MODEHTTP:
					break
				case models.MODEHOOK:
					break
				}
			}
		}

		result.Pipeline = record
		result.Pr = &runtimeData
		resChan <- result
		return
	}
}

func RunStep(ctx context.Context, step *models.PipelineTaskPivot, pr models.PipelineRuntimeData) *models.TaskRecords {
	record := &models.TaskRecords{}
	beginWith := time.Now()
	if step.Retries == 0 {
		record = run(ctx, step, pr)
	} else {
		for i := 0; i < step.Retries; i++ {
			record = run(ctx, step, pr)
			if record.Status == "finished" {
				break
			}
			time.Sleep(time.Duration(step.Interval) * time.Second)
		}
	}
	record.TaskId = step.Task.Id
	record.NodeId = service.Runtime.Id
	record.TaskName = step.Task.Name
	record.WorkerName = service.Runtime.Name
	if s, ok := pr.StepPR[step.Step].Param["content"]; ok {
		record.Content = s
	} else {
		record.Content = step.Task.Content
	}
	record.Mode = step.Task.Mode
	record.Hosts = step.Task.Hosts
	record.Url = step.Task.Url
	record.Method = step.Task.Method
	record.Script = step.Task.Script
	record.Timeout = step.Timeout
	record.Retries = step.Retries
	finishWith := time.Now()
	record.BeginWith = beginWith.Unix()
	record.FinishWith = time.Now().Unix()
	record.Duration = int64(finishWith.Sub(beginWith).Seconds())
	return record
}

func run(ctx context.Context, step *models.PipelineTaskPivot, pr models.PipelineRuntimeData) *models.TaskRecords {
	switch step.Task.Mode {
	case models.MODEANSIBLE:
		ansible := &plugins.Ansible{
			OrginHost:    step.Task.Hosts,
			OrginContent: step.Task.Content,
			OriginScript: step.Task.Script,
		}
		return ansible.Exec(ctx, step.Step, pr)
	case models.MODEHTTP:
		http := &plugins.Http{
			Url:     step.Task.Url,
			Method:  step.Task.Method,
			Content: step.Task.Content,
		}
		return http.Exec(ctx, step.Step, pr)
	case models.MODESHELL:
		shell := &plugins.Shell{
			User:    step.User,
			Env:     strings.Split(step.Environment, " "),
			Dir:     step.Directory,
			Command: step.Task.Content,
		}
		return shell.Exec(ctx, step.Step, pr)
	}
	return nil
}
