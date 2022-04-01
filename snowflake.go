package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10                      // 节点位数
	numberBits  uint8 = 12                      // 每个节点1毫秒内可生成的id序号的二进制位数
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点id最大值
	numberMax   int64 = -1 ^ (-1 << numberBits) // id最大值
	timeShift   = workerBits + numberBits       // 时间戳向左的偏移量
    workerShift = numberBits                    // 节点id向左的偏移量
    epoch       int64 = 1574915060000           // 单位毫秒
)

type SequenceResolver func(w *Worker, ms int64) (uint16, error)

type ID uint64

type Worker struct {
	mu sync.Mutex

	timestamp int64
	workerId  int64
	number    int64
	epoch     time.Time
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID out of range")
	}
	var curTime = time.Now()
	return &Worker{
		timestamp: 0,
		workerId: workerId,
		number: 0,
		epoch: curTime.Add(time.Unix(epoch / 1000, (epoch % 1000) * 1000000).Sub(curTime)),
	}, nil
}

func (w *Worker) Generate() ID {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Since(w.epoch).Nanoseconds() / 1e6
	if now == w.timestamp {
		w.number = (w.number + 1) & numberMax
		if w.number == 0 {
			for now <= w.timestamp {
				now = time.Since(w.epoch).Nanoseconds() / 1e6
			}
		}
	} else {
		w.number = 0
	}
	w.timestamp = now

	id := ID((now << timeShift) | (w.workerId << workerShift) | w.number)

	return id
}