package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/mongodb-job/config"
	"log"
	"time"
)

var (
	err        error
	SelfConfig config.Etcd
	Client     *clientv3.Client
)

// New ETCD V3 Client
func NewClient(c config.Etcd) {
	if Client, err = clientv3.New(clientv3.Config{
		Endpoints:   c.EndPoints,
		DialTimeout: 10 * time.Second,
	}); err != nil {
		log.Println(err)
	}
	SelfConfig = c
}
