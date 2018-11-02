package exector

import (

	//log "github.com/sunminghong/goutils/logger"
)

type Params interface{}

type IHandler interface{
  Handle(params Params) Params
}

//单通道执行者
type SingleChannelExector struct {
	taskchan chan Params //每个项目一个输入参数通道,容量为1，所以同时只能执行一个bind行为，即单线程队列
	inchan  chan Params //每个项目一个输入参数通道,容量为1，所以同时只能执行一个bind行为，即单线程队列
	outchan chan Params //每个项目一个返回结果通道,容量为1，所以同时只能执行一个bind行为，即单线程队列

	handler IHandler

	quit chan bool
}

func NewSingleChannelExector(handler IHandler) *SingleChannelExector {
  sce := &SingleChannelExector{
		taskchan:  make(chan Params),
		inchan:  make(chan Params),
		outchan: make(chan Params),

		handler: handler,
		quit:    make(chan bool),
	}

	go sce._exec()
	go sce._exec_non_block()
  return sce
}

func (sce *SingleChannelExector) Stop() {
	sce.quit <- true
}

func (sce *SingleChannelExector) Call(params Params) (result Params) {

	sce.inchan <- params
  result =<- sce.outchan
	return
}

//该方法是没有返回的，也就是非阻塞的
func (sce *SingleChannelExector) AddTask(params Params){
	go func(params Params){
    sce.taskchan<- params
  }(params)
}

func (sce *SingleChannelExector) _exec_non_block() {
	for {
		select {
		case in := <-sce.taskchan:

			sce.handler.Handle(in)
			//log.Debug("SingleChannelExector._exec_non_block() return" )

		case _ = <-sce.quit:
			//log.Debug("SingleChannelExector._exec() stop")
			return
		}
	}
}

func (sce *SingleChannelExector) _exec() {
	for {
		select {
		case in := <-sce.inchan:

			out := sce.handler.Handle(in)

			sce.outchan <- out
			//log.Debug("SingleChannelExector._exec() return:%v:", out)

		case _ = <-sce.quit:
			//log.Debug("SingleChannelExector._exec() stop")
			return
		}
	}
}
