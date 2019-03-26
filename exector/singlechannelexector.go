package exector

import (
  "sync"

log "github.com/sunminghong/goutils/logger"
)

type Params interface{}

type IHandler interface {
	Handle(params Params) Params
}

//单通道执行者
type SingleChannelExector struct {
	taskchan chan Params //每个项目一个输入参数通道,容量为1，所以同时只能执行一个bind行为，即单线程队列
	inchan   chan Params //每个项目一个输入参数通道,容量为1，所以同时只能执行一个bind行为，即单线程队列
	outchan  chan Params //每个项目一个返回结果通道,容量为1，所以同时只能执行一个bind行为，即单线程队列

	handler IHandler

	stop bool
  wg *sync.WaitGroup
}

func NewSingleChannelExector(handler IHandler) *SingleChannelExector {
	sce := &SingleChannelExector{
		taskchan: make(chan Params),
		inchan:   make(chan Params),
		outchan:  make(chan Params),

    wg: &sync.WaitGroup{},

		stop: false,
		handler: handler,
	}

	go sce._exec()
	go sce._exec_non_block()

	return sce
}

func (sce *SingleChannelExector) Quit() {
	sce.stop = true

  log.Info("single task start stop/////////////////////////////////")

	close(sce.taskchan)
	close(sce.inchan)
	close(sce.outchan)

  sce.wg.Wait()
}

func (sce *SingleChannelExector) Call(params Params) (result Params) {
	if sce.stop {
    return
	}
	sce.inchan <- params
	result = <-sce.outchan
	return
}

//该方法是没有返回的，也就是非阻塞的
func (sce *SingleChannelExector) AddTask(params Params) {
	if sce.stop {
    return
	}
	go func(params Params) {
		sce.taskchan <- params
	}(params)
}

func (sce *SingleChannelExector) _exec_non_block() {
  sce.wg.Add(1)

  defer sce.wg.Done()

	for {
		select {
		case in, ok := <-sce.taskchan:
			if ok {

				sce.handler.Handle(in)
				//log.Debug("SingleChannelExector._exec_non_block() return" )
			} else {
        log.Info("_exec_non_block task stoped/////////////////////////////////")
				return
			}
		}
	}
}

func (sce *SingleChannelExector) _exec() {
  sce.wg.Add(1)
  defer sce.wg.Done()

	for {
		select {
		case in, ok := <-sce.inchan:
			if ok {

				out := sce.handler.Handle(in)

				sce.outchan <- out
				//log.Debug("SingleChannelExector._exec() return:%v:", out)
			} else {
        log.Info("_exec_block task stoped/////////////////////////////////")
				return
			}
		}
	}
}
