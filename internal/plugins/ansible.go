package plugins

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/internal/utils"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

const (
	//_tempPath = "/Users/zhangzhen/go/src/github.com/mongodb-job/"
	_tempPath = "/app/logs/temporary/backup/"
)

type (
	Ansible struct { //只存放脚本执行所需要的东西，业务数据全部从etcd获取
		OrginHost string
		parseHost string

		OrginContent string
		parseContent string

		OriginScript string
		scriptPath   string

		tempPath string
	}
)

func (a *Ansible) Exec(ctx context.Context, step int, pr models.PipelineRuntimeData) *models.TaskRecords {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	record := &models.TaskRecords{}
	a.parseContent = Render(a.OrginContent, pr)
	a.parseHost = Render(a.OrginHost, pr)
	acc := models.ParamResult{Param: models.StringMap{"content": a.parseContent, "host": a.parseHost, "script": a.scriptPath}}

	spath, err := findScriptDir()
	if err != nil {
		record.Status = "failed"
		record.Result = err.Error()
		return record
	}
	a.scriptPath = spath + a.OriginScript                              // "/user/a/" + a.yml
	a.tempPath = _tempPath + utils.RandStringBytesMaskImprSrcUnsafe(6) // "/temp/clustereId/asdasd"
	err = os.MkdirAll(a.tempPath, 0644)
	if err != nil {
		record.Status = "failed"
		record.Result = err.Error()
		return record
	}

	defer func() {
		os.RemoveAll(a.tempPath)
	}()
	hostPath := a.tempPath + "/hosts"
	err = utils.WriteFile(hostPath, a.parseHost)
	if err != nil {
		record.Status = "failed"
		record.Result = err.Error()
		return record
	}

	contentpath := ""
	if a.parseContent != ""{
		contentpath = a.tempPath + "/hosts.json"
		err = utils.WriteFile(contentpath, a.parseContent)
		if err != nil {
			record.Status = "failed"
			record.Result = err.Error()
			return record
		}
	}

	result, err := oexec(a.OriginScript, pr.PipelineId, hostPath, a.scriptPath, contentpath)
	if err != nil {
		record.Status = "failed"
		record.Result = err.Error()
		return record
	}
	// 校验结果
	for _, m := range result {
		temp := make(models.StringMap)
		for k, v := range m {
			temp[k] = v
		}
		acc.Result = append(acc.Result, temp)
	}
	pr.StepPR[step] = acc
	marshal, _ := json.Marshal(result)
	record.Result = string(marshal)
	record.Status = "finished"
	logrus.Infof("Ansible result is: %s", record.Result)
	return record
}

func findScriptDir() (path string, err error) {
	path = utils.GetCurrentDirectory() + "/script/yaml_file/"
	if utils.CheckFileIsExist(path) {
		return
	} else {
		path = "/script/yaml_file/"
		if utils.CheckFileIsExist(path) {
			return
		} else {
			return "", errors.New("script path not found")
		}
	}

}

func oexec(action string, clusterId string, hostFile string, yamlFile string, content string) (result []map[string]string, err error) {
	// 验证文件是否存在，执行
	command := "ansible-playbook -i " + hostFile + " " + yamlFile + " -v --timeout=60"
	if content != "" {
		command = "ansible-playbook -i " + hostFile + " " + yamlFile + " -e @" + content + " -v --timeout=60"
	}
	cmd := exec.Command("bash", "-c", command)
	stdout, _ := cmd.StdoutPipe() //创建输出管道
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		logrus.WithFields(logrus.Fields{
			"id":       service.Runtime.Id,
			"host":     service.Runtime.Host,
			"desc":     service.Runtime.Description,
			"script":   yamlFile,
			"content":  content,
			"hostfile": hostFile,
			"error":    err.Error(),
		}).Error("Exec command failed")
	}
	loggerForScript.WithFields(logrus.Fields{
		"workerId": service.Runtime.Id,
		"host":     service.Runtime.Host,
		"desc":     service.Runtime.Description,
		"command":  cmd.Args,
		"script":   yamlFile,
		"content":  content,
		"hostfile": hostFile,
	}).Info("Exec command beginning")

	result = make([]map[string]string, 0)
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*1000)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		logInTime := string(bytes)
		loggerForScript.WithFields(logrus.Fields{
			"clusterid": clusterId,
		}).Info(logInTime)
		if strings.HasPrefix(logInTime, "    \"msg\":") {
			scriptReg := regexp.MustCompile("^\\s*\"msg\":\\s*\"(.*?)\"$")
			varRes := scriptReg.FindStringSubmatch(logInTime)
			if len(varRes) > 1 {
				tmpMap := make(map[string]string, 0)
				kvstr := strings.Split(varRes[1], ",")
				for _, s := range kvstr {
					sub := strings.Split(s, "=")
					if len(sub) == 2 {
						tmpMap[sub[0]] = sub[1]
					}
				}
				result = append(result, tmpMap)
			}
		}
	}
	if errscan := scanner.Err(); errscan != nil {
		loggerForScript.WithFields(logrus.Fields{
			"id":        service.Runtime.Id,
			"host":      service.Runtime.Host,
			"desc":      service.Runtime.Description,
			"clusterId": clusterId,
			"script":    yamlFile,
			"content":   content,
			"hostfile":  hostFile,
			"error":     errscan.Error(),
		}).Error("Exec command scan output failed")
	}
	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			_ = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
		}
		return nil, err
	}
	return
}
