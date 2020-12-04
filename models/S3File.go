package models

import "xorm.io/builder"

type S3File struct {
	FileId         string
	BackupId       string
	Path           string
	Size           int
	Md5            string
	CreateTime     int64
	Valid          int
	ReplicasetName string
}

func (sf *S3File) TableName() string {
	return "s3file"
}

func (sf *S3File) Store() error {
	_, err := Engine.InsertOne(sf)
	return err
}


type IncrementalBackupRecord struct {
	BackupId       string
	BackupName     string
	Path           string
	OplogStartTime int64
	OplogEndTime   int64
	StartTime      int64
	EndTime        int64
	Status         int
	CreateTime     int64
	Valid          int
	Deadline       int
	BackupAction   string
	BackupType     string
	ClusterId      string
	UserId         string
	ReplicaSetName string
	Size           int
	Md5            string
}

func (sf *IncrementalBackupRecord) TableName() string {
	return "incremental_backup_record"
}

func (sf *IncrementalBackupRecord) Store() error {
	_, err := Engine.InsertOne(sf)
	return err
}

func FindAllIncrBackupRecords() (bks []IncrementalBackupRecord, err error) {
	bks = make([]IncrementalBackupRecord, 0)
	err = Engine.Where(builder.Eq{"valid": 1, "backup_action": "Auto"}).Find(&bks)
	return
}

