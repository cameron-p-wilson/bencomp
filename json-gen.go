package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type JsonElement struct {
	Children []*JsonElement    `json:"children,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

// JsonGen implements the JsonGenerator interface
type JsonGen struct {
	config *JsonGenConfig
}

type JsonGenerator interface {
	JsonGenerate() (*JsonElement, error)
}

// returns a JsonGenConfig with default values
func NewJsonGenConfig() *JsonGenConfig {
	return &JsonGenConfig{}
}

func NewJsonGenerator(config *JsonGenConfig) JsonGenerator {
	return &JsonGen{
		config: config,
	}
}

type JsonGenConfig struct {
	FieldsPerNodeMin int
	FieldsPerNodeMax int
	DegreeMin        int
	DegreeMax        int
	DepthMax         int
	DictFile         string
	DictSize         int
	StrLenMin        int
	StrLenMax        int
	NetworkSpeed     uint64
}

func (conf *JsonGenConfig) numFieldsGetter() func() int {
	return func() int {
		return GetRandRange(conf.FieldsPerNodeMin, conf.FieldsPerNodeMax)
	}
}

func (conf *JsonGenConfig) numChildrenGetter() func(int) int {
	return func(depth int) int {
		if depth >= conf.DepthMax {
			return 0
		}
		return GetRandRange(conf.DegreeMin, conf.DegreeMax)
	}
}

func (conf *JsonGenConfig) strGetter() (func() string, error) {
	if conf.DictFile != "" {
		// user specified dictionary file
		file, err := os.Open(conf.DictFile)
		if err != nil {
			return nil, fmt.Errorf("error reading dictionary file: %v", err)
		}
		defer file.Close()
		dict := make([]string, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			dict = append(dict, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading dictionary file: %v", err)
		}
		return func() string {
			i := GetRandRange(0, len(dict))
			return dict[i]
		}, nil
	}
	if conf.DictSize > 0 {
		// user wants randomly generated dictionary
		dict := make([]string, conf.DictSize)
		for i := range dict {
			n := GetRandRange(conf.StrLenMin, conf.StrLenMax)
			dict[i] = RandNChars(n)
		}
		return func() string {
			i := GetRandRange(0, len(dict))
			return dict[i]
		}, nil
	}
	// user wants randomly generated strings
	return func() string {
		n := GetRandRange(conf.StrLenMin, conf.StrLenMax)
		return RandNChars(n)
	}, nil
}

func (gen *JsonGen) JsonGenerate() (*JsonElement, error) {
	strGetter, err := gen.config.strGetter()
	if err != nil {
		return nil, fmt.Errorf("failed to generate json: %v", err)
	}
	return gen.jsonGenerateRec(strGetter, 0), nil
}

func (gen *JsonGen) jsonGenerateRec(strGetter func() string, depth int) *JsonElement {
	numChildren := gen.config.numChildrenGetter()
	numFields := gen.config.numFieldsGetter()
	out := JsonElement{}
	nf := numFields()
	if nf >= 1 {
		out.Fields = make(map[string]string)
		for range nf {
			out.Fields[strGetter()] = strGetter()
		}
	}
	nc := numChildren(depth)
	if nc >= 1 {
		out.Children = make([]*JsonElement, 0, nc)
		for range nc {
			out.Children = append(out.Children, gen.jsonGenerateRec(strGetter, depth+1))
		}
	}
	return &out
}

func NumStatic(n int) func(int) int {
	return func(_ int) int {
		return n
	}
}
func NumByMaxDepth(max int) func(int) int {
	return func(depth int) int {
		return max - depth
	}
}

func RandNChars(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = GetRandLowercase()
	}
	return string(bytes)
}

func GetRandLowercase() byte {
	return byte(GetRandRange(97, 122))
}

func GetRandRange(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

func GenDictionary(wordSize, dictSize int) []string {
	s := make([]string, 0, dictSize)
	for _ = range dictSize {
		s = append(s, RandNChars(wordSize))
	}
	return s
}

func GetRandFromDictionary(dict []string) func(int) string {
	return func(_ int) string {
		return dict[rand.Intn(len(dict)-1)]
	}
}

func GetRandFromDictionaryAny(dict []string) func(int) any {
	return func(_ int) any {
		return dict[rand.Intn(len(dict)-1)]
	}
}
