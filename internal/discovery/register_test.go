package discovery

import (
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/service"
	"log"
	"testing"
)

func initClient(){
	etcd.NewClient(etcd.Etcd{
		Killer:    "/mjob/v1/killer",
		Locker:    "/mjob/v1/locker",
		Service:   "/mjob/v1/service",
		Pipeline:  "/mjob/v1/pipeline",
		Plan:    "/mjob/v1/config",
		EndPoints: []string{"101.251.219.229:8070"},
		Timeout:   1000,
	})
}

func TestService_Register(t *testing.T) {
	initClient()
	newService, e := NewService(&service.Instance{
		Id:          "00000001",
		Name:        "worker1",
		Host:        "0.0.0.0",
		Port:        8001,
		Mode:        "worker",
		Status:      "running",
		Version:     "0.1",
		Description: "worker for job",
	})
	if e != nil {
		log.Fatal(e)
	}

	e = newService.Register(5)
}