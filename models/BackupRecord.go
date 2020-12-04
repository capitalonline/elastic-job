package models

import (
	"fmt"
	"github.com/mongodb-job/config"
	"io/ioutil"
	"net/http"
	"time"
)
import "xorm.io/builder"

const (
	_deleteBackup     = "/inner/v1/delete_backup"
	_deleteIncrBackup = "/inner/v1/delete_incr_backup"
)

type BackupRecord struct {
	BackupId         string     `xorm:"not null pk comment('ID') CHAR(36)"`
	BackupName       string     `xorm:"not null comment('备份名称') VARCHAR(255)"`
	MetedataPath     string     `xorm:"comment('元数据路径，逻辑备份使用') VARCHAR(255)"`
	BackupType       string     `xorm:"not null comment('备份类型') VARCHAR(45)"`
	BackupAction     string     `xorm:"not null comment('备份方式') VARCHAR(45)"`
	Deadline         int        `xorm:"not null comment('保留时间') INT(30)"`
	Status           string     `xorm:"not null comment('状态') VARCHAR(45)"`
	FinishTime       *time.Time `xorm:"comment('完成时间') DATETIME"`
	DeleteTime       *time.Time `xorm:"comment('删除时间') DATETIME"`
	DeleteFinishTime *time.Time `xorm:"comment('删除完成时间') DATETIME"`
	UploadStatus     string     `xorm:"not null comment('上传状态') VARCHAR(45)"`
	StartTime        time.Time  `xorm:"comment('开始时间') DATETIME"`
	Valid            int        `xorm:"not null comment('有效') INT(22)"`
	ClusterId        string     `xorm:"not null comment('集群id') VARCHAR(45)"`
	UserId           string     `xorm:"not null comment('用户id') VARCHAR(45)"`
}

func (obj *BackupRecord) TableName() string {
	return "backup_record"
}

func FindAllBackupRecords() (bks []BackupRecord, err error) {
	bks = make([]BackupRecord, 0)
	err = Engine.Where(builder.Eq{"valid": 1, "backup_action": "Auto"}.
		And(builder.In("status", "Recovery OK", "OK"))).
		Find(&bks)
	return
}

func DeleteNoVaildBackupRecord(backId, userId string) (result string, err error) {
	url := fmt.Sprintf("http://%s%s?backup_id=%s&user_id=%s", config.Conf.App.Service, _deleteBackup, backId, userId)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	result = string(body)
	return
}

func DeleteNoVaildIncrBackupRecord(backId, userId string) (result string, err error) {
	url := fmt.Sprintf("http://%s%s?backup_id=%s&user_id=%s", config.Conf.App.Service, _deleteIncrBackup, backId, userId)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	result = string(body)
	return
}

func (obj *BackupRecord) Store() error {
	_, err := Engine.InsertOne(obj)
	return err
}
