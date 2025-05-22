package utils

import (
	"fmt"
	"hash/fnv"
	"net"
	"sync"
	"time"
)

// 常量
const (
	basisTime        int64 = 1288834974657 // 起始时间点
	workerIDBits     uint  = 5             // 机器 ID 位数
	datacenterIDBits uint  = 5             // 数据中心 ID 位数
	sequenceBits     uint  = 12            // 序列号位数

	maxWorkerID     int64 = -1 ^ (-1 << workerIDBits)     // 机器 ID 最大值
	maxDatacenterID int64 = -1 ^ (-1 << datacenterIDBits) // 数据中心 ID 最大值
	sequenceMask    int64 = -1 ^ (-1 << sequenceBits)     // 序列号掩码

	workerIDShift     = sequenceBits                                   // 机器 ID 左移位数
	datacenterIDShift = sequenceBits + workerIDBits                    // 数据中心 ID 左移位数
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits // 时间戳左移位数
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex // 互斥锁，多个线程并发生成id时不发生数据竞争
	lastTimestamp int64      // 上次生成id的时间戳
	sequence      int64      // 序列号
	workerID      int64      // 机器id
	datacenterID  int64      // 数据中心id
}

// NewSnowflake 自动生成 workerID 和 datacenterID
func NewSnowflake() (*Snowflake, error) {
	workerID := generateWorkerID()
	datacenterID := generateDatacenterID()

	return &Snowflake{
		lastTimestamp: 0,
		sequence:      0,
		workerID:      workerID,
		datacenterID:  datacenterID,
	}, nil
}

// generateWorkerID 生成 workerID（基于主机名哈希）
func generateWorkerID() int64 {
	host, err := net.LookupHost("localhost") // 获取本机主机名对应的ip地址
	if err != nil || len(host) == 0 {
		return 1 // 默认值
	}
	return hashToRange(host[0], maxWorkerID) // 对ip地址做哈希运算，确保唯一性
}

// generateDatacenterID 生成 datacenterID（基于 MAC 地址哈希）
func generateDatacenterID() int64 {
	interfaces, err := net.Interfaces() // 获取所有网卡信息
	if err != nil || len(interfaces) == 0 {
		return 1 // 默认值
	}
	return hashToRange(interfaces[0].HardwareAddr.String(), maxDatacenterID)
}

// hashToRange 哈希值到固定范围
func hashToRange(input string, max int64) int64 {
	h := fnv.New32a() // 使用 FNV-1a 哈希算法
	h.Write([]byte(input))
	return int64(h.Sum32()) % (max + 1) // 确保 ID 在 0~max之间
}

// GenerateID 唯一 ID 生成
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock() // 让每次生成id都是串行执行，防止冲突
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli()

	// 处理时钟回拨
	offset := s.lastTimestamp - timestamp
	if offset > 0 {
		if offset <= 5 { // 若回拨小于等于 5ms，等待时钟追上
			time.Sleep(time.Duration(offset) * time.Millisecond)
			timestamp = time.Now().UnixMilli()
		} else {
			// 若回拨超过 5ms，则直接 panic
			panic(fmt.Sprintf("时钟回拨，拒绝生成 ID: %d 毫秒", offset))
		}
	}

	// 处理序列号
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			timestamp = waitNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequence = 1
	}

	s.lastTimestamp = timestamp

	return ((timestamp - basisTime) << timestampShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence
}

// waitNextMillis 等待下一毫秒
func waitNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixMilli()
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixMilli()
	}
	return timestamp
}
