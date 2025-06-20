package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"time"

	"github.com/klauspost/compress/zstd"
)

// Zstder implements the Benchmarker interface
type Zstder struct {
}

func NewZstdRunner() *Zstder {
	return &Zstder{}
}

func (z *Zstder) RunBenchmark(input []byte) (*BenchmarkResult, error) {
	return runZstd(input, zlib.BestSpeed)
}

func runZstd(input []byte, level int) (*BenchmarkResult, error) {
	zstdBytes, zstdCompTime, err := compressZstd(input, level)
	if err != nil {
		return nil, err
	}
	zstdSize := len(zstdBytes)
	_, zstdDecompTime, err := decompressZstd(zstdBytes)
	if err != nil {
		return nil, err
	}
	zstdRatio := float64(zstdSize) / float64(len(input))
	res := BenchmarkResult{
		DecompressTime: zstdDecompTime,
		CompressTime:   zstdCompTime,
		CompressedSize: zstdSize,
		Ratio:          zstdRatio,
		Name:           "zstd",
	}
	return &res, nil
}

func compressZstd(input []byte, level int) ([]byte, time.Duration, error) {
	var buf bytes.Buffer
	t0 := time.Now()
	zw, err := zstd.NewWriter(&buf)
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

func decompressZstd(inputBytes []byte) ([]byte, time.Duration, error) {
	t0 := time.Now()
	reader := bytes.NewReader(inputBytes)
	zr, err := zstd.NewReader(reader)
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
