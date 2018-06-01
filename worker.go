package goes

//
// Author: 陈永佳 chenyongjia@parkingwang.com, yoojiachen@gmail.com
//

type GoTask func() // 任务

// Worker包含一个任务列表和停止控制信号
type worker struct {
	taskQueue  chan GoTask          // 任务列表
	idleNotify func(worker *worker) // 完成任务后空闲通知
	stop       chan struct{}        // 停止信号
	state      chan struct{}        // 运行状态
}

// 启动内部协程，异步循环处理任务列表
func (slf *worker) Run() {
	if nil == slf.idleNotify {
		panic("Func < idleNotify(*worker) > for worker is nil")
	}
	go func() {
		for {
			select {
			case <-slf.stop:
				close(slf.state)
				return

			case task := <-slf.taskQueue:
				func(){
					defer slf.idleNotify(slf)
					task()
				}()
			}
		}
	}()
}

func newWorker() *worker {
	return &worker{
		taskQueue: make(chan GoTask, 1), // 1: for sync send task
		stop:      make(chan struct{}),
		state:     make(chan struct{}),
	}
}
