package scheduler

import (
	"container/list"
	"context"
	"github.com/mongodb-job/models"
	"sync"
	"time"
)

type (
	TimeWheel struct {
		randomID uint64
		interval time.Duration
		ticker   *time.Ticker

		slotNum    int
		slots      []*list.List
		currentPos int

		onceStart  sync.Once
		addTaskC   chan *models.Pipeline
		stopC      chan struct{}
		taskRecord sync.Map
	}

	TimeWheelScheduler struct { //时间轮调度
		EventsChan chan *Event         // 事件通道
		ResultChan chan *models.Result // 结果通道
		// 小时轮
		HourQueue [24]TimeWheel
		// 分钟轮
		MinuteQueue []TimeWheel
		// 秒轮
		SecondQueue []models.Pipeline
	}
)

func (tws *TimeWheelScheduler) Run(ctx context.Context) {
	//TODO 禁止直接在newtimer后获取channel
	timer := time.NewTimer(5 * time.Second)
	select {
	case e := <-tws.EventsChan:
		tws.eventHandler(e)
	case <-timer.C:
	}
	tws.TrySchedule(ctx)
	// TODO 必须在调度后重置时间, 防止memory leak
	timer.Reset(5 * time.Second)
}

func (tws *TimeWheelScheduler) DispatchEvent(event *Event) {
	tws.EventsChan <- event
}

func (tws *TimeWheelScheduler) eventHandler(event *Event) {
	switch event.Type {
	case PUT:

	case DEL:

	}
}
func (tws *TimeWheelScheduler) ResultHandler(result *models.Result) {

}

func (tws *TimeWheelScheduler) TrySchedule(ctx context.Context) {
	for ; ; {

	}
}
