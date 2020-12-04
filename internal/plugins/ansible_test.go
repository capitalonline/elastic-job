package plugins

import (
	"bytes"
	"fmt"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	"log"
	"testing"
	"text/template"
)

func TestRender(t *testing.T) {
	const tbl = `
		{
			"backup_path":"/data/backup/a",
			"backup_host_ip":"a",
			"type":"a",
			"tar_filename":"{{nowtime}}"
 			"content": "{"backup_path":"/data/data-backup/CustomerID/strings/physicsbackup","type":"PhysicsBackup","backup_host_ip":"0.0.0.0","tar_filename":"{{nowtime}}","host_ip":"100.131.0.45 mongo-pre.ae327e0452c545349e4b6b41478d72b5.oss-cnbj01.cdsgss.com","start_time":"{{stub . 1 "end_time" "Rpl-HJsQUr"}}","end_time":"{{nowunix}}","super_user_info":{"username":"cds_root","passwd":"S18NAMXPkDis"}}"	
		}				
	`

	const tpl = `123 ansible_ssh_pass=123\n345 ansible_ssh_pass=345`

	var b bytes.Buffer

	runtimeData := models.PipelineRuntimeData{
		PipelineId: "b0a00c59-220d-4573-bfe9-0249dae651dd",
		StepPR:     make(models.StepParamResult),
	}
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{"nowtime":Nowtime, "fromcontext":FromReplicaSetBackupMap, "stub": FindStubByPipeline, "nowunix": NowTimestamp}).Parse(tbl))
	if err := tmpl.Execute(&b, runtimeData); err != nil {
		log.Fatal(err)
	}
	print(b.String())
}

func TestAnsible(t *testing.T) {
	//service.Runtime = &service.Instance{
	//	Id:          "test1",
	//	Name:        "test1",
	//	Host:        "0.0.0.0",
	//	Port:        30,
	//	Mode:        "worker",
	//	Status:      "up",
	//	Version:     "0.1",
	//	Description: "n",
	//}
	//a := Ansible{
	//	OrginHost:    "127.0.0.1 ansible_ssh_pass=zhangzhen \n127.0.0.1 ansible_ssh_pass=zhangzhen \n",
	//	OrginContent: "",
	//	OriginScript: "ping.yaml",
	//}
	//result := make(models.StepParamResult)
	//a.Exec(context.TODO(), 1, result)
}

func TestAnsible1(t *testing.T) {
	service.Runtime = &service.Instance{
		Id:          "test1",
		Name:        "test1",
		Host:        "0.0.0.0",
		Port:        30,
		Mode:        "worker",
		Status:      "up",
		Version:     "0.1",
		Description: "n",
	}
	result, err := oexec("", "", "/Users/zhangzhen/go/src/github.com/mongodb-job/orHJik/hosts", "/Users/zhangzhen/go/src/github.com/mongodb-job/script/yaml_file/ping.yaml", "")
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(result)
}
