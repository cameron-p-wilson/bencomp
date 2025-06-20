package main

import "testing"

type JsonConstraint func(t *testing.T, depth int, elem *JsonElement)

func checkConstraints(t *testing.T, elem *JsonElement, constraints []JsonConstraint) {
	type StackTuple struct {
		elem  *JsonElement
		depth int
	}
	stack := []StackTuple{{elem: elem, depth: 0}}
	for len(stack) > 0 {
		curTuple := stack[len(stack)-1]
		curElem := curTuple.elem
		curDepth := curTuple.depth
		stack = stack[:len(stack)-1]
		for _, constraint := range constraints {
			constraint(t, curDepth, curElem)
		}
		for _, child := range curElem.Children {
			newTuple := StackTuple{elem: child, depth: curDepth + 1}
			stack = append(stack, newTuple)
		}
	}
}

func getFieldConstraint(min, max int) JsonConstraint {
	return func(t *testing.T, _ int, elem *JsonElement) {
		n := len(elem.Fields)
		if n < min || n > max {
			t.Errorf("constraint failed: element had %d fields, but expected between %d and %d", n, min, max)
			t.FailNow()
		}
	}
}

func getChildrenConstraint(min, max, maxDepth int) JsonConstraint {
	return func(t *testing.T, depth int, elem *JsonElement) {
		n := len(elem.Children)
		if depth == maxDepth {
			if n != 0 {
				t.Errorf("constraint failed: tree exceeded depth %d", maxDepth)
				t.FailNow()
			}
			return
		}
		if n < min || n > max {
			t.Errorf("constraint failed: element had %d children, but expected between %d to %d", n, min, max)
			t.FailNow()
		}
	}
}

func getStrConstraint(min, max int) JsonConstraint {
	return func(t *testing.T, _ int, elem *JsonElement) {
		for k, v := range elem.Fields {
			n := len(k)
			if n < min || n > max {
				t.Errorf("constraint failed: element had key of size %d, but expected between %d to %d", n, min, max)
				t.FailNow()
			}
			n = len(v)
			if n < min || n > max {
				t.Errorf("constraint failed: element had value of size %d, but expected between %d to %d", n, min, max)
				t.FailNow()
			}
		}
	}
}

func TestJsonGenerate(t *testing.T) {
	type testData struct {
		name        string
		input       *JsonGenConfig
		constraints []JsonConstraint
		wantErr     bool
	}
	tests := []testData{
		{
			name: "default",
			input: &JsonGenConfig{
				FieldsPerNodeMin: 3,
				FieldsPerNodeMax: 3,
				DegreeMin:        4,
				DegreeMax:        4,
				DepthMax:         5,
				StrLenMin:        16,
				StrLenMax:        16,
			},
			constraints: []JsonConstraint{
				getFieldConstraint(3, 3),
				getChildrenConstraint(4, 4, 5),
				getStrConstraint(16, 16),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsongen := NewJsonGenerator(test.input)
			output, err := jsongen.JsonGenerate()
			if test.wantErr {
				if err == nil {
					t.Errorf("expected error, but did not get one")
				}
				return
			}
			checkConstraints(t, output, test.constraints)
		})
	}
}
