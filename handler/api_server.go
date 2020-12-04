package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/mongodb-job/internal/dispatcher"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/render"
	"github.com/mongodb-job/internal/scheduler"
	"github.com/mongodb-job/internal/worker_manager"
	"github.com/mongodb-job/mcode"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	validate = validator.New()
)

// @Summary save job
// @Produce json
// @Success 200 {object} render.JSON{code:200}
// @Router /inner/v1/save [post]
func HandleJobSave(c *gin.Context) {
	r := render.New(c)
	var pipeline models.Pipeline
	err := c.Bind(&pipeline)
	if err != nil {
		r.JSON(err.Error(), mcode.RequestErr)
		return
	}
	if err = validate.Struct(pipeline); err != nil {
		logrus.Error(err)
		r.JSON(err.Error(), err)
		return
	}
	now := time.Now().Unix()
	pipeline.CreatedAt = now
	pipeline.UpdatedAt = now

	for _, step := range pipeline.Steps {
		step.UpdatedAt = now
		step.CreatedAt = now
		step.Task.Content = strings.ReplaceAll(step.Task.Content, "\\\"","\"")
		step.Task.Hosts = strings.ReplaceAll(step.Task.Hosts, "\\\"","\"")
	}

	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Pipeline, pipeline.Id)
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		logrus.Error(err)
		return
	}
	if _, err = etcd.Client.KV.Put(c, key, string(bytes)); err != nil {
		logrus.Error(err)
		return
	}
	r.JSON("ok", mcode.OK)
	return
}

// @Summary delete job
// @Produce json
// @Success 200 {object} render.JSON{code:200}
// @Router /inner/v1/delete [delete]
func HandlerJobDelete(c *gin.Context) {
	r := render.New(c)
	pipelineId := c.Query("pipeline_id")
	if pipelineId == "" {
		r.JSON(nil, mcode.RequestErr)
		return
	}
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Pipeline, pipelineId)
	logrus.Info("delete Pipeline, pipeline id: %s", pipelineId)
	if _, err := etcd.Client.KV.Delete(c, key); err != nil {
		r.JSON(err.Error(), err)
		return
	}
	r.JSON(nil, mcode.OK)
}

// @Summary update job
// @Produce json
// @Success 200 {object} render.JSON{code:200}
// @Router /inner/v1/update [post]
func HandleJobUpdate(c *gin.Context) {
	r := render.New(c)
	var pipelineVo models.PipelineVo
	err := c.Bind(&pipelineVo)
	if err != nil {
		r.JSON(err.Error(), mcode.RequestErr)
		return
	}
	if err = validate.Struct(pipelineVo); err != nil {
		logrus.Error(err)
		r.JSON(err.Error(), err)
		return
	}
	// TODO 参数校验

	//
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Pipeline, pipelineVo.PipelineId)
	getResponse, err := etcd.Client.KV.Get(c, key)
	if err != nil {
		logrus.Error(err)
		return
	}
	if len(getResponse.Kvs) == 0 {
		return
	}
	value := getResponse.Kvs[0].Value
	pipeline := new(models.Pipeline)
	josnErr := json.Unmarshal(value, pipeline)
	if josnErr != nil {
		return
	}
	pipeline.Spec = pipelineVo.Spec
	pipeline.Status = models.UPDATE
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		return
	}
	if _, err = etcd.Client.KV.Put(c, key, string(bytes)); err != nil {
		return
	}
	r.JSON("ok", mcode.OK)
	return
}

// @Summary kill job
// @Produce json
// @Success 200 {object} render.JSON{code:200}
// @Router /inner/v1/save [post]
func HandlerJobKill(c *gin.Context) {
	r := render.New(c)
	// parse params
	r.JSON("", mcode.OK)
	return
}

func ShowPipeline(c *gin.Context) {
	r := render.New(c)
	if resp, err := etcd.Client.KV.Get(c, etcd.SelfConfig.Pipeline, clientv3.WithPrefix()); err != nil {
		r.JSON(nil, mcode.RequestErr)
		return
	} else {
		result := make([]models.Pipeline, 0)
		for _, kv := range resp.Kvs {
			var p models.Pipeline
			_ = json.Unmarshal(kv.Value, &p)
			result = append(result, p)
		}
		r.JSON(result, mcode.OK)
		return
	}
}

