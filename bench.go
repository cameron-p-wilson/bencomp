package main

import (
	"compress/zlib"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	newJsonGenerator = NewJsonGenerator
)

type PrintOptions struct {
	ShouldPrintInput bool
	ShouldPrintCTime bool
	ShouldPrintDTime bool
	NetworkSpeed     uint64
	NetworkPayloads  int
}

func NewBenchCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "",
		Short: "Utility for comparing compression performance",
		Long: `Runs several different compression algorithms then
reports performance statistics.

You can choose to run the benchmark against a specific file or randomly
generate a JSON input with a variey of parameters. The benchmark runs gzip,
3 different levels of zlib, and zstd.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runBenchmark(cmd)
			if err != nil {
				return err
			}
			return nil
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	createBenchFlags(root)
	return root
}

// Performs the benchmark according to user flags
func runBenchmark(cmd *cobra.Command) error {
	count, _ := getCountFlag(cmd)
	printOptions, err := getPrintOptions(cmd)
	if err != nil {
		return fmt.Errorf("error while preparing benchmark: %v", err)
	}
	nResults := make([][]*BenchmarkResult, 0, count)
	var input []byte
	for _ = range count {
		input, err = getBenchmarkInput(cmd)
		if err != nil {
			return fmt.Errorf("error while preparing benchmark: %v", err)
		}
		if len(input) == 0 {
			return fmt.Errorf("error while preparing benchmark: there is nothing to compress")
		}
		benchmarkers := getBenchmarkers()
		results := make([]*BenchmarkResult, len(benchmarkers))
		for i, benchmarker := range benchmarkers {
			result, err := benchmarker.RunBenchmark(input)
			if err != nil {
				return fmt.Errorf("error while running benchmark: %v", err)
			}
			results[i] = result
		}
		nResults = append(nResults, results)
	}
	aggResults := aggregateResults(nResults)
	printResults(input, printOptions, aggResults)
	return nil
}

func getBenchmarkers() []Benchmarker {
	return []Benchmarker{
		NewGzipRunner(),
		NewZlibRunner(zlib.DefaultCompression),
		NewZlibRunner(zlib.BestCompression),
		NewZlibRunner(zlib.BestSpeed),
		NewZstdRunner(),
	}
}

func aggregateResults(nResults [][]*BenchmarkResult) []*BenchmarkResult {
	n := len(nResults) // number of test runs
	if n == 0 {
		return nil
	}
	m := len(nResults[0]) // number of compression libraries tested
	final := make([]*BenchmarkResult, 0, m)
	for libraryID := range m {
		// for each compression library
		libName := nResults[0][libraryID].Name
		compTimes := make([]time.Duration, n)
		decompTimes := make([]time.Duration, n)
		sizes := make([]int, n)
		ratios := make([]float64, n)
		for resultIndex, results := range nResults {
			result := results[libraryID]
			compTimes[resultIndex] = result.CompressTime
			decompTimes[resultIndex] = result.DecompressTime
			sizes[resultIndex] = result.CompressedSize
			ratios[resultIndex] = result.Ratio
		}
		medCompTime := findMedianTime(compTimes)
		medDecompTime := findMedianTime(decompTimes)
		medSize := findMedianInt(sizes)
		medRatio := findMedianFloat64(ratios)
		libraryAgg := &BenchmarkResult{
			Name:           libName,
			CompressTime:   medCompTime,
			DecompressTime: medDecompTime,
			CompressedSize: medSize,
			Ratio:          medRatio,
		}
		final = append(final, libraryAgg)
	}
	return final
}

func findMedianTime(times []time.Duration) time.Duration {
	slices.SortFunc(times, func(a, b time.Duration) int {
		return int(a - b)
	})
	n := len(times)
	if n%2 == 0 {
		return (times[(n/2)-1] + times[n/2]) / 2
	}
	return times[n/2]
}

func findMedianInt(ints []int) int {
	slices.SortFunc(ints, func(a, b int) int {
		return a - b
	})
	n := len(ints)
	if n%2 == 0 {
		return (ints[(n/2)-1] + ints[n/2]) / 2
	}
	return ints[n/2]
}

func findMedianFloat64(floats []float64) float64 {
	slices.SortFunc(floats, func(a, b float64) int {
		return int(a - b)
	})
	n := len(floats)
	if n%2 == 0 {
		return (floats[(n/2)-1] + floats[n/2]) / 2
	}
	return floats[n/2]
}

func printResults(input []byte, opts *PrintOptions, results []*BenchmarkResult) {
	tw := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)
	printers := printResultHeader(tw, input, opts)
	for _, result := range results {
		printResultRow(tw, result, printers)
	}
	tw.Flush()
}

// prints the top row of the result table, and returns a list of formatting functions for all other rows
func printResultHeader(tw *tabwriter.Writer, input []byte, opts *PrintOptions) []func(*BenchmarkResult) string {
	if opts.ShouldPrintInput {
		fmt.Printf("Input data: %v\n", input)
	}
	fmt.Printf("Original data size: %s\n", formatBytes(len(input)))
	fields := []string{"Compression-Library"}
	printers := []func(*BenchmarkResult) string{
		func(br *BenchmarkResult) string {
			return br.Name
		},
	}
	if opts.ShouldPrintCTime {
		fields = append(fields, "Compression-Time")
		printers = append(printers, func(br *BenchmarkResult) string {
			return br.CompressTime.String()
		})
	}
	if opts.ShouldPrintDTime {
		fields = append(fields, "Decompression-Time")
		printers = append(printers, func(br *BenchmarkResult) string {
			return br.DecompressTime.String()
		})
	}
	fields = append(fields, "Total-Time")
	printers = append(printers, func(br *BenchmarkResult) string {
		return br.GetTotalTime().String()
	})
	fields = append(fields, "Compressed-Size")
	printers = append(printers, func(br *BenchmarkResult) string {
		return formatBytes(br.CompressedSize)
	})
	fields = append(fields, "Ratio")
	printers = append(printers, func(br *BenchmarkResult) string {
		return formatRatio(br.Ratio)
	})
	if opts.NetworkSpeed != 0 {
		fields = append(fields, fmt.Sprintf("%d-Payloads", opts.NetworkPayloads))
		printers = append(printers, func(br *BenchmarkResult) string {
			return br.GetBatchTime(opts.NetworkPayloads, opts.NetworkSpeed).String()
		})
	}
	fmt.Fprintln(tw, strings.Join(fields, "\t"))
	return printers
}

func printResultRow(tw *tabwriter.Writer, result *BenchmarkResult, printers []func(*BenchmarkResult) string) {
	entries := []string{}
	for _, printer := range printers {
		entries = append(entries, printer(result))
	}
	outStr := strings.Join(entries, "\t")
	fmt.Fprintln(tw, outStr)
}

// convert float64 to string
func formatRatio(r float64) string {
	return fmt.Sprintf("%.2f%%", (r * 100.0))
}

// convert file size to string
func formatBytes(size int) string {
	unitLadder := []string{"B", "KB", "MB", "GB"}
	unit := 0
	f := float64(size)
	for f >= 1000 {
		f /= 1000
		unit++
	}
	return fmt.Sprintf("%.4f %s", f, unitLadder[unit])
}

// Gets input for compression and decompression
func getBenchmarkInput(cmd *cobra.Command) ([]byte, error) {
	isRand, err := cmd.Flags().GetBool(isRandInput)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	var input []byte
	if isRand {
		// user wants to randomly generate JSON input
		input, err = getRandInput(cmd)
		if err != nil {
			return nil, fmt.Errorf("error while generating random input: %v", err)
		}
	} else {
		// user wants to read input from a file
		filename, err := cmd.Flags().GetString(fileInputFlag)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		input, err = os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error while reading input file: %v", err)
		}
	}
	return input, nil
}

// randomly generates a JSON object according to user flags
func getRandInput(cmd *cobra.Command) ([]byte, error) {
	generator, err := getJsonGenerator(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to setup random JSON generator: %v", err)
	}
	randJson, err := generator.JsonGenerate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random JSON: %v", err)
	}
	return json.Marshal(randJson)
}

func getJsonGenerator(cmd *cobra.Command) (JsonGenerator, error) {
	numFieldMin, numFieldMax, err := getNumFields(cmd)
	if err != nil {
		return nil, err
	}
	numChMin, numChMax, maxDepth, err := getNumChildren(cmd)
	if err != nil {
		return nil, err
	}
	minStrLen, maxStrLen, dictSize, dictFile, err := getStrFlags(cmd)
	if err != nil {
		return nil, err
	}
	jsonConfig := NewJsonGenConfig()
	jsonConfig.FieldsPerNodeMax = numFieldMax
	jsonConfig.FieldsPerNodeMin = numFieldMin
	jsonConfig.DegreeMin = numChMin
	jsonConfig.DegreeMax = numChMax
	jsonConfig.DepthMax = maxDepth
	jsonConfig.StrLenMin = minStrLen
	jsonConfig.StrLenMax = maxStrLen
	jsonConfig.DictFile = dictFile
	jsonConfig.DictSize = dictSize
	return newJsonGenerator(jsonConfig), nil
}
