package discovery

import (
	"log"
	"testing"
)

func TestService_WatchWorker(t *testing.T) {
	initClient()
	log.Println("初始化连接成功")
	WatchWorker()
}
