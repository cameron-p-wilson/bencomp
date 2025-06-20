package main

import (
	"testing"
	"time"
)

func TestBatchTime(t *testing.T) {
	type testData struct {
		name  string
		size  int
		cTime int64
		dTime int64
		n     int
		speed uint64
		exp   int64
	}
	tests := []testData{
		{
			name:  "all times are equal",
			size:  10000,
			cTime: 100,
			dTime: 100,
			n:     100,
			speed: 100,
			exp:   100 + 100 + (100 * 100), // 10200
		},
		{
			name:  "compression is slowest",
			size:  10000,
			cTime: 250,
			dTime: 100,
			n:     100,
			speed: 100,
			exp:   100 + 100 + (250 * 100), // 25200
		},
		{
			name:  "network is slowest",
			size:  10000,
			cTime: 100,
			dTime: 100,
			n:     100,
			speed: 10,
			exp:   100 + 100 + (1000 * 100), // 100200
		},
		{
			name:  "decompression is slowest",
			size:  10000,
			cTime: 100,
			dTime: 400,
			n:     100,
			speed: 100,
			exp:   100 + 100 + (400 * 100), // 40200
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			br := BenchmarkResult{
				CompressTime:   time.Duration(test.cTime),
				DecompressTime: time.Duration(test.dTime),
				CompressedSize: test.size,
			}
			out := br.GetBatchTime(test.n, test.speed)
			if out != time.Duration(test.exp) {
				t.Errorf("expected %d but got %d", time.Duration(test.exp), out)
			}
		})
	}
}
