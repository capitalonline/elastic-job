package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/service"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
)

type (
	Service struct {
		instance   *service.Instance
		instanceId string
		leaseID    clientv3.LeaseID
		close      chan struct{}
		wg         sync.WaitGroup
	}
)


func NewService(instance *service.Instance) (*Service, error) {
	return &Service{
		close:    make(chan struct{}),
		instance: instance,
	}, nil
}

// Register service
func (service *Service) Register(ttlSecond int64) error {
	logrus.Printf("worker ID为: %s", service.instance.Id)
	res, err := etcd.Client.Grant(context.TODO(), ttlSecond)
	if err != nil {
		return err
	}

	service.leaseID = res.ID

	val, err := json.Marshal(&service.instance)
	if err != nil {
		logrus.Println(err)
	}
	log.Printf("worker path: %s", etcd.SelfConfig.Service)
	key := fmt.Sprintf("%s/%s", etcd.SelfConfig.Service, service.instance.Id)

	if _, err = etcd.Client.Put(context.TODO(), key, string(val), clientv3.WithLease(service.leaseID)); err != nil {
		return err
	}

	logrus.Printf("启动成功, ID为: %s", service.instance.Id)

	ch, err := etcd.Client.KeepAlive(context.TODO(), service.leaseID)
	if nil != err {
		return err
	}

	service.wg.Add(1)
	defer service.wg.Done()

	for {
		select {
		case <-service.close:
			return service.revoke()
		case <-etcd.Client.Ctx().Done():
			return errors.New("server closed")
		case c, ok := <-ch:
			if !ok {
				logrus.Println(c.String())
				return service.revoke()
			}
		}
	}
}

// Stop service
func (service *Service) Stop() {
	close(service.close)
	service.wg.Wait()
	if err := etcd.Client.Close(); err != nil {
		logrus.Println(err)
	}
}

// Revoke service leaseID
func (service *Service) revoke() error {
	_, err := etcd.Client.Revoke(context.TODO(), service.leaseID)
	if err != nil {
		logrus.Printf("[discovery] Service revoke %s error: %s", service.instance.Id, err.Error())
	} else {
		logrus.Printf("[discovery] Service revoke successfully %s", service.instance.Id)
	}

	return err
}
