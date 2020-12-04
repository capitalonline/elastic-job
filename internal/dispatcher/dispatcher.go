package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/worker_manager"
	"github.com/mongodb-job/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log"
	"time"

	"sync"
)

type (
	Dispatcher struct {
		sync.Mutex
	}
)

var (
	DispatcherInstance *Dispatcher
)

func init() {
	DispatcherInstance = &Dispatcher{}
}

func (d *Dispatcher) DispatchPipeline(p models.Pipeline) {
	d.Lock()
	defer d.Unlock()
	defer func() {
		if r := recover(); r != nil {
			logrus.Println(r)
		}
	}()
	worker := worker_manager.WorkerManagerInstance.SelectWorker()
	if worker == nil {
		logrus.Errorf("No online worker, please deal it with console, pipeline id is %s", p.Id)
		return
	}
	p.Nodes = worker.Id

	WorkerPipelineInitStatus := models.WorkerPipelineCurrentStatus{
		Id:         uuid.NewV4().String(),
		PipelineId: p.Id,
		NodeId:     worker.Id,
		WorkerName: worker.Name,
		Status:     models.Knowleage,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if saveErr := WorkerPipelineInitStatus.Store(); saveErr != nil {
		logrus.WithFields(logrus.Fields{
			"pipeline_id": p.Id,
			"worker_id": worker.Id,
			"error": saveErr.Error(),
		}).Error("Init worker's pipeline staus failed")
		return
	}

	// 保存任务分配的节点信息
    foreignKey := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign,p.Id)
	if _, foreignErr := etcd.Client.KV.Put(context.Background(), foreignKey, p.Nodes); foreignErr != nil {
		logrus.Errorf("Add foreign failed, pipelineId: %s, workerId: %s", p.Id, worker.Id)
		return
	}

	// 开始分配任务到指定节点
	progress := models.Progress{
		SchdulePipeline: &p,
		Status:          models.DISPATCH,
		MeteData:        make(map[string]string),
	}
	key := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, p.Nodes , p.Id)
	if _, err := etcd.Client.KV.Put(context.Background(), key, progress.ToString()); err != nil {
		logrus.Errorf("Dispatch pipeline failed, error: %v", err.Error())
		etcd.Client.KV.Delete(context.Background(), foreignKey)
		return
	}
	worker_manager.WorkerManagerInstance.CountAdd(worker.Id)
}

func (d *Dispatcher) DeletePipeline(p models.Pipeline){
	d.Lock()
	defer d.Unlock()
	defer func() {
		if r := recover(); r != nil {
			logrus.Println(r)
		}
	}()
	foreignKey := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign,p.Id)
	response, err := etcd.Client.KV.Get(context.Background(), foreignKey)
	if err != nil {
		logrus.Errorf("Get foreign failed, error: %v, key: %s", err, foreignKey)
		return
	}
	if len(response.Kvs) == 0 {
		return
	}
	//node id
	nodeId := string(response.Kvs[0].Value)
	worker := worker_manager.WorkerManagerInstance.Have(nodeId)

	key := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, nodeId, p.Id)
	if _, err := etcd.Client.KV.Delete(context.Background(), key); err != nil {
		logrus.Errorf("Delete plan failed, error: %v, key: %s", err, key)
		return
	}
	if _, err := etcd.Client.KV.Delete(context.Background(), foreignKey); err != nil {
		logrus.Errorf("Delete foreign failed, error: %v, key: %s", err, foreignKey)
		return
	}
	worker_manager.WorkerManagerInstance.CountSub(worker.Id)
}

func (d *Dispatcher) ModifyPipeline(p models.Pipeline){
	d.Lock()
	defer d.Unlock()
	defer func() {
		if r := recover(); r != nil {
			logrus.Println(r)
		}
	}()
	log.Println("Begin modify pipeline")

	// 根据外键找到任务节点
	foreignKey := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign,p.Id)
	workerWithPipeline, err := etcd.Client.KV.Get(context.Background(), foreignKey)
	if err != nil || len(workerWithPipeline.Kvs) == 0{
		logrus.Errorf("Not found task node, pipelineId: %s", p.Id)
		return
	}
	//node id
	nodeId := string(workerWithPipeline.Kvs[0].Value)
	log.Println("任务所在Worker:", nodeId)

	//获取节点上任务信息
	workerKey := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, nodeId, p.Id)
	pipelineProgress, err := etcd.Client.KV.Get(context.Background(), workerKey)
	if err != nil || len(pipelineProgress.Kvs) == 0 {
		logrus.Errorf("Get pipeline foreign failed, error: %+v, pid: %s, nid: %s", err, p.Id, nodeId)
		return
	}

	oldProgress := new(models.Progress)
	if err = json.Unmarshal(pipelineProgress.Kvs[0].Value, oldProgress);err != nil{
		logrus.Errorf("Parse etcd to progress failed, error: %+v", err)
		return
	}

	// 修改并分配任务
	newProgress := new(models.Progress)
	newProgress.SchdulePipeline = &p
	newProgress.Status = models.UPDATETIME
	newProgress.MeteData = oldProgress.MeteData

	newProgressStr := newProgress.ToString()
	oldProgressStr := oldProgress.ToString()

	if newProgressStr == "" || oldProgressStr == "" {
		logrus.Errorf("Unknow failed, newProgress: %+v, oldProgress: %+v", newProgress, oldProgress)
		return
	}

	txnResponse, err := etcd.Client.Txn(context.Background()).If(clientv3.Compare(clientv3.Value(workerKey), "=", oldProgressStr)).
		Then(clientv3.OpPut(workerKey, newProgressStr)).
		Commit()
	if err != nil || !txnResponse.Succeeded{
		logrus.Errorf("Update pipeline failed, error: %+v", err)
		return
	}
}

