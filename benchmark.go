package main

import "time"

type Benchmarker interface {
	RunBenchmark([]byte) (*BenchmarkResult, error)
}

type BenchmarkResult struct {
	Name           string
	CompressTime   time.Duration
	DecompressTime time.Duration
	CompressedSize int
	Ratio          float64
}

func (br *BenchmarkResult) GetTotalTime() time.Duration {
	return br.CompressTime + br.DecompressTime
}

func (br *BenchmarkResult) GetBatchTime(n int, speed uint64) time.Duration {
	if speed == 0 {
		return 0
	}
	networkTime := time.Duration((float64(br.CompressedSize) / float64(speed)))
	ops := []time.Duration{br.CompressTime, networkTime, br.DecompressTime}
	slowest := 0
	for i := 1; i < len(ops); i++ {
		if ops[i] > ops[slowest] {
			slowest = i
		}
	}
	var total time.Duration
	// total duration of compressing, sending, and decompressing n payloads
	// can be estimated by (slowest operation * n) + sum of all other operations
	for i, optime := range ops {
		if i == slowest {
			total += time.Duration(n) * optime
		} else {
			total += optime
		}
	}
	return total
}
