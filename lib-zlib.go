package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"time"
)

// Zlibber implements the Benchmarker interface
type Zlibber struct {
	level int
}

func NewZlibRunner(level int) *Zlibber {
	return &Zlibber{
		level: level,
	}
}

func (z *Zlibber) RunBenchmark(input []byte) (*BenchmarkResult, error) {
	return runZlib(input, z.level)
}

func runZlib(input []byte, level int) (*BenchmarkResult, error) {
	zlibBytes, zlibCompTime, err := compressZlib(input, level)
	if err != nil {
		return nil, err
	}
	zlibSize := len(zlibBytes)
	_, zlibDecompTime, err := decompressZlib(zlibBytes)
	if err != nil {
		return nil, err
	}
	zlibRatio := float64(zlibSize) / float64(len(input))
	var runnerName string
	switch level {
	case zlib.DefaultCompression:
		runnerName = "zlib-default"
	case zlib.BestCompression:
		runnerName = "zlib-best-compression"
	case zlib.BestSpeed:
		runnerName = "zlib-best-speed"
	default:
		runnerName = fmt.Sprintf("zlib-%d", level)
	}
	res := BenchmarkResult{
		DecompressTime: zlibDecompTime,
		CompressTime:   zlibCompTime,
		CompressedSize: zlibSize,
		Ratio:          zlibRatio,
		Name:           runnerName,
	}
	return &res, nil
}

func compressZlib(input []byte, level int) ([]byte, time.Duration, error) {
	var buf bytes.Buffer
	t0 := time.Now()
	zw, err := zlib.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, 0, err
	}
	_, err = zw.Write(input)
	if err != nil {
		return nil, 0, err
	}
	if err := zw.Close(); err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), time.Since(t0), nil
}

func decompressZlib(inputBytes []byte) ([]byte, time.Duration, error) {
	t0 := time.Now()
	reader := bytes.NewReader(inputBytes)
	zr, err := zlib.NewReader(reader)
	if err != nil {
		return nil, 0, err
	}
	defer zr.Close()
	outBytes, err := io.ReadAll(zr)
	if err != nil {
		return nil, 0, err
	}
	return outBytes, time.Since(t0), nil
}
