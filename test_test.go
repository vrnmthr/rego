package rego

import (
	"testing"
	"fmt"
)

func TestExpectedUndefined(t *testing.T) {
	test := TestCase{
		Target:   "garbage",
		Rules:    nil,
		Expected: UNDEF,
	}
	test.Run(t, nil, nil)
}

func TestUnexpectedUndefined(t *testing.T)  {
	test := TestCase{
		Rules: []string{"t { false }"},
		Expected: true,
	}
	err := runTestCase(nil, nil, &test)
	if err == nil {
		t.Fatalf("did not catch unexpected undefined")
	}
}

func TestExpectedError(t *testing.T) {
	test := TestCase{
		Rules: []string{"t { x := http.send({}) }"},
		Expected: fmt.Errorf("bad args"),
	}
	err := runTestCase(nil, nil, &test)
	if err == nil {
		t.Fatalf("did not expect error")
	}
}

func TestUnexpectedError(t *testing.T)  {
	test := TestCase{
		Rules: []string{"t { x := http.send({}) }"},
		Expected: "fine!",
	}
	err := runTestCase(nil, nil, &test)
	if err == nil {
		t.Fatalf("did not catch unexpected error")
	}
}

func TestDefaultRuleWorks(t *testing.T) {
	test := TestCase{
		Rules: []string{"t = 10 {true}", "default b = 8"},
		Expected: 10,
	}
	test.Run(t, nil, nil)
}

func TestStringIntDifference(t *testing.T) {
	test := TestCase{
		Rules: []string{"t = 10 {true}", "default b = 8"},
		Expected: 10,
	}
	err := runTestCase(nil, nil, &test)
	if err != nil {
		t.Fatalf(err.Error())
	}

	test2 := TestCase{
		Rules: []string{"t = 10 {true}", "default b = 8"},
		Expected: "10",
	}
	err = runTestCase(nil, nil, &test2)
	if err == nil {
		t.Fatalf("did not catch type difference")
	}
}

func TestComplexEquality(t *testing.T) {
	hippo := struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Friends []string `json:"friends"`
	}{
		Name: "jim",
		Age: 123,
		Friends: []string{"tom", "ben"},
	}

	test := TestCase{
		Target:   "hippo",
		Rules:    []string{`hippo = {"name":"jim", "age": 123, "friends":["tom", "ben"]} {true}`},
		Expected: hippo,
	}

	err := runTestCase(nil, nil, &test)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestComplexInequality(t *testing.T) {
	hippo := struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Friends []string `json:"friends"`
	}{
		Name: "jim",
		Age: 123,
		Friends: []string{"tommy", "ben"},
	}

	test := TestCase{
		Target:   "hippo",
		Rules:    []string{`hippo = {"name":"jim", "age": 123, "friends":["tom", "ben"]} {true}`},
		Expected: hippo,
	}

	err := runTestCase(nil, nil, &test)
	if err == nil {
		t.Fatalf("did not catch error in field difference")
	}
}

func TestInputs(t *testing.T)  {
	inputs := map[string]interface{} {
		"arg": 10,
		"arg2": 12,
	}
	test := TestCase{
		Rules: []string{"t = x {x := input.arg + input.arg2}"},
		Expected: 22,
	}
	test.Run(t, inputs, nil)
}

func TestData(t *testing.T) {
	data := map[string]interface{} {
		"arg": 10,
		"arg2": 12,
	}
	test := TestCase{
		Rules: []string{"t = x {x := data.arg + data.arg2}"},
		Expected: 22,
	}
	test.Run(t, nil, data)
}

func TestInputsAndData(t *testing.T) {
	inputs := map[string]interface{} {
		"arg": 10,
		"arg2": 12,
	}
	data := map[string]interface{} {
		"arg3": 1,
		"arg4": 2,
	}
	test := TestCase{
		Rules: []string{"t = x {x := input.arg + input.arg2 + data.arg3 + data.arg4}"},
		Expected: 25,
	}
	test.Run(t, inputs, data)
}