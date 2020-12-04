package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type (
	Etcd struct {
		Killer    string   `json:"killer" yaml:"killer" validate:"required"`
		Locker    string   `json:"locker" yaml:"locker" validate:"required"`
		Service   string   `json:"service" yaml:"service" validate:"required"`     //worker节点注册
		Pipeline  string   `json:"pipeline" yaml:"pipeline" validate:"required"`   //任务分配到哪个节点调度数据存储
		Snapshot  string   `json:"snapshot" yaml:"snapshot" validate:"required"`   //任务分配到哪个节点调度数据存储
		Plan      string   `json:"plan" yaml:"plan" validate:"required"`           //任务元数据存储
		Foreign   string   `json:"foreign" yaml:"foreign" validate:"required"`     //任务元数据存储
		EndPoints []string `json:"endpoints" yaml:"endpoints" validate:"required"` // etcd
		Timeout   int64    `json:"timeout" yaml:"timeout" validate:"required"`
	}
	Database struct {
		Host string `json:"host" yaml:"host" validate:"required"`
		Port int    `json:"port" yaml:"port" validate:"required"`
		Name string `json:"name" yaml:"name" validate:"required"`
		User string `json:"user" yaml:"user" validate:"required"`
		Pass string `json:"pass" yaml:"pass" validate:"required"`
		Char string `json:"char" yaml:"char" validate:"required"`
	}

	App struct {
		Service string `json:"service" yaml:"service"`
	}

	Api struct {
		GicNewMonitorSystemUrl string `json:"gic_new_monitor_system_url" yaml:"gicNewMonitorSystemUrl"`
		GicUserUrl             string `json:"gic_user_url" yaml:"gicUserUrl"`
	}

	Config struct {
		Database `json:"database" yaml:"database"`
		Etcd     `json:"etcd" yaml:"etcd"`
		App      `json:"app" yaml:"app"`
		Api      `json:"api" yaml:"api"`
	}
)

var (
	Conf *Config
	Path string
)

func Init(configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// 检查配置文件是否存在
func CheckConfigFile(path string) (bool, error) {
	_, err := os.Stat(path)
	exist := !os.IsNotExist(err)
	return exist, err
}

// 创建配置文件目录
func CreateConfigDir(dir string) {
	_, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

// 检查配置文件目录是否有权限
func CheckConfigDirPermisson(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		log.Println(err)
	}
	mode := info.Mode()
	perm := mode.Perm()
	flag := perm & os.FileMode(493)
	if flag == 493 {
		return true
	}

	return false
}

// 写入配置文件
func WriteConfigToFile(file string, content []byte) bool {
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		log.Println(os.IsNotExist(err))
		log.Println(err)
	}
	return true
}
