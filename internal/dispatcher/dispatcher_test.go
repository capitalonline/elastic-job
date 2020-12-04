package dispatcher

import (
	"github.com/mongodb-job/internal/etcd"
	"testing"
)
func initClient(){
	etcd.NewClient(etcd.Etcd{
		Killer:    "/mjob/v1/killer",
		Locker:    "/mjob/v1/locker",
		Service:   "/mjob/v1/service",
		Pipeline:  "/mjob/v1/pipeline",
		Snapshot:  "/mjob/v1/snapshot",
		Plan:      "/mjob/v1/plan",
		Foreign:   "/mjob/v1/foreign",
		EndPoints: []string{"101.251.219.229:8070"},
		Timeout:   1000,
	})

}
func TestDispatcher_DispatchPipeline(t *testing.T) {

	initClient()


}
