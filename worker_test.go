package goes

import (
	"testing"
	"time"
)

//
// Author: 陈永佳 chenyongjia@parkingwang.com, yoojiachen@gmail.com
//

func TestWorker_Run(t *testing.T) {
	w := newWorker()
	defer close(w.stop)

	isNotify := false
	w.idleNotify = func(w *worker) {
		isNotify = true
	}
	w.Run()

	done := make(chan bool)
	w.taskQueue <- func() {
		done <- true
	}

	time.AfterFunc(time.Millisecond, func() {
		t.Error("Timeout")
	})

	<-done

	if !isNotify {
		t.Error("Notify not call")
	}else{
		t.Log("Test ok")
	}
}

func BenchmarkWorker(b *testing.B) {
	w := newWorker()
	defer close(w.stop)
	w.idleNotify = func(w *worker) {}
	w.Run()

	count := 0
	for i := 0; i < b.N; i++ {
		w.taskQueue <- func() {
			count ++
		}
	}
}