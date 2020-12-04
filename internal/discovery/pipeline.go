package discovery

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/mongodb-job/internal/dispatcher"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"time"
)

// master用来监听新增流水线
func WatchPipelines() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic: %+v", r)
			return
		}
	}()
	var curRevision int64 = 0
	historyPipelines, err := etcd.Client.Get(context.TODO(), etcd.SelfConfig.Pipeline, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	logrus.Infof("Current task count is %d", historyPipelines.Count)
	curRevision = historyPipelines.Header.Revision + 1

	watchChan := etcd.Client.Watch(context.TODO(), etcd.SelfConfig.Pipeline, clientv3.WithPrefix(), clientv3.WithRev(curRevision), clientv3.WithPrevKV())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			var pipeline models.Pipeline
			switch event.Type {
			case mvccpb.PUT:
				logrus.Info("Enter add pipeline")
				if err := json.Unmarshal(event.Kv.Value, &pipeline); err != nil {
					logrus.Println(err)
					continue
				}
				logrus.Infof("Begin deal pipeline, pipeline id: %s, name: %s, status: %d", pipeline.Id, pipeline.Name, pipeline.Status)
				if pipeline.Status==models.CREATE{
					dispatcher.DispatcherInstance.DispatchPipeline(pipeline)
				}
				if pipeline.Status==models.UPDATE{
					dispatcher.DispatcherInstance.ModifyPipeline(pipeline)
				}

			case mvccpb.DELETE:
				logrus.Info("Enter delete pipeline")
				if err := json.Unmarshal(event.PrevKv.Value, &pipeline); err != nil {
					logrus.Println(err)
					continue
				}
				dispatcher.DispatcherInstance.DeletePipeline(pipeline)

			}
		}

		time.Sleep(1 * time.Second)
	}
}