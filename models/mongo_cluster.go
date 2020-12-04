package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mongodb-job/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
	"xorm.io/builder"
)

const (
	_removeTempBackup = "/inner/v1/delete_temp"
	_updateClusterStatusRunning = "UPDATE mongo_cluster SET `status` = 'Running' WHERE `cluster_id` = ?"
)
type MongoCluster struct {
	ClusterId     string    `gorm:"primary_key;column:cluster_id;type:varchar(36);not null"` // 主键
	ClusterName   string    `gorm:"column:cluster_name;type:varchar(255);not null"`
	CustomerId    string    `gorm:"column:customer_id;type:varchar(36);not null"`
	UserId        string    `gorm:"column:user_id;type:varchar(36);not null"`
	UserSiteId    string    `gorm:"column:user_site_id;type:varchar(36);not null"`
	UserSiteName  string    `gorm:"column:user_site_name;type:varchar(72);not null"`
	UserAppId     string    `gorm:"column:user_app_id;type:varchar(36);not null"`
	UserAppName   string    `gorm:"column:user_app_name;type:varchar(72);not null"`
	UserPipeId    string    `gorm:"column:user_pipe_id;type:varchar(36)"`
	Status        string    `gorm:"column:status;type:varchar(30);not null"`
	ClusterType   string    `gorm:"column:cluster_type;type:varchar(36)"`
	Version       string    `gorm:"column:version;type:varchar(36)"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null"`
	IsValid       int8      `gorm:"column:is_valid;type:tinyint(1);not null"`
	Passwd        string    `gorm:"column:passwd;type:varchar(36)"`
	Config        string    `gorm:"column:config;type:varchar(255)"`
	Ips           string    `gorm:"column:ips;type:varchar(36)"`
	ImageId       string    `gorm:"column:image_id;type:varchar(36);not null"`
	ManagerSiteId string    `gorm:"column:manager_site_id;type:varchar(36)"`
	TaskId        string    `gorm:"column:task_id;type:varchar(2048);not null"`
	Detail        string    `gorm:"column:detail;type:varchar(512)"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null"`
	ManagerAppId  string    `gorm:"column:manager_app_id;type:varchar(36)"`
	GpnPipeId     string    `gorm:"column:gpn_pipe_id;type:varchar(36)"`
	SuborderId    string    `gorm:"column:suborder_id;type:varchar(36);not null"`
	GoodId        int       `gorm:"column:good_id;type:int(11);not null"`
	Port          int       `gorm:"column:port;type:int(11);not null"`
	OplogSize     int       `gorm:"column:oplog_size;type:int(11);not null"`
	GoodName      string    `gorm:"column:good_name;type:varchar(128);not null"`
	SuperPwd      string    `gorm:"column:super_pwd;type:varchar(128);not null"`
	ConnectionUri string    `gorm:"column:connection_uri;type:varchar(128);not null"`
}

type InnerDeleteTempClusterVO struct {
	ChildCsId  string `json:"child_id" validate:"required"`
	UserId     string `json:"user_id"`
	CustomerId string `json:"customer_id"`
}


func (obj *MongoCluster) TableName() string {
	return "mongo_cluster"
}

func FindTempCluster()(cls []MongoCluster, err error) {
	cls = make([]MongoCluster, 0)
	err = Engine.Where(builder.Eq{"is_valid": 2}).Find(&cls)
	return
}

func FindClusterByID(clusterId string)(cls MongoCluster, err error) {
	_, err = Engine.Where(builder.Eq{"cluster_id": clusterId, "is_valid": 1}).Get(&cls)
	return
}

func ToDeleteTempcluster(userId, customerId, clusterId string)(result string){
	url := fmt.Sprintf("http://%s%s", config.Conf.App.Service, _removeTempBackup)
	client := &http.Client{}
	param := &InnerDeleteTempClusterVO{
		ChildCsId:  clusterId,
		UserId:     userId,
		CustomerId: customerId,
	}

	j, err := json.Marshal(param)
	if err != nil {
		return err.Error()
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(j))
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	result = string(body)
	return
}

func UpdateClusterStatusToRunning(clusterId string){
	if clusterId == ""{
		return
	}
	_, err := Engine.Exec(_updateClusterStatusRunning, clusterId)
	if err != nil {
		logrus.Errorf("db.Exec(%s) Failed, cluster_id = %s", _updateClusterStatusRunning, clusterId)
	}
}