func ShowPipelineHistory(c *gin.Context) {
	r := render.New(c)
	pipelineId := c.Query("pipeline_id")
	if pipelineId == "" {
		r.JSON(nil, mcode.RequestErr)
		return
	}
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Snapshot, pipelineId)
	if resp, err := etcd.Client.KV.Get(c, key, clientv3.WithPrefix(), clientv3.WithLimit(30)); err != nil {
		r.JSON(nil, mcode.RequestErr)
		return
	} else {
		result := make([]models.PipelineRecords, 0)
		for _, kv := range resp.Kvs {
			var p models.Result
			_ = json.Unmarshal(kv.Value, &p)
			result = append(result, *p.Pipeline)
		}
		r.JSON(result, mcode.OK)
		return
	}
}

func ShowWorkerPipeline(c *gin.Context) {
	r := render.New(c)
	workerId := c.Query("worker_id")
	if workerId == "" {
		r.JSON(nil, mcode.RequestErr)
		return
	}
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Plan, workerId)
	if resp, err := etcd.Client.KV.Get(c, key, clientv3.WithPrefix()); err != nil {
		r.JSON(nil, mcode.RequestErr)
		return
	} else {
		result := make([]models.Progress, 0)
		for _, kv := range resp.Kvs {
			var p models.Progress
			_ = json.Unmarshal(kv.Value, &p)
			result = append(result, p)
		}
		r.JSON(result, mcode.OK)
		return
	}
}

func JobLogSearch(c *gin.Context) {
	r := render.New(c)
	var (
		total int64
		err   error
	)

	scene := c.Query("scene")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if size == 0 {
		size = 20
	}
	switch scene {
	case "pipeline":
		logs := make([]models.PipelineRecords, 0)
		if search != "" {
			total, err = models.Engine.Where(builder.Like{"pipeline_id", search}).Limit(size, page).Desc("created_at").FindAndCount(&logs)
		} else {
			total, err = models.Engine.Limit(size, page).Desc("created_at").FindAndCount(&logs)
		}
		result := make(map[string]interface{})
		result["total"] = total
		result["page"] = page
		result["size"] = size
		result["data"] = logs
		r.JSON(result, err)
		return
	case "task":
		logs := make([]models.TaskRecords, 0)
		if search != "" {
			field := c.Query("field")
			if field == "" {
				field = "pipeline_record_id"
			}
			total, err = models.Engine.Where(builder.Like{field, search}).Limit(size, page).Desc("created_at").FindAndCount(&logs)
		} else {
			total, err = models.Engine.Limit(size, page).Desc("created_at").FindAndCount(&logs)
		}
		result := make(map[string]interface{})
		result["total"] = total
		result["page"] = page
		result["size"] = size
		result["data"] = logs
		r.JSON(result, err)
		return
	}

}

func WorkerStatus(c *gin.Context) {
	r := render.New(c)
	result := make([]map[string]interface{}, 0)
	workers := worker_manager.WorkerManagerInstance.All()
	for _, worker := range workers {
		workerStatus := make([]models.WorkerPipelineCurrentStatus, 0)
		err := models.Engine.Where(builder.Eq{"node_id": worker.Id}).Find(&workerStatus)
		if err != nil {
			logrus.Error("Exist search worker status failed")
			continue
		}
		resultUnit := make(map[string]interface{})
		resultUnit["worker"] = worker
		resultUnit["pipelines"] = workerStatus

		result = append(result, resultUnit)
	}

	r.JSON(result, nil)
	return
}

func JobMigrate(c *gin.Context) {
	r := render.New(c)
	pid := c.Query("pipeline_id")
	wid := c.Query("worker_id")
	dispatcher.DispatcherInstance.MigrateSinglePipeline(pid, wid)
	r.JSON("", mcode.OK)
	return
}

// FOR WORKER
func WorkersJob(c *gin.Context)  {
	r := render.New(c)
	r.JSON(scheduler.Instance.Plan, nil)
	return
}
