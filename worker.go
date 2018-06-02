package goes

//
// Author: 陈永佳 chenyongjia@parkingwang.com, yoojiachen@gmail.com
//

type GoTask func() // 任务

// Worker包含一个任务列表和停止控制信号
type worker struct {
	tasks   chan GoTask   // 任务列表
	onIdle  func(*worker) // 完成任务后空闲通知
	stop    chan struct{} // 停止信号
	running chan struct{} // 运行状态
}

// 启动内部协程，异步循环处理任务列表
func (slf *worker) Run() {
	if nil == slf.onIdle {
		panic("Func < onIdle(*worker) > for worker is nil")
	}
	go func() {
		for {
			select {
			case <-slf.stop:
				close(slf.running)
				return

			case fun := <-slf.tasks:
				func() {
					defer slf.onIdle(slf)
					fun()
				}()
			}
		}
	}()
}

func newWorker() *worker {
	return &worker{
		tasks:   make(chan GoTask, 1), // 1: for sync send task
		stop:    make(chan struct{}),
		running: make(chan struct{}),
	}
}
