package plugins

import (
	"context"
	"github.com/mongodb-job/models"
	"github.com/prometheus/common/log"
	"os/exec"
)

type (
	Shell struct {
		User    string
		Env     []string
		Dir     string
		Command string
	}
)

func (s *Shell) Exec(ctx context.Context, step int, pr models.PipelineRuntimeData) *models.TaskRecords {
	s.Command = Render(s.Command, pr)
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", s.Command)
	acc := models.ParamResult{Param: models.StringMap{"content": s.Command}}
	record := &models.TaskRecords{}
	if s.User != "" {
		//credential, err := getCredential(s.User)
		//if err != nil {
		//	record.Status = "failed"
		//	record.Result = err.Error()
		//	return record
		//}
		//cmd.SysProcAttr.Credential = credential
	}

	resChan := make(chan struct {
		output []byte
		err    error
	})

	go func() {
		output, err := cmd.CombinedOutput()
		resChan <- struct {
			output []byte
			err    error
		}{output: output, err: err}
	}()
	res := <-resChan
	record.Result = string(res.output)
	if res.err != nil {
		log.Info(res.err.Error())
		record.Status = "failed"
	} else {
		acc.Result = []models.StringMap{{"result": record.Result}}
		pr.StepPR[step] = acc
		record.Status = "finished"
	}
	return record
}

// 获取执行证书
//func getCredential(username string) (*syscall.Credential, error) {
//	sysuser, err := user.Lookup(username)
//	if err != nil {
//		return nil, err
//	}
//	uid, err := strconv.Atoi(sysuser.Uid)
//	gid, err := strconv.Atoi(sysuser.Gid)
//	return &syscall.Credential{
//		Uid:         uint32(uid),
//		Gid:         uint32(gid),
//		Groups:      nil,
//		NoSetGroups: false,
//	}, nil
//}
