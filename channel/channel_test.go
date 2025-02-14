package channel

import (
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
