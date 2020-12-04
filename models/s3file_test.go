package models

import (
	"github.com/mongodb-job/config"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
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
	Connection()

}
func TestIncrementalBackupRecord_Store(t *testing.T) {
	record := IncrementalBackupRecord{
		BackupId:       uuid.NewV4().String(),
		BackupName:     "test",
		Path:           "/data/test",
		OplogStartTime: time.Now().Unix(),
		OplogEndTime:   time.Now().Unix(),
		StartTime:      time.Now().Unix(),
		EndTime:        time.Now().Unix(),
		Status:         1,
		CreateTime:     time.Now().Unix(),
		Valid:          1,
		Deadline:       8,
		BackupAction:   "M",
		BackupType:     "Auto",
		ClusterId:      uuid.NewV4().String(),
		UserId:         "U0000000000",
		ReplicaSetName: "RPL-serfs",
		Size:           0,
		Md5:            uuid.NewV4().String(),
	}

	record.Store()
}
