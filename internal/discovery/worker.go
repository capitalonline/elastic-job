package discovery

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/internal/worker_manager"
	"github.com/sirupsen/logrus"
)

func WatchWorker() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Println(r)
		}
	}()
	var curRevision int64 = 0

	for {
		rangeResp, err := etcd.Client.Get(context.TODO(), etcd.SelfConfig.Service, clientv3.WithPrefix())

		if err != nil {
			logrus.Error(err)
			continue
		}
		// 从当前版本开始订阅
		for _, kv := range rangeResp.Kvs {
			var worker service.Instance
			if err := json.Unmarshal(kv.Value, &worker); err != nil {
				logrus.Errorf("Init worker failed, error: %s", err.Error())
				continue
			}
			if worker.Mode == service.WORKER {
				worker_manager.WorkerManagerInstance.AddWorker(&worker)
			}
		}
		curRevision = rangeResp.Header.Revision + 1
		break
	}

	watchChan := etcd.Client.Watch(context.Background(), etcd.SelfConfig.Service, clientv3.WithPrefix(), clientv3.WithRev(curRevision), clientv3.WithPrevKV())

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			var worker service.Instance
			switch event.Type {
			case mvccpb.PUT:
				logrus.Infof("Enter add worker")
				if err := json.Unmarshal(event.Kv.Value, &worker); err != nil {
					logrus.Errorf("Add worker failed, error: %s", err.Error())
					continue
				}
				if worker.Mode == service.WORKER {
					worker_manager.WorkerManagerInstance.AddWorker(&worker)
				}
			case mvccpb.DELETE:
				logrus.Infof("Enter rm worker")
				if err := json.Unmarshal(event.PrevKv.Value, &worker); err != nil {
					logrus.Errorf("Remove worker failed, error: %s", err.Error())
					continue
				}
				worker_manager.WorkerManagerInstance.Down(worker.Id)
			}

		}
	}
}
