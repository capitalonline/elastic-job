package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/mongodb-job/internal/exector"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	PUT  = 1
	DEL  = 2
	KILL = 3
)

type (
	Event struct {
		Type     int              // 事件类型
		Pipeline *models.Pipeline // 流水线
	}
	Contract interface {
		Run(ctx context.Context)             // 运行调度器
		DispatchEvent(event *Event)          // 分发事件
		eventHandler(event *Event)           // 事件处理
		ResultHandler(result *models.Result) // 调度结果处理
	}
)

type Scheduler struct {
	EventsChan chan *Event                 // 事件通道
	ResultChan chan *models.Result         // 执行结果通道
	Plan       map[string]*models.Pipeline // 调度计划
	Running    sync.Map // 正在运行的流水线
}

var Instance *Scheduler

func (sc *Scheduler) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic:%v", r)
		}
	}()
	logrus.Info("Scheduler is running")
	scheduleTimer := time.NewTimer(5 * time.Second)
	for {
		select {
		case event := <-sc.EventsChan:
			sc.eventHandler(event)
		case <-scheduleTimer.C:
		case result := <-sc.ResultChan:
			sc.ResultHandler(result)
		}
		sc.TrySchedule(ctx)
		//TODO 想办法减少空转
		scheduleTimer.Reset(5 * time.Second)
	}
}

func (sc *Scheduler) TrySchedule(ctx context.Context) (after time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic:%v", r)
		}
	}()

	if len(sc.Plan) == 0 {
		after = 1 * time.Second
		return
	}

	now := time.Now()

	for _, pipeline := range sc.Plan {
		if pipeline.ScheduleType == "delay"{
			if pipeline.Status != 3 {
				if time.Unix(pipeline.ExecTime, 0).Before(now) {
					logrus.Infof("Begin run pipeline, id: %s, name: %s, clusterId: %s, exec: %d", pipeline.Id, pipeline.Name, pipeline.ClusterId, pipeline.ExecTime)
					go exector.RunPipeline(ctx, *pipeline, sc.ResultChan)
					pipeline.Status = 3
				}
			} else {
				//任务执行完成后两天删除
				if time.Unix(pipeline.ExecTime, 0).Add(48*time.Hour).Before(now){
					logrus.Warnf("Delete local delay task, pid: %s, name: %s", pipeline.Id, pipeline.Name)
					delete(sc.Plan, pipeline.Id)
				}
			}
		} else if pipeline.ScheduleType == "cron" {
			if pipeline.NextTime.Before(now) || pipeline.NextTime.Equal(now) {
				logrus.Infof("Begin run pipeline, id: %s, name: %s, clusterId: %s, spec: %s", pipeline.Id, pipeline.Name, pipeline.ClusterId, pipeline.Spec)
				runningPipeline, ok := sc.Running.Load(pipeline.Id)
				if !ok {
					go exector.RunPipeline(ctx, *pipeline, sc.ResultChan)
					sc.Running.Store(pipeline.Id, pipeline)
					pipeline.NextTime = pipeline.Expression.Next(now)
				} else {
					logrus.Infof("This pipeline is running, please wait, info: %+v", runningPipeline)
				}
			}
		}
	}
	return
}
func (sc *Scheduler) DispatchEvent(event *Event) {
	sc.EventsChan <- event
}

func (sc *Scheduler) eventHandler(event *Event) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic:%v", r)
		}
	}()
	switch event.Type {
	case PUT:
		logrus.Println("处理新增或者修改任务:", event.Pipeline.Id)
		if event.Pipeline.ScheduleType == "cron" {
			event.Pipeline.Expression = cronexpr.MustParse(event.Pipeline.Spec)
			event.Pipeline.NextTime = event.Pipeline.Expression.Next(time.Now())
		}
		sc.Plan[event.Pipeline.Id] = event.Pipeline
		a, ok := sc.Plan[event.Pipeline.Id]
		if !ok {
			logrus.Error("no pipeline")
		}
		logrus.Println(a)
	case DEL:
		delete(sc.Plan, event.Pipeline.Id)
	case KILL:
		//TODO Kill
	}
}