func (d *Dispatcher) MigratePipeline(workerId string){
	d.Lock()
	defer d.Unlock()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	worker_manager.WorkerManagerInstance.Down(workerId)
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Plan, workerId)
	taskOnWorkers, err := etcd.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, leacyTask := range taskOnWorkers.Kvs {

		var progess models.Progress
		if err := json.Unmarshal(leacyTask.Value, &progess); err != nil {
			log.Println(err)
			continue
		}

		newworker := worker_manager.WorkerManagerInstance.SelectWorker()
		if newworker == nil {
			return
		}
		progess.SchdulePipeline.Nodes = newworker.Id
		progess.Status = models.DISPATCH

		key := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, progess.SchdulePipeline.Nodes , progess.SchdulePipeline.Id)
		bytes, err := json.Marshal(progess)
		if err != nil {
			log.Println(err)
			continue
		}
		if _, err = etcd.Client.KV.Put(context.Background(), key, string(bytes)); err != nil {
			log.Println(err)
			continue
		}
		worker_manager.WorkerManagerInstance.CountAdd(newworker.Id)

		_, err = etcd.Client.Delete(context.Background(), string(leacyTask.Key))
		if err != nil {
			log.Println(err)
			continue
		}
		//todo 更新Foreign nodeId
		foreignKey := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign,progess.SchdulePipeline.Id)
		getResponse, err := etcd.Client.KV.Get(context.Background(), foreignKey)
		if len(getResponse.Kvs) == 0 {
			return
		}
		txnResponse, err := etcd.Client.Txn(context.Background()).If(clientv3.Compare(clientv3.Value(foreignKey), "=", string(getResponse.Kvs[0].Value[:]))).
			Then(clientv3.OpPut(foreignKey, progess.SchdulePipeline.Nodes)).
			Commit()
		if err != nil {
			return
		}
		if !txnResponse.Succeeded {
			return
		}

		worker_manager.WorkerManagerInstance.CountSub(workerId)
	}

	worker_manager.WorkerManagerInstance.RemoveWorker(workerId)
}

func (d *Dispatcher) MigrateSinglePipeline(pipelineId string, workerId string){
	d.Lock()
	defer d.Unlock()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign, pipelineId)
	taskOnWorkers, err := etcd.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		log.Println(err.Error())
		return
	}

	nodeId := taskOnWorkers.Kvs[0].Value
	key2 := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, nodeId, pipelineId)
	progressJson, _ := etcd.Client.Get(context.Background(), key2)
	var progess models.Progress
	if err := json.Unmarshal(progressJson.Kvs[0].Value, &progess); err != nil {
		log.Println(err)
		return
	}
	progess.SchdulePipeline.Nodes = workerId
	progess.Status = models.DISPATCH

	key3 := fmt.Sprintf("%s/%s/%s", etcd.SelfConfig.Plan, progess.SchdulePipeline.Nodes , progess.SchdulePipeline.Id)
	bytes, err := json.Marshal(progess)
	if err != nil {
		log.Println(err)
	}
	if _, err = etcd.Client.KV.Put(context.Background(), key3, string(bytes)); err != nil {
		log.Println(err)
	}
	worker_manager.WorkerManagerInstance.CountAdd(workerId)

	_, err = etcd.Client.Delete(context.Background(), string(key2))

	foreignKey := fmt.Sprintf("%s/%s", etcd.SelfConfig.Foreign,progess.SchdulePipeline.Id)
	getResponse, err := etcd.Client.KV.Get(context.Background(), foreignKey)
	if len(getResponse.Kvs) == 0 {
		return
	}
	txnResponse, err := etcd.Client.Txn(context.Background()).If(clientv3.Compare(clientv3.Value(foreignKey), "=", string(getResponse.Kvs[0].Value[:]))).
		Then(clientv3.OpPut(foreignKey, progess.SchdulePipeline.Nodes)).
		Commit()
	if err != nil {
		return
	}
	if !txnResponse.Succeeded {
		return
	}

	worker_manager.WorkerManagerInstance.CountSub(workerId)
}

