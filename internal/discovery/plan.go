package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/scheduler"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func WatchPlan(worker *service.Instance) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var curRevision int64 = 0
	listenKey := fmt.Sprintf("%s/%s/", etcd.SelfConfig.Plan, worker.Id)
	rangeResp, err := etcd.Client.Get(context.TODO(), listenKey, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	curRevision = rangeResp.Header.Revision + 1

	for _, obj := range rangeResp.Kvs {
		var progess models.Progress
		if err := json.Unmarshal(obj.Value, &progess); err != nil {
			logrus.Errorf("Parse etcd to progress failed, error: %v, string: %s", err, string(obj.Value))
			continue
		}
        if progess.SchdulePipeline.ScheduleType == "delay" {
			now := time.Now()
			if time.Unix(progess.SchdulePipeline.ExecTime, 0).Before(now){
				continue
			}
		}
		logrus.Println("更新任务WatchPlan : ", progess.SchdulePipeline.Id,progess.SchdulePipeline.Spec)
		scheduler.Instance.DispatchEvent(&scheduler.Event{
			Type:     scheduler.PUT,
			Pipeline: progess.SchdulePipeline,
		})
	}

	watchChan := etcd.Client.Watch(context.TODO(), listenKey, clientv3.WithPrefix(), clientv3.WithRev(curRevision), clientv3.WithPrevKV())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			var progess models.Progress
			switch event.Type {
			case mvccpb.PUT:
				if err := json.Unmarshal(event.Kv.Value, &progess); err != nil {
					logrus.Println(err)
					continue
				}
				if progess.Status == models.DISPATCH {
					scheduler.Instance.DispatchEvent(&scheduler.Event{
						Type:     scheduler.PUT,
						Pipeline: progess.SchdulePipeline,
					})
					// TODO 取代成MYSQL
					progess.Status = models.Knowleage
					updateStatusProgress, _ := json.Marshal(progess)
					_, _ = etcd.Client.Put(context.TODO(), string(event.Kv.Key), string(updateStatusProgress))

					models.ExecUpdateWorkerPipelineStatus(models.Knowleage, 0, "","",progess.SchdulePipeline.Id,worker.Id)

				} else if progess.Status == models.UPDATETIME {
					logrus.Println("更新任务: ", progess.SchdulePipeline.Id,progess.SchdulePipeline.Spec)
					scheduler.Instance.DispatchEvent(&scheduler.Event{
						Type:     scheduler.PUT,
						Pipeline: progess.SchdulePipeline,
					})
				}
			case mvccpb.DELETE:
				if err := json.Unmarshal(event.PrevKv.Value, &progess); err != nil {
					logrus.Println(err)
					continue
				}
				scheduler.Instance.DispatchEvent(&scheduler.Event{
					Type:     scheduler.DEL,
					Pipeline: progess.SchdulePipeline,
				})
			}
		}

		time.Sleep(1 * time.Second)
	}
}

//func WatchKiller() {
//	var curRevision int64 = 0
//
//	for {
//		rangeResp, err := discover.Client.Get(context.TODO(), config.Conf.Etcd.Pipeline, clientv3.WithPrefix())
//
//		if err != nil {
//			continue
//		}
//		curRevision = rangeResp.Header.Revision + 1
//		break
//	}
//
//	watchChan := discover.Client.Watch(context.TODO(), "", clientv3.WithPrefix(), clientv3.WithRev(curRevision))
//	for watchResp := range watchChan {
//		for _, event := range watchResp.Events {
//			var pipeline models.Pipeline
//			if err := json.Unmarshal(event.Kv.Value, &pipeline); err != nil {
//				log.Println(err)
//			}
//
//			switch event.Type {
//			case mvccpb.PUT:
//			case mvccpb.DELETE:
//			}
//		}
//	}
//
//}
