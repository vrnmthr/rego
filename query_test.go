package eval

import (
	"testing"
	"github.com/open-policy-agent/opa/ast"
)

func setup(policy string) (*ast.Compiler) {
	m, err := ParseBytes("test", []byte(policy))
	if err != nil {
		panic(err)
	}
	cmp := NewCompiler()
	err = Compile(cmp, map[string]*ast.Module{"test": m})
	if err != nil {
		panic(err)
	}
	return cmp
}

func validate(t *testing.T, res, expected interface{}) {
	ok, err := areEqualJson(res, expected)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !ok {
		t.Fatalf("got %v, expected %v", res, expected)
	}
}

func TestQueryRuleUndefined(t *testing.T) {
	policy := `
	package test
	eval { false }
	`
	cmp := setup(policy)
	_, err := QueryRule(cmp, "test", "eval", nil, nil)
	if err == nil {
		t.Fatalf("did not catch undefined error")
	}
	if !IsUndefined(err) {
		t.Fatalf("incorrect error type")
	}
}

func TestQueryRuleNotFound(t *testing.T) {
	policy := `
	package test
	eval { false }
	`
	cmp := setup(policy)
	_, err := QueryRule(cmp, "test", "asdfasd", nil, nil)
	if err == nil {
		t.Fatalf("did not catch undefined error")
	}
	if !IsUndefined(err) {
		t.Fatalf("incorrect error type")
	}
}

func TestQueryRuleSingle(t *testing.T) {
	policy := `
	package test
	
	tester {true}
	`
	cmp := setup(policy)
	res, err := QueryRule(cmp, "test", "tester", nil, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := true
	validate(t, res, expected)
}

func TestQueryRuleMultiple(t *testing.T) {
	policy := `
	package test
	
	lst = [1,2,3]
	eval[v] { v = lst[_] }
	`
	cmp := setup(policy)
	res, err := QueryRule(cmp, "test", "eval", nil, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	validate(t, res, []int{1,2,3})
}

func TestQueryRuleInputs(t *testing.T) {
	policy := `
	package test
	
	lst = [1,2,3]
	eval = x { x := input.a + input.b }
	`
	cmp := setup(policy)
	inputs := map[string]interface{}{"a":1, "b":2}
	res, err := QueryRule(cmp, "test", "eval", inputs, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	validate(t, res, 3)
}

func TestQueryRuleEvalErr(t *testing.T)  {
	policy := `
	package test
	eval { http.send({}) }
	`
	cmp := setup(policy)
	_, err := QueryRule(cmp, "test", "eval", nil, nil)
	if err == nil {
		t.Fatalf("did not catch eval error")
	}
	if !IsEvalErr(err) {
		t.Fatalf("incorrect error type")
	}
}