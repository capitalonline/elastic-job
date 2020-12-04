package plugins

import (
	"bytes"
	"encoding/json"
	"github.com/mongodb-job/models"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"
	"xorm.io/builder"
)

func Render(tpl string, pr models.PipelineRuntimeData) string {
	var b bytes.Buffer
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{"nowtime": Nowtime, "find_rps_map": FromReplicaSetBackupMap, "s3_remote": CaculateRemoteS3Path, "find_index_map": FromKVByIndex, "stub": FindStubByPipeline, "nowunix": NowTimestamp}).Parse(tpl))
	if err := tmpl.Execute(&b, pr); err != nil {
		log.Println(err)
	}
	return b.String()
}

func Nowtime() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

func FromReplicaSetBackupMap(pr models.PipelineRuntimeData, step int, rsn string, key string) string {
	for _, rmap := range pr.StepPR[step].Result {
		if rmap["replica_set_name"] == rsn {
			return strings.TrimSpace(rmap[key])
		}
	}
	return ""
}

func CaculateRemoteS3Path(pr models.PipelineRuntimeData, step int, rsn string) string {
	path := FromReplicaSetBackupMap(pr, step, rsn, "path")
	remotePath := strings.TrimLeft(path, "/")
	index := strings.Split(remotePath, "/")
	log.Printf("CaculateRemoteS3Path remotePath : %s, rsn : %s", remotePath, rsn)
	s := append(index[:len(index)-1],rsn ,index[len(index)-1])
	log.Printf("CaculateRemoteS3Path s : %s", s)
	path = strings.Join(s, "/")
	//return strings.Replace(path, " ", "", -1)
	return strings.TrimSpace(path)
}

func FromKVByIndex(pr models.PipelineRuntimeData, step int, index int, key string) string {
	s, ok := pr.StepPR[step].Result[index][key]
	if !ok {
		return ""
	} else {
		return s
	}
}

func FindStubByPipeline(pr models.PipelineRuntimeData, step int, key string, rsn string) string {
	var r models.PipelineRecords
	get, err2 := models.Engine.Where(builder.Eq{"pipeline_id": pr.PipelineId, "status": 1}).Desc("finish_with").Get(&r)
	if err2 != nil {
		//return strconv.FormatInt(time.Now().Add(-time.Minute * 50).Unix(), 10)
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	if !get {
		//return strconv.FormatInt(time.Now().Add(-time.Minute * 50).Unix(), 10)
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	var p models.TaskRecords
	yes, err := models.Engine.Where(builder.Eq{"pipeline_record_id": r.Id, "task_id": step, "status": "finished"}).Desc("finish_with").Limit(1).Get(&p)
	if err != nil || !yes{
		//return strconv.FormatInt(time.Now().Add(-time.Minute * 50).Unix(), 10)
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	temp := make([]map[string]string, 0)
	_ = json.Unmarshal([]byte(p.Result), &temp)
	for _, m := range temp {
		if m["replica_set_name"] == rsn {
			return m[key]
		}
	}
	//return strconv.FormatInt(time.Now().Add(-time.Minute * 50).Unix(), 10)
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func NowTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}