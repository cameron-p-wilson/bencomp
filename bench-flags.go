package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	// random JSON generation flags
	isRandInput                = "rand-gen"
	isRandInputShort           = "r"
	randJsonNumFieldsFlag      = "json-num-fields"
	randJsonNumFieldsRangeFlag = "json-num-fields-range"
	randJsonMaxDepthFlag       = "json-max-depth"
	randJsonDegreeFlag         = "json-degree"
	randJsonDegreeRangeFlag    = "json-degree-range"
	randJsonStrLenFlag         = "json-str-len"
	randJsonStrLenRangeFlag    = "json-str-len-range"
	randJsonDictFile           = "json-dict-file"
	randJsonDictSize           = "json-dict-size"

	// file input
	fileInputFlag      = "file"
	fileInputFlagShort = "f"

	// optional stats
	networkSpeedFlag    = "network-bandwidth"
	networkPayloadsFlag = "network-payloads"
	printCTimeFlag      = "show-compress-time"
	printDTimeFlag      = "show-decompress-time"
	printJsonFlag       = "show-input"
	countFlag           = "count"
	countFlagShort      = "c"

	// default values
	defaultFieldNum   = 3
	defaultMaxDepth   = 5
	defaultDegree     = 4
	defaultJsonStrLen = 16
)

func getPrintOptions(cmd *cobra.Command) (*PrintOptions, error) {
	speed, err := getSpeedFlag(cmd)
	if err != nil {
		return nil, err
	}
	shouldPrint, _ := cmd.Flags().GetBool(printJsonFlag)
	shouldPrintCTime, _ := cmd.Flags().GetBool(printCTimeFlag)
	shouldPrintDTime, _ := cmd.Flags().GetBool(printDTimeFlag)
	numPayloads := 1
	if numPayloadsInput, err := cmd.Flags().GetInt(networkPayloadsFlag); err == nil && numPayloadsInput != 0 {
		numPayloads = numPayloadsInput
	}
	return &PrintOptions{
		NetworkSpeed:     speed,
		ShouldPrintInput: shouldPrint,
		NetworkPayloads:  numPayloads,
		ShouldPrintCTime: shouldPrintCTime,
		ShouldPrintDTime: shouldPrintDTime,
	}, nil
}

func getCountFlag(cmd *cobra.Command) (int, error) {
	count, err := cmd.Flags().GetInt(countFlag)
	if err != nil {
		return 0, err
	}
	if count <= 0 {
		return 0, fmt.Errorf("value for %s must be a positive integer", countFlag)
	}
	return count, nil
}

func getSpeedFlag(cmd *cobra.Command) (uint64, error) {
	speedStr, err := cmd.Flags().GetString(networkSpeedFlag)
	if err != nil {
		return 0, err
	}
	if speed, err := parseSpeed(speedStr); err != nil {
		return 0, err
	} else {
		return speed, nil
	}
}

func parseSpeed(speedStr string) (out uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			out, err = 0, fmt.Errorf("invalid value '%s' for %s: %v", speedStr, networkSpeedFlag, r)
		}
	}()
	if speedStr == "" {
		return 0, nil
	}
	upper := strings.ToUpper(speedStr)
	mult := uint64(1)
	numEndIndex := len(upper)
	if upper[len(upper)-1] == 'B' {
		if len(upper) <= 2 {
			return 0, fmt.Errorf("invalid value '%s' for %s", speedStr, networkSpeedFlag)
		}
		char := upper[len(upper)-2]
		switch {
		case char == 'K':
			mult = 1000
			numEndIndex--
		case char == 'M':
			mult = 1000000
			numEndIndex--
		case char == 'G':
			mult = 1000000000
			numEndIndex--
		case char >= 48 && char <= 57:
			mult = 1
		default:
			return 0, fmt.Errorf("invalid value '%s' for %s", speedStr, networkSpeedFlag)
		}
	}
	valStr := upper[:numEndIndex-1]
	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return val * mult, nil
}

// returns the min string length, max string length, size of dictionary, and dictionary file name
func getStrFlags(cmd *cobra.Command) (minStrLen, maxStrLen, dictSize int, dictFile string, err error) {
	if jsonDictFile, err := cmd.Flags().GetString(randJsonDictFile); err == nil && jsonDictFile != "" {
		return defaultJsonStrLen, defaultJsonStrLen, 0, jsonDictFile, nil
	}
	strLenRangeStr, err := cmd.Flags().GetString(randJsonStrLenRangeFlag)
	if err == nil && strLenRangeStr != "" {
		min, max, err := parseRange(strLenRangeStr)
		if err != nil {
			return -1, -1, -1, "", fmt.Errorf("invalid argument for %s: %v", randJsonStrLenRangeFlag, err)
		}
		if max <= 0 || min <= 0 {
			return -1, -1, -1, "", fmt.Errorf("invalid argument for %s: must be 1 or greater", randJsonStrLenRangeFlag)
		}
		minStrLen = min
		maxStrLen = max
	}
	if jsonDictSize, err := cmd.Flags().GetInt(randJsonDictSize); err == nil && jsonDictSize != 0 {
		if jsonDictSize <= 0 {
			return -1, -1, -1, "", fmt.Errorf("invalid argument for %s: must be 1 or greater", randJsonDictSize)
		}
		dictSize = jsonDictSize
	}

	strLen, err := cmd.Flags().GetInt(randJsonStrLenFlag)
	if err == nil && strLenRangeStr == "" {
		if strLen <= 0 {
			return -1, -1, -1, "", fmt.Errorf("invalid argument for %s: must be 1 or greater", randJsonStrLenFlag)
		}
		minStrLen = strLen
		maxStrLen = strLen
	}
	return minStrLen, maxStrLen, dictSize, "", nil
}

