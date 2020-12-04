package models

import (
	"github.com/mongodb-job/config"
	uuid "github.com/satori/go.uuid"
	"log"
	"testing"
	"time"
	"xorm.io/builder"
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

func TestPipelineRecordsADD(t *testing.T) {
	for i := 0; i < 100; i++ {
		a := PipelineRecords{
			Id:           uuid.NewV4().String(),
			PipelineId:   "7d5f86a5-5b7c-45ac-9e66-c4b1fa64c509",
			NodeId:       uuid.NewV4().String(),
			WorkerName:   "test-node-1",
			ScheduleType: "cron",
			Spec:         "*/5 * * * *",
			Status:       1,
			Duration:     500,
			BeginWith:    time.Now().Unix(),
			FinishWith:   time.Now().Unix(),
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		_ = a.Store()
	}

}

func TestPipelineRecordsDelete(t *testing.T) {
	Engine.Delete(&PipelineRecords{Id: "0c512bf0-b7bd-4ba3-8d31-37b1258c851d"})
}

func TestPipelineRecordsSearch(t *testing.T) {
	size := 40
	page := 20
	search := "7d5f86a5-5b7c-45ac-9e66-c4b1fa64c509"
	logs := make([]PipelineRecords, 0)
	if search != "" {
		total, _ := Engine.Where(builder.Eq{"pipeline_id": search}).Limit(size, page).Desc("created_at").FindAndCount(&logs)
		log.Println(total)
		log.Println(len(logs))
	} else {
		_, _ = Engine.Limit(size, page).Desc("created_at").FindAndCount(&logs)
	}

}

func TestTaskRecordsADD(t *testing.T) {
	for i := 0; i < 100; i++ {
		a := TaskRecords{
			PipelineRecordId: uuid.NewV4().String(),
			TaskId:           "TfdsDFDSfmgvk",
			NodeId:           uuid.NewV4().String(),
			TaskName:         "B诶根",
			WorkerName:       "node-test-1",
			Mode:             "test",
			Url:              "http://baidu.com",
			Method:           "get",
			Script:           "asdas.yaml",
			Hosts:            "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogon",
			Content:          "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogon",
			Timeout:          0,
			Retries:          0,
			Status:           "finish",
			Result:           "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogonssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcI5eJXXf0+NYIyfVsxCdEXpr70zti9XcwDIFT9QZ4bj5YwOdmFBbpppoUnnvPT68PtqrUZ+PiG59yUUlI4pTaz66wOeopTk7P7B0cMjclN9rL1iTt6GSjhJmeuBry7eeSi16d1A1BoZTxL/erGo5NJ/yUHs2PIha8QJYu8xcNf6dCmGqpBpPCcrBsDq3vlmyQkS1Yxa7+YtAwgYS30snqL8fgTPHFGReC53oQLFB3Of8e+9/XWGV7CnyibBH73lwvyw0wrM7jnDzPrtJ1AhGGGcr1vrC53VonKE164GBiboIKgY0Brrarj1Ty3xHbXhqPd8t+AP7dQWBevYdXlqn9 zhangzhen@bogon",
			Duration:         100,
			BeginWith:        time.Now().Unix(),
			FinishWith:       time.Now().Unix(),
			CreatedAt:        time.Now().Unix(),
		}

		if err := a.Store(); err != nil {
			log.Fatal(err)
		}
	}
}

func TestTaskRecordsDelete(t *testing.T) {

}

func TestTaskRecordsSearch(t *testing.T) {

}