func (sc *Scheduler) ResultHandler(result *models.Result) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic error: %v", r)
		}
	}()

	sc.Running.Delete(result.Pipeline.PipelineId)
	logrus.Infof("Release running lock for pipeline: %s", result.Pipeline.PipelineId)

	logrus.Infof("Deal with result, pipelineId: %s, startTime: %d", result.Pipeline.Id, result.Pipeline.ExecTime)
	if err := result.Pipeline.Store(); err != nil {
		logrus.WithFields(logrus.Fields{
			"mode": service.Runtime.Mode,
			"id": service.Runtime.Id,
			"desc": service.Runtime.Description,
			"host": service.Runtime.Host,
			"msg": result.Pipeline.ToString(),
		}).Error("Save result to Mysql failed")
	}
	for _, step := range result.Pipeline.Steps {
		if err := step.Store(); err != nil {
			logrus.WithFields(logrus.Fields{
				"mode": service.Runtime.Mode,
				"id": service.Runtime.Id,
				"desc": service.Runtime.Description,
				"host": service.Runtime.Host,
				"msg": step.ToString(),
			}).Error("Save result to Mysql failed")
		}
	}
	if result.Pipeline.Status == 1 {
		p := sc.Plan[result.Pipeline.PipelineId]
		if strings.HasPrefix(p.Name, "FROMMONGOSERVICEFULL") {
			logrus.Infof("增加全量自动备份结果, 集群: %s", p.ClusterId)
			// 获取副本集size
			if len(result.Pipeline.Steps) < 3 {
				logrus.Infof("结果出现异常, 请在OP上查看结果")
			}
			backupMap := make([]map[string]string, 0)
			backupOrigin := result.Pipeline.Steps[1].Result
			err2 := json.Unmarshal([]byte(backupOrigin), &backupMap)
			if err2 != nil {
				logrus.Errorf("Parse task result to map failed, error: %v", err2)
				return
			}

			// 获取副本集
			uploadMap := make([]map[string]string, 0)
			uploadOrigin := result.Pipeline.Steps[2].Result
			err := json.Unmarshal([]byte(uploadOrigin), &uploadMap)
			if err != nil {
				logrus.Errorf("Parse task result to map failed, error: %v", err)
				return
			}

			converBackuptMap := make(map[string]map[string]string)
			for _, v := range backupMap {
				if rsn, ok := v["replica_set_name"]; !ok {
					logrus.Errorf("结果异常，未找到副本集参数")
					continue
				} else {
					converBackuptMap[rsn] = v
				}
			}

			converUploadMap := make(map[string]map[string]string)
			for _, v := range uploadMap {
				converUploadMap[v["replica_set_name"]] = v
			}
			finishTime := time.Unix(result.Pipeline.FinishWith, 10)
			startTime := time.Unix(result.Pipeline.BeginWith, 10)

			br := models.BackupRecord{
				BackupId:     uuid.NewV4().String(),
				BackupName:   "BACKUPAUTO" + strconv.FormatInt(result.Pipeline.BeginWith, 10),
				BackupType:   "Physics",
				BackupAction: "Auto",
				Deadline:     8,
				Status:       "OK",
				FinishTime:   &finishTime,
				UploadStatus: "OK",
				StartTime:    startTime,
				Valid:        1,
				ClusterId:    p.ClusterId,
				UserId:       p.UserId,
			}
			if err = br.Store(); err != nil {
				logrus.Errorf("Save full backup failed, data: %+v, error: %+v", br, err)
			}

			for K, V := range converUploadMap {
				size, _ := strconv.Atoi(converBackuptMap[K]["size"])
				s3 := &models.S3File{
					FileId:         uuid.NewV4().String(),
					BackupId:       br.BackupId,
					Path:           V["remote_path"],
					Size:           size,
					Md5:            V["md5"],
					CreateTime:     time.Now().Unix(),
					Valid:          1,
					ReplicasetName: K,
				}
				if err = s3.Store(); err != nil{
					logrus.Errorf("Save full backup failed, data: %+v, error: %+v", s3, err)
				}
			}

		}
		if strings.HasPrefix(p.Name, "FROMMONGOSERVICEINCR") {
			logrus.Infof("增加增量自动备份结果, 集群: %s", p.ClusterId)
			backupMap := make([]map[string]string, 0)
			backupOrigin := result.Pipeline.Steps[0].Result
			err2 := json.Unmarshal([]byte(backupOrigin), &backupMap)
			if err2 != nil {
				logrus.Errorf("Parse task result to map failed, error: %v", err2)
				return
			}

			uploadMap := make([]map[string]string, 0)
			uploadOrigin := result.Pipeline.Steps[1].Result
			err := json.Unmarshal([]byte(uploadOrigin), &uploadMap)
			if err != nil {
				logrus.Errorf("Parse task result to map failed, error: %v", err)
				return
			}

			converBackuptMap := make(map[string]map[string]string)
			for _, v := range backupMap {
				converBackuptMap[v["replica_set_name"]] = v
			}

			converUploadMap := make(map[string]map[string]string)
			for _, v := range uploadMap {
				converUploadMap[v["replica_set_name"]] = v
			}


			for K, V := range converUploadMap {
				size, _ := strconv.Atoi(converBackuptMap[K]["size"])
				oplogStartTime, _ := strconv.ParseInt(converBackuptMap[K]["start_time"], 10, 64)
				oplogEndTime, _ := strconv.ParseInt(converBackuptMap[K]["end_time"], 10, 64)
				s3 := &models.IncrementalBackupRecord{
					BackupId:       uuid.NewV4().String(),
					BackupName:     "BACKUPAUTOOPLOG-" + K + "-" + converBackuptMap[K]["start_time"],
					Path:           V["remote_path"],
					OplogStartTime: oplogStartTime,
					OplogEndTime:   oplogEndTime,
					Status:         1,
					CreateTime:     time.Now().Unix(),
					Valid:          1,
					Deadline:       9,
					BackupAction:   "Auto",
					BackupType:     "Oplog",
					ClusterId:      p.ClusterId,
					UserId:         p.UserId,
					ReplicaSetName: K,
					Size:           size,
					Md5:            V["md5"],
					StartTime:      result.Pipeline.BeginWith,
					EndTime:        result.Pipeline.FinishWith,
				}
				if err = s3.Store(); err != nil{
					logrus.Errorf("Save full backup failed, data: %+v, error: %+v", s3, err)
				}
			}
		}
	} else if result.Pipeline.Status == 0 {
		p := sc.Plan[result.Pipeline.PipelineId]
		if strings.HasPrefix(p.Name, "FROMMONGOSERVICEFULL") {
			// 更新集群状态
			models.UpdateClusterStatusToRunning(p.ClusterId)
			// 增加失败记录
			finishTime := time.Unix(result.Pipeline.FinishWith, 10)
			startTime := time.Unix(result.Pipeline.BeginWith, 10)
			br := models.BackupRecord{
				BackupId:     uuid.NewV4().String(),
				BackupName:   "BACKUPAUTO" + strconv.FormatInt(result.Pipeline.BeginWith, 10),
				BackupType:   "Physics",
				BackupAction: "Auto",
				Deadline:     8,
				Status:       "Failed",
				FinishTime:   &finishTime,
				UploadStatus: "Failed",
				StartTime:    startTime,
				Valid:        1,
				ClusterId:    p.ClusterId,
				UserId:       p.UserId,
			}
			if err := br.Store(); err != nil {
				logrus.Errorf("Save full backup failed, data: %+v, error: %+v", br, err)
			}

			//通知GIC
			reason := result.Pipeline.Steps[len(result.Pipeline.Steps)-1].Result
			go notifyGIC("全量备份失败:"+reason, *p)
		}
		if strings.HasPrefix(p.Name, "FROMMONGOSERVICEINCR"){
			//通知GIC
			reason := result.Pipeline.Steps[len(result.Pipeline.Steps)-1].Result
			go notifyGIC("增量备份失败:"+reason, *p)
		}
	}
}

