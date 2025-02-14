package channel

import (
	"sync"
	"testing"
)

func Test_buffer_lock(t *testing.T) {
	// 最终 panic
	bc := make(chan int)
	bc <- 1
	bc <- 1 // 死锁
	<-bc
}

func Test_buffer_lock1(t *testing.T) {
	// 最终 panic
	bc := make(chan int)
	bc <- 1
	<-bc
	<-bc // 死锁
	bc <- 1
}

func Test_buffer_panic(t *testing.T) {
	bc := make(chan int)
	close(bc)
	bc <- 1 // panic: send on closed channel
}

func Test_buffer_lock_3(t *testing.T) {
	bc := make(chan int, 2)
	bc <- 1
	bc <- 1
	<-bc
	t.Log("end0")
	<-bc
	t.Log("end1")
	<-bc // panic: test timed out after 30s
	t.Log("end2")
}

func Test_buffer_lock_4(t *testing.T) {
	bc := make(chan int, 2)
	bc <- 1
	bc <- 1
	t.Log("end0")
	bc <- 1 // panic: test timed out after 30s
	t.Log("end1")
	<-bc
	t.Log("end2")
	<-bc
	t.Log("end3")
	<-bc
	t.Log("end4")
}

func Test_buffer_lock_5(t *testing.T) {
	bc := make(chan int, 2)
	t.Log("end0")
	<-bc // panic: test timed out after 30s
	t.Log("end1")
	bc <- 1
	t.Log("end2")
}

func Test_wait_group(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	t.Log("end0") // panic: sync: negative WaitGroup counter [recovered]
	wg.Done()
	t.Log("end1") // panic: sync: negative WaitGroup counter
	wg.Done()
	t.Log("end2")
	wg.Done()
	t.Log("end3")
	wg.Wait()
	t.Log("end4")
}

func Test_cond_1(t *testing.T) {
	var mu sync.Mutex
	var cond = sync.NewCond(&mu)
	mu.Lock()
	cond.Wait() // panic: test timed out after 30s
	t.Log("end0")
	cond.Wait()
	t.Log("end1")
	cond.Wait()
	t.Log("end2")
	mu.Unlock()

	go func() {
		mu.Lock()
		cond.Signal() // signal waiting goroutine
		mu.Unlock()
	}()

}
