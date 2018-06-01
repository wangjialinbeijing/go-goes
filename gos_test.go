package gos

import (
	"testing"
	"strconv"
	"sync"
	"time"
	"runtime"
)

//
// Author: 陈永佳 chenyongjia@parkingwang.com, yoojiachen@gmail.com
// Go Pool Test
//

func init() {
	println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestAddTask(t *testing.T) {
	god := NewGos(100, 10)
	god.Start()
	defer god.Shutdown()

	done := make(chan bool)
	scheduled := false

	god.Add(func() {
		println("Called")
		scheduled = true
		close(done)
	})

	<-done
	if !scheduled {
		t.Error("Task is not be scheduled")
	}
}

func TestTask2(t *testing.T) {
	dispatcher := NewGos(100, 10)
	dispatcher.Start()
	defer dispatcher.Shutdown()

	wg := new(sync.WaitGroup)

	TASKS := int(1000*100)
	for i := 0; i < TASKS; i++ {
		wg.Add(1)
		dispatcher.Add(func() {
			wg.Done()
		})
	}

	timer := time.AfterFunc(time.Millisecond*100, func() {
		t.Error("Timeout 100ms")
	})

	wg.Wait()
	timer.Stop()

	t.Log("Send tasks:" + strconv.Itoa(TASKS))
}

func BenchmarkCpuNumWorkers(b *testing.B) {
	MakeBenchmarkWith(runtime.NumCPU(), b)
}

func Benchmark1KWorkers(b *testing.B) {
	MakeBenchmarkWith(1000, b)
}

func MakeBenchmarkWith(numWorkers int,b *testing.B) {
	dispatcher := NewGos(numWorkers, numWorkers)
	dispatcher.Start()
	defer dispatcher.Shutdown()

	wg := new(sync.WaitGroup)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		dispatcher.Add(func() {
			wg.Done()
		})
	}

	timer := time.AfterFunc(time.Millisecond*100, func() {
		b.Error("Timeout 100ms")
	})

	wg.Wait()
	timer.Stop()
}