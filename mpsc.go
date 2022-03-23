package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type TaskNode struct {
	Data interface{} `json:"data"`
	Next *TaskNode   `json:"Next"`
}

var UNCONNECTED *TaskNode = new(TaskNode)

func NewExecutionQueue(_func func(interface{})) *ExecutionQueue {
	return &ExecutionQueue{
		Head:          nil,
		_execute_func: _func,
		locker:        sync.Mutex{},
		pool: &sync.Pool{New: func() interface{} {
			return new(TaskNode)
		}},
	}
}

type ExecutionQueue struct {
	Head          *TaskNode         `json:"Head"`
	_execute_func func(interface{}) `json:"-"` // 消费者函数
	locker        sync.Mutex        `json:"-"`
	pool          *sync.Pool        `json:"-"`
}

func (ex *ExecutionQueue) AddTaskNode(data interface{}) {
	node := ex.pool.Get().(*TaskNode)
	node.Data = data
	node.Next = UNCONNECTED

	preHead := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ex.Head)), unsafe.Pointer(node))

	if preHead != nil {
		node.Next = (*TaskNode)(preHead)
		return
	}

	node.Next = nil
	// 任务不多直接执行，防止线程切换
	ex._execute_func(node.Data)
	if !ex.moreTasks(node) {
		return
	}
	go ex.consumeTasks(node)

}

func (ex *ExecutionQueue) moreTasks(oldNode *TaskNode) bool {

	newHead := oldNode

	if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ex.Head)), unsafe.Pointer(newHead), nil) {
		return false
	}
	newHead = (*TaskNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ex.Head))))
	var tail *TaskNode
	p := newHead
	for {
		for {
			if p.Next != UNCONNECTED {
				break
			} else {
				runtime.Gosched()
			}
		}
		saved_next := p.Next
		p.Next = tail
		tail = p
		p = saved_next

		if p == oldNode {
			oldNode.Next = tail
			return true
		}
	}
}

func (ex *ExecutionQueue) consumeTasks(taskNode *TaskNode) {
	for {
		tmp := taskNode

		taskNode = taskNode.Next
		tmp.Next = nil
		ex.pool.Put(tmp)
		ex._execute_func(taskNode.Data)

		if taskNode.Next == nil && !ex.moreTasks(taskNode) {
			return
		}
	}
}

var count int64 = 0

func print(data interface{}) {
	atomic.AddInt64(&count, 1)
}

func producer() {
	var singalexit = sync.WaitGroup{}
	ex := NewExecutionQueue(print)
	for i := 0; i < 10000; i++ {
		singalexit.Add(1)
		go func(i int, singalexit *sync.WaitGroup) {
			defer singalexit.Done()
			for j := 0; j < 1000; j++ {
				ex.AddTaskNode(i*100 + j)
			}
		}(i, &singalexit)
	}

	singalexit.Wait()
	time.Sleep(2 * time.Second)
	fmt.Println(atomic.LoadInt64(&count))
}

func main() {
	for i := 0; i < 10; i++ {
		count = 0
		producer()
	}
}
