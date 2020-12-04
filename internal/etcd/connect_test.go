package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"testing"
)

func TestNewClient(t *testing.T) {
	NewClient(Etcd{
		Killer:    "/mjob/v1/killer",
		Locker:    "/mjob/v1/locker",
		Service:   "/mjob/v1/service",
		Pipeline:  "/mjob/v1/pipeline",
		Plan:    	"/mjob/v1/plan",
		Snapshot:  "/mjob/v1/snatshot",
		Foreign:   "/mjob/v1/foreign",
		EndPoints: []string{"101.251.219.229:8070"},
		Timeout:   1000,
	})

	//list, exection := Client.MemberList(context.Background())
	//if exection != nil {
	//	log.Fatal(exection)
	//}
	//log.Println(list)

	foreignKey := fmt.Sprintf("%s/%s", "/mjob/v1/foreign","eb26a832-7532-40b8-9412-e2aaa2a5e80b")
	getResponse, err := Client.KV.Get(context.Background(), foreignKey)
	if len(getResponse.Kvs) == 0 {
		return
	}

	oldNodeId := string(getResponse.Kvs[0].Value[:])
	txnResponse, err := Client.Txn(context.Background()).If(clientv3.Compare(clientv3.Value(foreignKey), "=", oldNodeId)).
		Then(clientv3.OpPut(foreignKey, "5a14fd18-b0d0-474d-911c-452ae8f95a10")).
		Commit()
	if err != nil {
		return
	}
	if !txnResponse.Succeeded {
		return
	}
}
