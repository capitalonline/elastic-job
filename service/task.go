package task

import (
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"time"
)

func AutoDelete() {
	logrus.Info(">>>>>>>>>> Begin auto delete")
	ticker := time.NewTicker(40 * time.Second)
	defer func() {
		logrus.Info(">>>>>>>>>> Exit auto delete")
		ticker.Stop()
		logrus.Info(">>>>>>>>>> Exit auto delete success")
	}()

	for {
		select {
		case <-ticker.C:
			deleteExpireBackup()
		}
	}
}

func deleteExpireBackup() {
	allBackup, err := models.FindAllBackupRecords()
	if err != nil {
		logrus.Errorf("Find all auto full physics backup failed, error: %v", err.Error())
		return
	}
	for _, record := range allBackup {
		if record.StartTime.AddDate(0,0,record.Deadline).Before(time.Now()){
			logrus.Info(">>>>>>>>>> Auto delete expire backup: %+v", record)

			deleteBackup, err := models.DeleteNoVaildBackupRecord(record.BackupId, record.UserId)
			if err != nil {
				logrus.Error(">>>>>>>>>> Auto delete backup error: %v, param: %+v", err, record)
				continue
			}
			logrus.Println(">>>>>>>>>>" + deleteBackup)
		}
	}
}

func AutoDeleteIncrBackup() {
	logrus.Info(">>>>>>>>>> Begin auto delete incr")
	ticker := time.NewTicker(23 * time.Second)
	defer func() {
		logrus.Info(">>>>>>>>>> exit auto delete")
		ticker.Stop()
		logrus.Info(">>>>>>>>>> exit auto delete success")
	}()

	for {
		select {
		case <-ticker.C:
			deleteExpireIncrBackup()
		}
	}
}

func deleteExpireIncrBackup() {
	allIncrBackup, err := models.FindAllIncrBackupRecords()
	if err != nil {
		logrus.Error("Find all incr backup failed")
		return
	}
	for _, record := range allIncrBackup {
		recordTime := time.Unix(record.StartTime, 10)
		if recordTime.AddDate(0,0,record.Deadline).Before(time.Now()){
			logrus.Info(">>>>>>>>>> auto delete expire incr backup: %+v", record)

			deleteBackup, err := models.DeleteNoVaildIncrBackupRecord(record.BackupId, record.UserId)
			if err != nil {
				logrus.Error(">>>>>>>>>> auto delete incr backup error: %v, param: %+v", err, record)
				continue
			}
			logrus.Println(">>>>>>>>>>" + deleteBackup)
		}
	}
}

func AutoDeleteTempCluster(){
	logrus.Info(">>>>>>>>>> Begin auto delete temp")
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		logrus.Info(">>>>>>>>>> exit auto delete temp")
	}()
	for {
		select {
		case <-ticker.C:
			deleteTempCls()
		}
	}
}

func deleteTempCls(){
	clusters, _ := models.FindTempCluster()
	for _, cluster := range clusters {
		if cluster.CreateTime.AddDate(0,0,2).Before(time.Now()) && cluster.IsValid==2{
			logrus.Info(">>>>>>>>>> auto delete temp cluster: %+v", cluster)
			tempcluster := models.ToDeleteTempcluster(cluster.UserId, cluster.CustomerId, cluster.ClusterId)
			logrus.Info(">>>>>>>>>> auto delete temp result: %v, param: %+v", tempcluster, cluster)
		}
	}

}