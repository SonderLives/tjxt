package idgen

import (
	"sync"
	"time"
)

// snowflake 参数
const (
	epoch        = 1609459200000 // 2021-01-01 00:00:00 UTC 起算
	machineBits  = 10
	sequenceBits = 12
	maxSequence  = 1<<sequenceBits - 1
	maxMachine   = 1<<machineBits - 1
	machineShift = sequenceBits
	timeShift    = machineBits + sequenceBits
)

var (
	machineID int64 = 1 // 单机部署固定为 1，多实例需按实例分配
	mu        sync.Mutex
	sequence  int64
	lastTime  int64
)

// NextID 返回下一个唯一雪花 ID。
func NextID() int64 {
	now := time.Now().UnixMilli()

	mu.Lock()
	defer mu.Unlock()

	if now < lastTime {
		// 时钟回拨，退回到上一次时间，保证仍单调
		now = lastTime
	}
	if now == lastTime {
		sequence = (sequence + 1) & maxSequence
		// 同一毫秒内序号耗尽，等待下一毫秒
		if sequence == 0 {
			for now <= lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sequence = 0
	}
	lastTime = now

	id := (now-epoch)<<timeShift | machineID<<machineShift | sequence
	return id
}
