package main

import (
	"reflect"
	"testing"
)

// TestJsonGenerator implements the JsonGenerator interface for mocking
type TestJsonGenerator struct {
	config *JsonGenConfig
}

func (tgen *TestJsonGenerator) JsonGenerate() (*JsonElement, error) {
	return &JsonElement{
		Fields: map[string]string{
			"foo": "bar",
		},
	}, nil
}

func testConfigEqual(t *testing.T, actual, exp *JsonGenConfig) {
	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("actual and expected configs are not equal, actual: %v, exp: %v", actual, exp)
	}
}

func TestBenchCmd(t *testing.T) {
	type testData struct {
		name          string
		args          []string
		expConfig     JsonGenConfig
		wantErr       bool
		wantNilConfig bool
	}
	tests := []testData{
		{
			name: "default",
			args: []string{"--rand-gen"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name:          "file input test",
			args:          []string{"--file", "./bench_test.go"},
			wantNilConfig: true,
		},
		{
			name: "static num fields",
			args: []string{"--rand-gen", "--json-num-fields", "5"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: 5,
				FieldsPerNodeMax: 5,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "range num fields",
			args: []string{"--rand-gen", "--json-num-fields-range", "4-6"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: 4,
				FieldsPerNodeMax: 7, // range is half open, so "4-6" => [4,7)
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "max depth",
			args: []string{"--rand-gen", "--json-max-depth", "2"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         2,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "static degree",
			args: []string{"--rand-gen", "--json-degree", "7"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        7,
				DegreeMax:        7,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "range degree",
			args: []string{"--rand-gen", "--json-degree-range", "1-8"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        1,
				DegreeMax:        9, // range is half open, so "1-8" => [1,9)
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "static string length",
			args: []string{"--rand-gen", "--json-str-len", "12"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        12,
				StrLenMax:        12,
			},
		},
		{
			name: "range string length",
			args: []string{"--rand-gen", "--json-str-len-range", "4-12"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        4,
				StrLenMax:        13, // range is half open, so "4-12" => [4,13)
			},
		},
		{
			name: "json file dictionary",
			args: []string{"--rand-gen", "--json-dict-file", "./bench_test.go"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
				DictFile:         "./bench_test.go",
			},
		},
		{
			name: "json random dictionary",
			args: []string{"--rand-gen", "--json-dict-size", "10"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
				DictSize:         10,
			},
		},
		{
			name: "network speed 1000",
			args: []string{"--rand-gen", "--network-bandwidth", "1000"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "network speed 128KB",
			args: []string{"--rand-gen", "--network-bandwidth", "128KB"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "network speed 256MB",
			args: []string{"--rand-gen", "--network-bandwidth", "256MB"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name: "network speed 10GB",
			args: []string{"--rand-gen", "--network-bandwidth", "10GB"},
			expConfig: JsonGenConfig{
				FieldsPerNodeMin: defaultFieldNum,
				FieldsPerNodeMax: defaultFieldNum,
				DegreeMin:        defaultDegree,
				DegreeMax:        defaultDegree,
				DepthMax:         defaultMaxDepth,
				StrLenMin:        defaultJsonStrLen,
				StrLenMax:        defaultJsonStrLen,
			},
		},
		{
			name:    "network speed invalid",
			args:    []string{"--rand-gen", "--network-bandwidth", "10 B"},
			wantErr: true,
		},
		{
			name:    "network speed overflow",
			args:    []string{"--rand-gen", "--network-bandwidth", "100000 GB"},
			wantErr: true,
		},
		{
			name:    "error no mode",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "error both modes",
			args:    []string{"--rand-gen", "--file", "./bench_test.go"},
			wantErr: true,
		},
	}
	var outConfig *JsonGenConfig
	newJsonGenerator = func(config *JsonGenConfig) JsonGenerator {
		outConfig = config
		return &TestJsonGenerator{
			config: config,
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outConfig = nil
			cmd := NewBenchCmd()
			cmd.SetArgs(test.args)
			err := cmd.Execute()
			if test.wantErr {
				if err == nil {
					t.Errorf("expected error, but did not get one")
				}
				return
			}
			if test.wantNilConfig {
				if outConfig != nil {
					t.Errorf("wantNilConfig is true but got %v", outConfig)
				}
				return
			}
			testConfigEqual(t, outConfig, &test.expConfig)
		})
	}
}
