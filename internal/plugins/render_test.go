package plugins

import (
	"github.com/mongodb-job/config"
	"github.com/mongodb-job/models"
	"log"
	"strings"
	"testing"
)
func init() {
	config.Conf = &config.Config{
		Database: config.Database{
			Host: "101.251.219.226",
			Port: 3306,
			Name: "cds_mongo",
			User: "root",
			Pass: "123Abc,.;",
			Char: "utf8mb4",
		},
		Etcd: config.Etcd{},
	}
	models.Connection()

}
func TestCaculateRemoteS3Path(t *testing.T) {
	remotePath := "data/data-backup/CustomerID/strings/physicsbackup/rsn/2020-08-12_17-30-00.tar.gz "
	index := strings.Split(remotePath, "/")
	log.Printf("CaculateRemoteS3Path remotePath : %s, rsn : %s", remotePath, "rsn")
	s := append(index[:len(index)-1], index[len(index)-1])
	log.Printf("CaculateRemoteS3Path s : %s", s)
	path := strings.Join(s, "/")
	path = strings.Replace(path, " ", "", -1)
	println(path)
}

func TestFindStubByPipeline(t *testing.T) {
	var result models.PipelineRuntimeData
	result.PipelineId = "80a00c59-220d-4573-bfe9-0249dae651dd"
	pipeline := FindStubByPipeline(result, 2, "remote_path", "Rpl-HJsQUr")
	log.Println("result:" +pipeline)
}