func InitScheduler() {
	Instance = &Scheduler{
		EventsChan: make(chan *Event, 1000),
		ResultChan: make(chan *models.Result, 1000),
		Plan:       make(map[string]*models.Pipeline),
		//Running:    sync.Map{},
	}
}

func notifyGIC(presult string, pipeline models.Pipeline) {
	c, err := models.FindClusterByID(pipeline.ClusterId)
	if err != nil || c.ClusterId != pipeline.ClusterId {
		logrus.Errorf("处理告警失败，集群信息未查到，集群ID(%s), 错误信息(%v), 备份失败信息(%s)", pipeline.ClusterId, err.Error(), "请查询数据库")
		return
	}

	customer := models.GetCustomer(c.CustomerId)

	infoTmp := "任务ID：%s\r\n任务类型：%s\r\n账户编号：%s\r\n账户名称：%s\r\n实例编号：%s\r\n实例名称：%s\r\n错误信息：%s"
	msg := fmt.Sprintf(infoTmp, pipeline.Id, "Backup", c.CustomerId, customer.CustomerName, c.ClusterId, c.ClusterName, presult)
	report := models.FailureReport{
		Hostname:     "云数据库MongoDB备份失败",
		SubObject:    fmt.Sprintf("MongoDB(%s)", c.ClusterType),
		Ip:           c.ConnectionUri,
		Level:        "Alert",
		LogTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		Customername: customer.CustomerName,
		Tag1:         "Backup Failure",
		Message:      msg,
	}
	report.SendFailureToGICNewMonitor()
}

