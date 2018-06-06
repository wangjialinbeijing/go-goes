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
	defer close(w.tasks)

	isNotify := false
	w.onIdle = func(w *worker) {
		isNotify = true
	}
	go w.startWork()

	done := make(chan bool)
	w.tasks <- func() {
		done <- true
	}

	time.AfterFunc(time.Millisecond, func() {
		t.Error("Timeout")
	})

	<-done

	if !isNotify {
		t.Error("Notify not call")
	} else {
		t.Log("Test ok")
	}
}

func BenchmarkWorker(b *testing.B) {
	w := newWorker()
	defer close(w.tasks)

	w.onIdle = func(w *worker) {}
	go w.startWork()

	count := 0
	for i := 0; i < b.N; i++ {
		w.tasks <- func() {
			count++
		}
	}
}
