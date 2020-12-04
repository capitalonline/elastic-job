package service

const (
	ONLINE  = "up"
	OFFLINE = "offline"
	MASTER  = "master"
	WORKER  = "worker"
)
type (
	Instance struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Mode        string `json:"mode"`
		Status      string `json:"status"`
		Version     string `json:"version"`
		Description string `json:"description"`
	}
)

var (
	ConfigKey string    // TODO 动态配置文件的Key
	Runtime   *Instance // 运行服务的星系
	EndPoints []string  // ETCD 节点信息
)
