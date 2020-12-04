package cmd

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mongodb-job/config"
	"github.com/mongodb-job/handler"
	"github.com/mongodb-job/internal/discovery"
	"github.com/mongodb-job/internal/etcd"
	"github.com/mongodb-job/internal/service"
	"github.com/mongodb-job/models"
	task "github.com/mongodb-job/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	masterConfigPath = ""

	masterCmd = &cobra.Command{
		Use:   "master",
		Short: "Run a master node service.",
		Long:  "Run a master node service on this server.",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}
	master = &service.Instance{
		Mode:        service.MASTER,
		Status:      service.ONLINE,
		Version:     rootCmd.Version,
		Description: "mongodb-job master",
	}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.AddCommand(masterCmd)
	masterCmd.Flags().StringVar(&masterConfigPath, "conf", "", "Get config path")
	masterCmd.Flags().StringVar(&master.Host, "serviceName", "", "ServiceName")
	master.Id = os.Getenv("MJOBUNIQUEKEY")
}

func start() {
	var err error
	logPath := "/app/logs/info.log"
	if writer, err := rotatelogs.New(
		logPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(360)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	);err != nil {
		panic(err)
	} else {
		logrus.SetOutput(writer)
	}

	// 初始化服务本身信息
	if master.Id == "" {
		log.Fatal("Master id not found")
	}

	if master.Host == ""{
		log.Fatal("Master host not found")
	}

	if master.Name == "" {
		master.Name, err = os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
	}
	service.Runtime = master

	// 处理启动前的依赖
	config.Init(masterConfigPath)
	models.Connection()
	etcd.NewClient(config.Conf.Etcd)

	registry, err := discovery.NewService(master)
	if err != nil {
		logrus.Fatal(err)
	}

	go func(ser *discovery.Service) {
		if err := ser.Register(5); err != nil {
			logrus.Fatal(err)
		}
	}(registry)

	_, cancelFunc := context.WithCancel(context.Background())


	// 监听pipeline
	go discovery.WatchWorker()
	// 监听worker
	go discovery.WatchPipelines()
	// 启动server
	go handler.Start()

	go http.ListenAndServe("0.0.0.0:6060", nil)

	go task.AutoDelete()
	go task.AutoDeleteIncrBackup()
	go task.AutoDeleteTempCluster()
	// 放推出
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
