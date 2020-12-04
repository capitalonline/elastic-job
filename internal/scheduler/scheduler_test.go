package scheduler

import (
	"encoding/json"
	"github.com/mongodb-job/config"
	"github.com/mongodb-job/models"
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

		Api:  config.Api{
			GicNewMonitorSystemUrl: "http://localhost:50000",
			GicUserUrl:             "http://10.13.2.235:6003",
		},
	}
	models.Connection()

}

func TestNotifyGIC(t *testing.T) {
	pipelineStr := `{"id":"a523c5a0-7550-4923-a29d-40de68b63374","clusterid":"e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","user_id":"663686","name":"FROMMONGOSERVICEFULL_cluster_id:e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","sched_type":"cron","spec":"50 5 * * 0,1,2,3,4,5,6","status":1,"c_at":1600659521,"u_at":1600659521,"node":"","steps":[{"id":"0","pipeline_id":"a523c5a0-7550-4923-a29d-40de68b63374","step":0,"c_at":1600659521,"u_at":1600659521,"task":{"id":"0","name":"update cluster e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","mode":"http","url":"http://mongodb-service/inner/v1/updateCluster","method":"POST","content":"{\"cluster_id\":\"e4e7acd4-1ad7-42ff-9fad-f9cf463e1802\",\"status\":\"Backuping\"}","desc":"更新集群备份状态"}},{"id":"1","pipeline_id":"a523c5a0-7550-4923-a29d-40de68b63374","step":1,"c_at":1600659521,"u_at":1600659521,"task":{"id":"1","name":"backup Physics e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","mode":"ansible","script":"physics_backup.yaml","hosts":"202.202.0.20 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI vm_role=hidden private_ip=10.240.29.14 replica_set_name=Rpl-jNLDpg \n202.202.0.23 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI vm_role=hidden private_ip=10.240.29.17 replica_set_name=Rpl-WoYSoP \n202.202.0.27 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI vm_role=hidden private_ip=10.240.29.21 replica_set_name=Config \n","content":"{\"backup_path\":\"/data/data-backup/E889999/e4e7acd41ad742ff9fadf9cf463e1802/physicsbackup\",\"backup_host_ip\":\"0.0.0.0\",\"type\":\"physicsbackup\",\"tar_filename\":\"{{nowtime}}\",\"db_info_path\":\"\",\"host_ip\":\"100.131.0.45 oss-cnbj01.cdsgss.com\",\"super_user_info\":{\"username\":\"cds_root\",\"passwd\":\"S18OTaxzOwdYscJmKPiJ\"}}","desc":"Physics e4e7acd4-1ad7-42ff-9fad-f9cf463e1802"}},{"id":"2","pipeline_id":"a523c5a0-7550-4923-a29d-40de68b63374","step":2,"c_at":1600659521,"u_at":1600659521,"task":{"id":"2","name":"upload Physics e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","mode":"ansible","script":"upload.yaml","hosts":"202.202.0.20 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI backup_file={{find_rps_map . 1 \"Rpl-jNLDpg\" \"path\"}} remote_backup_file={{s3_remote . 1 \"Rpl-jNLDpg\"}} backup_host_ip=202.202.0.20 vm_role=hidden private_ip=10.240.29.14 replica_set_name=Rpl-jNLDpg \n202.202.0.23 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI backup_file={{find_rps_map . 1 \"Rpl-WoYSoP\" \"path\"}} remote_backup_file={{s3_remote . 1 \"Rpl-WoYSoP\"}} backup_host_ip=202.202.0.23 vm_role=hidden private_ip=10.240.29.17 replica_set_name=Rpl-WoYSoP \n202.202.0.27 ansible_ssh_pass=V!1RoLxsNugDphnyhqGI backup_file={{find_rps_map . 1 \"Config\" \"path\"}} remote_backup_file={{s3_remote . 1 \"Config\"}} backup_host_ip=202.202.0.27 vm_role=hidden private_ip=10.240.29.21 replica_set_name=Config \n","content":"{\"access_key\":\"c13429db9c7b5b4d890342437ade933e\",\"secret_key\":\"572989b801a55735b3a21a37d8ca8486\",\"host_ip\":\"100.131.0.45 oss-cnbj01.cdsgss.com\",\"total_host_ip\":\"100.131.0.45 mongo-pre.ae327e0452c545349e4b6b41478d72b5.oss-cnbj01.cdsgss.com\",\"endpoint_url\":\"http://oss-cnbj01.cdsgss.com\",\"bucket_name\":\"mongo-pre\",\"action\":\"upload\"}","desc":"upload Physics e4e7acd4-1ad7-42ff-9fad-f9cf463e1802"}},{"id":"3","pipeline_id":"a523c5a0-7550-4923-a29d-40de68b63374","step":3,"c_at":1600659521,"u_at":1600659521,"task":{"id":"3","name":"update cluster e4e7acd4-1ad7-42ff-9fad-f9cf463e1802","mode":"http","url":"http://mongodb-service/inner/v1/updateCluster","method":"POST","content":"{\"cluster_id\":\"e4e7acd4-1ad7-42ff-9fad-f9cf463e1802\",\"status\":\"Running\"}","desc":"还原集群备份状态"}}]}`
	b := models.Pipeline{}

	_ = json.Unmarshal([]byte(pipelineStr), &b)
	notifyGIC("exit status 1", b)
}