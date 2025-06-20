package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"time"
)

// Gzipper implements the Benchmarker interface
type Gzipper struct {
}

func NewGzipRunner() *Gzipper {
	return &Gzipper{}
}

func (gzip *Gzipper) RunBenchmark(input []byte) (*BenchmarkResult, error) {
	return runGzip(input)
}

func runGzip(input []byte) (*BenchmarkResult, error) {
	gzipBytes, gzipCompTime, err := compressGzip(input)
	if err != nil {
		return nil, err
	}
	gzipSize := len(gzipBytes)
	_, gzipDecompTime, err := decompressGzip(gzipBytes)
	if err != nil {
		return nil, err
	}
	gzipRatio := float64(gzipSize) / float64(len(input))
	res := BenchmarkResult{
		DecompressTime: gzipDecompTime,
		CompressTime:   gzipCompTime,
		CompressedSize: gzipSize,
		Ratio:          gzipRatio,
		Name:           "gzip",
	}
	return &res, nil
}

func compressGzip(input []byte) ([]byte, time.Duration, error) {
	var buf bytes.Buffer
	t0 := time.Now()
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(input)
	if err != nil {
		return nil, 0, err
	}
	if err := zw.Close(); err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), time.Since(t0), nil
}

func decompressGzip(inputBytes []byte) ([]byte, time.Duration, error) {
	t0 := time.Now()
	reader := bytes.NewReader(inputBytes)
	zr, err := gzip.NewReader(reader)
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
