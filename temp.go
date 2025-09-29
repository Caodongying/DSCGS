package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// id: 1(不用) + 41(timestamp，毫秒，69年) + 10(机器号) + 12(序列号，毫秒内，4096个)
// 当新的毫秒到来，需要将序列号置0吗，还是继续分发序列号，直到达到4096、需要等待下一个序列号

const (
	// TODO: 改成12
	sequenceBits int64 = 3 // 应该是uint8
	workerIDBits int64 = 10
	maxSequenceNumber int64 = 1 << sequenceBits - 1// 一毫秒内一共可以分配的序列号数量
	// 下面不对，因为^在Go语言里不是幂运算符，而是按位异或运算符，同0异1，结果是14
	// maxSequenceNumber int64 = 2 ^ sequenceBits // 一毫秒内一共可以分配的序列号数量

	// 二进制里，表示正数的负数，是取反再加1，比如111111（若干1）就是-1
	maxWorker int64 = 1 << workerIDBits // 应该是 -1 ^ (-1<<workerIDBits)

	workerIDShift int64 = sequenceBits
	timestampShift int64 = workerIDShift + workerIDBits
)

type Worker struct {
	mu sync.Mutex
	workerID int64
	lastAssignTimeStamp int64 // 上一次分发ID时的毫秒级时间戳
	count int64 // 在当前毫秒内已经分发的ID数
}

func NewWorker(workerID int64) (*Worker, error) { // int和int64区别在哪
	if workerID >= maxWorker {
		return nil, errors.New("超出最大worker限制")
	}
	if workerID < 0 {
		return nil, errors.New("workerID不可以小于0")
	}
	return &Worker{
		// mu: sync.Mutex{}, 为什么不用初始化
		workerID: workerID,
		lastAssignTimeStamp: 0,
		count: 0,
	}, nil

}

func (worker *Worker) GetID() int64 {
	worker.mu.Lock()
	defer worker.mu.Unlock()

	currentTimeStamp := int64(time.Now().Unix())
	println("Current timestamp is", currentTimeStamp)
	if currentTimeStamp == worker.lastAssignTimeStamp {
		// 和上次分配在同一毫秒内发生
		if worker.count < maxSequenceNumber {
			// 当前毫秒内序列号没有用完
			worker.count++
		} else {
			// 当前毫秒内序列号已经用完,等待下一个毫秒到来
			for currentTimeStamp !=  worker.lastAssignTimeStamp + 1 {
				currentTimeStamp = int64(time.Now().Unix())
			}
			worker.count = 0
		}
	} else {
		// 在新的毫秒内发生
		worker.count = 0
	}

	// 开始拼接ID
	id := currentTimeStamp << timestampShift | worker.workerID << workerIDShift | worker.count // currentTimeStamp-startTime，为什么

	worker.lastAssignTimeStamp = currentTimeStamp

	return id
}

func main() {
	worker, err := NewWorker(1)
	if err != nil {
		panic(err)
	}

	fmt.Println("maxSequenceNumber is ", maxSequenceNumber)
	fmt.Println("maxWorker is ", maxWorker)
	fmt.Println("timestampShift is ", timestampShift)

	for i:=0; i<100; i++{
		fmt.Printf("%019b\n", worker.GetID()) // worker不用解引用吗
	}

}