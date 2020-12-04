package plugins

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"time"
)

type (

	Contract interface {
		Exec(ctx context.Context) *models.TaskRecords
	}
)

var (
	loggerForScript *logrus.Logger
)

func init() {
	loggerForScript = logrus.New()
	logPath := "/app/logs/script.log"
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(240)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err == nil {
		loggerForScript.SetOutput(writer)
	} else {
		panic(err)
	}
}


