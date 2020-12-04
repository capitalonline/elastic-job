package worker_manager

import (
	"github.com/mongodb-job/internal/service"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
)

type (
	JobNode struct {
		*service.Instance
		JobCount int
	}
	WorkerManager struct {
		clients map[string]*JobNode
		sync.Mutex
	}
)

var (
	WorkerManagerInstance *WorkerManager
)

func init() {
	WorkerManagerInstance = &WorkerManager{
		clients: make(map[string]*JobNode),
	}
}

func (wm *WorkerManager) SelectWorker() *service.Instance {
	wm.Lock()
	defer wm.Unlock()
	max := 9999
	var result *service.Instance
	for _, node := range wm.clients {
		if node.Status == "up" {
			if max > node.JobCount {
				max = node.JobCount
				result = node.Instance
			}
		}
	}
	return result
}

func (wm *WorkerManager) AddWorker(w *service.Instance) {
	wm.Lock()
	defer wm.Unlock()
	if node := wm.clients[w.Id]; node != nil {
		if node.Status == "down" {
			node.Status = "up"
		}
		logrus.Warnf("Worker %s exist", w.Id)
		return
	}
	wm.clients[w.Id] = &JobNode{
		Instance: w,
		JobCount: 0,
	}
	logrus.Info("Add worker success, worker id: %s, host: %s", w.Id, w.Host)
	return
}
func (wm *WorkerManager) RemoveWorker(workerId string) {
	wm.Lock()
	defer wm.Unlock()
	node := wm.clients[workerId]
	if node == nil {
		return
	}
	if node.JobCount != 0{
		log.Println("Warn: 节点存在任务就被删除了", node.Id)
	}
	delete(wm.clients, workerId)
	return
}

func (wm *WorkerManager) CountAdd(workerId string) {
	wm.Lock()
	defer wm.Unlock()
	if node := wm.clients[workerId]; node == nil {
		return
	} else {
		node.JobCount++
	}
	return
}

func (wm *WorkerManager) CountSub(workerId string) {
	wm.Lock()
	defer wm.Unlock()
	if node := wm.clients[workerId]; node == nil {
		return
	} else {
		node.JobCount--
	}
	return
}

func (wm *WorkerManager) Have(workerId string) *service.Instance {
	wm.Lock()
	defer wm.Unlock()
	if node := wm.clients[workerId]; node == nil {
		return nil
	} else {
		return node.Instance
	}
}
func (wm *WorkerManager) All() map[string]*JobNode {
	//wm.Lock()
	//defer wm.Unlock()
	return wm.clients
}

func (wm *WorkerManager) Down(workerId string){
	wm.Lock()
	defer wm.Unlock()
	if node := wm.clients[workerId]; node == nil {
		logrus.Warnf("Wish down node not exist, worker id: %s", workerId)
		return
	} else {
		wm.clients[workerId].Status = "down"
	}
}
