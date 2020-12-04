package cmd

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mongodb-job/config"
	"github.com/mongodb-job/handler"
	"github.com/mongodb-job/internal/discovery"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/scheduler"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	workerConfigPath = ""
	workerCmd = &cobra.Command{
		Use:   "worker",
		Short: "Run a worker node service",
		Long:  "Run a worker node service on this server",
		Run: func(cmd *cobra.Command, args []string) {
			listen()
		},
	}

	worker = &service.Instance{
		Mode:        service.WORKER,
		Status:      service.ONLINE,
		Version:     rootCmd.Version,
		Description: "mongodb-job worker",
	}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringVar(&workerConfigPath, "conf", "./mongodb-job.yaml", "Config")
	workerCmd.Flags().StringVar(&worker.Host, "serviceName", "", "ServiceName")
	worker.Id = os.Getenv("MJOBUNIQUEKEY")

}

func listen() {
	var err error
	logPath := "/app/logs/info.log"
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(240)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err == nil {
		logrus.SetOutput(writer)
	} else {
		panic(err)
	}

	// 初始化服务本身信息
	if worker.Id == "" {
		log.Fatal("worker id not found")
	}

	if worker.Host == ""{
		log.Fatal("worker host not found")
	}

	if worker.Name == "" {
		worker.Name, err = os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
	}

	service.Runtime = worker
	logrus.WithFields(logrus.Fields{
		"mode": service.Runtime.Mode,
		"id": service.Runtime.Id,
		"desc": service.Runtime.Description,
		"host": service.Runtime.Host,
	}).Info("MongoDB-Job Worker Started")
	// 处理启动前的依赖
	config.Init(workerConfigPath)
	models.Connection()
	etcd.NewClient(config.Conf.Etcd)

	registry, err := discovery.NewService(worker)
	if err != nil {
		logrus.Fatal(err)
	}

	go func(ser *discovery.Service) {
		if err := ser.Register(5); err != nil {
			logrus.Fatal(err)
		}
	}(registry)

	scheduler.InitScheduler()

	ctx, cancelFunc := context.WithCancel(context.Background())

	go scheduler.Instance.Run(ctx)
	go discovery.WatchPlan(worker)
	go handler.StartWorker()

	sign := make(chan os.Signal, 1)

	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		receiver := <-sign
		log.Printf("get a signal %s", receiver.String())
		switch receiver {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			registry.Stop()
			cancelFunc()
			return
		}
	}
}