// returns a function to determine the number of children for each element in the JSON tree
func getNumChildren(cmd *cobra.Command) (int, int, int, error) {
	maxDepth, err := cmd.Flags().GetInt(randJsonMaxDepthFlag)
	if err != nil {
		return -1, -1, -1, fmt.Errorf("invalid argument for %s: %v", randJsonMaxDepthFlag, err)
	}
	degreeRangeStr, err := cmd.Flags().GetString(randJsonDegreeRangeFlag)
	if err == nil && degreeRangeStr != "" {
		min, max, err := parseRange(degreeRangeStr)
		if err != nil {
			return -1, -1, -1, fmt.Errorf("invalid argument for %s: %v", randJsonDegreeFlag, err)
		}
		return min, max, maxDepth, nil
	}
	degreeStatic, err := cmd.Flags().GetInt(randJsonDegreeFlag)
	if err == nil && degreeStatic != 0 {
		return degreeStatic, degreeStatic, maxDepth, nil
	}
	return defaultDegree, defaultDegree, maxDepth, nil
}

// helper function to parse flags which define a range of values [min, max]
func parseRange(r string) (min, max int, err error) {
	split := strings.Split(r, "-")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("invalid range %s", r)
	}
	minStr := split[0]
	maxStr := split[1]
	if min, err = strconv.Atoi(minStr); err != nil {
		return 0, 0, fmt.Errorf("invalid range %s", r)
	}
	if max, err = strconv.Atoi(maxStr); err != nil {
		return 0, 0, fmt.Errorf("invalid range %s", r)
	}
	max++
	if min > max {
		return 0, 0, fmt.Errorf("min '%d' cannot be greater than max '%d'", min, max)
	}
	return min, max, nil
}

// returns a function to determine the number of key-value pairs for each element in the JSON tree
func getNumFields(cmd *cobra.Command) (int, int, error) {
	flagNum, err := cmd.Flags().GetInt(randJsonNumFieldsFlag)
	if err == nil && flagNum != 0 {
		return flagNum, flagNum, nil
	}
	flagRange, err := cmd.Flags().GetString(randJsonNumFieldsRangeFlag)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid argument for %s", randJsonNumFieldsRangeFlag)
	}
	if flagRange == "" {
		return defaultFieldNum, defaultFieldNum, nil
	}
	min, max, err := parseRange(flagRange)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid argument for %s: %v", randJsonDegreeRangeFlag, err)
	}
	return min, max, nil
}

func createBenchFlags(benchCmd *cobra.Command) {
	// random input
	benchCmd.Flags().BoolP(isRandInput, isRandInputShort, false, "If BenComp should randomly generate input data for benchmarking")
	benchCmd.Flags().Int(randJsonNumFieldsFlag, 0, "Fixed number of fields to populate in each JSON node")
	benchCmd.Flags().String(randJsonNumFieldsRangeFlag, "", "Min and max number of fields to populate in each JSON node")
	benchCmd.MarkFlagsMutuallyExclusive(randJsonNumFieldsFlag, randJsonNumFieldsRangeFlag)
	benchCmd.Flags().Int(randJsonMaxDepthFlag, defaultMaxDepth, "Maximum depth of the JSON tree")
	benchCmd.Flags().Int(randJsonDegreeFlag, 0, "Fixed number of children of each JSON node")
	benchCmd.Flags().String(randJsonDegreeRangeFlag, "", "Min and max number of children of each JSON node")
	benchCmd.MarkFlagsMutuallyExclusive(randJsonDegreeFlag, randJsonDegreeRangeFlag)
	benchCmd.Flags().Int(randJsonStrLenFlag, 16, "Fixed number of ASCII characters in each JSON string value")
	benchCmd.Flags().String(randJsonStrLenRangeFlag, "", "Min and max number of ASCII characters in JSON string values")
	benchCmd.MarkFlagsMutuallyExclusive(randJsonStrLenFlag, randJsonStrLenRangeFlag)
	// random input from dictionary
	benchCmd.Flags().String(randJsonDictFile, "", "File containing a dictionary of words to use in JSON string values")
	benchCmd.Flags().Int(randJsonDictSize, 0, "Number of words to randomly generate as a dictionary for JSON string values")
	benchCmd.MarkFlagsMutuallyExclusive(randJsonDictFile, randJsonDictSize)
	// debug
	benchCmd.Flags().Bool(printJsonFlag, false, "If BenComp should print the JSON it used in benchmarking")

	// file input
	benchCmd.Flags().StringP(fileInputFlag, fileInputFlagShort, "", "Specify a file to be used for compression benchmarking")
	benchCmd.MarkFlagsOneRequired(isRandInput, fileInputFlag)
	benchCmd.MarkFlagsMutuallyExclusive(isRandInput, fileInputFlag)

	// optional output
	benchCmd.Flags().String(networkSpeedFlag, "", "Number of bytes (not bits) per second on the wire, e.g. 128KB")
	benchCmd.Flags().Int(networkPayloadsFlag, 0, "Number of payloads used in system performance estimate")
	benchCmd.Flags().Bool(printCTimeFlag, false, "If set, will display time spent compressing in a separate column")
	benchCmd.Flags().Bool(printDTimeFlag, false, "If set, will display time spent decompressing in a separate column")
	benchCmd.Flags().IntP(countFlag, countFlagShort, 1, "Repeat the benchmark multiple times and record the median values")
}
