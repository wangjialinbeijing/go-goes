package gos

//
// Author: 陈永佳 yoojiachen@gmail.com
//

// Dispatcher内部维护一组Worker。
// 如果接收到外部任务，向空闲Worker发送任务
type Gos struct {
	workers    []*worker     // 所有Worker
	ready      chan *worker  // 空闲
	tasks      chan GoTask   // 待调度任务列表
	stop       chan struct{} // 停止信号
	terminated chan struct{} // 终止状态
}

func (slf *Gos) Start() {
	go func() {
		for {
			select {
			case <-slf.stop:
				//// 关闭所有Worker
				for _, worker := range slf.workers {
					close(worker.stop)
					<-worker.state
				}
				close(slf.terminated)
				return

			case task := <-slf.tasks:
				// 取空闲Worker，发任务给它处理
				worker := <-slf.ready
				worker.taskQueue <- task
			}
		}
	}()
}

// 关闭Go调度器，等待所有Go协程完成后返回
func (slf *Gos) Shutdown() {
	close(slf.stop)
	// 等待所有任务完成
	<-slf.terminated
}

// 添加需要调度的任务
func (slf *Gos) Add(task GoTask) {
	slf.tasks <- task
}

func NewGos(numWorkers int, taskQueueSize int) *Gos {
	numWorkers = max(1, numWorkers)
	taskQueueSize = max(1, taskQueueSize)
	d := &Gos{
		workers:    make([]*worker, numWorkers),
		ready:      make(chan *worker, numWorkers),
		tasks:      make(chan GoTask, taskQueueSize),
		stop:       make(chan struct{}),
		terminated: make(chan struct{}),
	}

	//// 初始化Worker列表
	notify := func(worker *worker) {
		d.ready <- worker
	}
	//
	for i := 0; i < numWorkers; i++ {
		worker := newWorker()
		worker.idleNotify = notify
		worker.Run()
		d.workers[i] = worker
		d.ready <- worker
	}

	return d
}

func max(a int, b int) int {
	if a > b {
		return a
	}else{
		return b
	}
}
