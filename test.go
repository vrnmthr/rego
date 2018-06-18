package rego

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"strings"
	"testing"
	"bytes"
	"reflect"
	"encoding/json"
)

const (
	UNDEF = "---undefined---"
)

func compileModules(input []string) *ast.Compiler {

	mods := map[string]*ast.Module{}

	for idx, i := range input {
		id := fmt.Sprintf("testMod%d", idx)
		mods[id] = ast.MustParseModule(i)
	}

	c := ast.NewCompiler()
	if c.Compile(mods); c.Failed() {
		panic(c.Errors)
	}

	return c
}

func compileRules(pkg string, input []string) (*ast.Compiler, error) {

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("package %v\n", pkg))
	buf.WriteString(strings.Join(input, " \n\n"))

	parsed, err := ParseBytes("test", buf.Bytes())
	if err != nil {
		panic(err)
	}

	c := ast.NewCompiler()
	if c.Compile(map[string]*ast.Module{"testMod": parsed}); c.Failed() {
		return nil, c.Errors
	}

	return c, nil
}

// Parses object into JSON type
func toJson(obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(obj)
	return buf.Bytes(), err
}

// Credit to @turtlemonvh: https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func areEqualJson(arg1, arg2 interface{}) (bool, error) {

	a, err := toJson(arg1)
	if err != nil {
		return false, err
	}
	b, err := toJson(arg2)
	if err != nil {
		return false, err
	}

	var o1 interface{}
	var o2 interface{}

	err = json.Unmarshal(a, &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1: %s", err.Error())
	}
	err = json.Unmarshal(b, &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2: %s", err.Error())
	}
	return reflect.DeepEqual(o1, o2), nil
}

// TestCase represents a single test. Rule is the rule to be queried for. It defaults to "t".
// Rules should be Rego rules.
type TestCase struct {
	Rule string
	Rules []string
	Expected interface{}
}

// RunTestCase runs the given test with the given inputs and data document. It annotates the test with note.
// Although test.Expected is an interface{}, Rego queries can only return JSON values, so format your interface
// accordingly (i.e should be map[string]interface{}, []interface{}, etc.). An incorrect expected value will
// simply return false.
func (test *TestCase) Run(t *testing.T, inputs, data map[string]interface{}, note string) {
	t.Run(note, func(t2 *testing.T) {
		err := runTestCase(inputs, data, test)
		if err != nil {
			t2.Fatalf(err.Error())
		}
	})
}

func runTestCase(inputs, data map[string]interface{}, test *TestCase) (error) {
	pkg := "testing"
	compiler, err := compileRules(pkg, test.Rules)
	if err != nil {
		return err
	}

	if len(test.Rule) == 0 {
		test.Rule = "t"
	}

	var store storage.Store = nil
	if data != nil {
		store = inmem.NewFromObject(data)
	}

	path := "data." + pkg
	return assertWithPath(compiler, inputs, store, test.Rule, path, test.Expected)
}

// TestCase the given file. Rule must be specified
func RunTestFile(t *testing.T, inputs, data map[string]interface{}, file, rule, note string, expected interface{}) {
	module, err := ParseBytes("test", []byte(file))
	if err != nil {
		t.Fatalf(err.Error())
	}
	cmp := NewCompiler()
	err = Compile(cmp, map[string]*ast.Module{"testMod": module})
	if err != nil {
		t.Fatalf(err.Error())
	}

	var store storage.Store = nil
	if data != nil {
		store = inmem.NewFromObject(data)
	}

	assertWithPath(cmp, inputs, store, rule, module.Package.Path.String(), expected)
}

func assertWithPath(compiler *ast.Compiler, inputs map[string]interface{}, store storage.Store,
	rule, path string, expected interface{}) (error) {

	q := fmt.Sprintf("%v.%v", path, rule)

	switch e := expected.(type) {
	case error:
		rs, err := Query(compiler, q, inputs, &store)
		if err == nil {
			return fmt.Errorf("expected error but got: %v", rs)
		}
		if !strings.Contains(err.Error(), e.Error()) {
			return fmt.Errorf("expected error %v but got: %v", e, err)
		}
	default:

		rs, err := Query(compiler, q, inputs, &store)

		if err != nil {
			return fmt.Errorf("unexpected error: %v", err)
		}

		if expected == UNDEF {
			if len(rs) != 0 {
				return fmt.Errorf("expected undefined result but got: %v", rs)
			}
			return nil
		}

		if len(rs) == 0 {
			return fmt.Errorf("expected %v but got undefined", e)
		}

		// compare the two
		if len(rs[0].Expressions) == 0 {
			return fmt.Errorf("no expressions found upon evaluation")
		}

		result := rs[0].Expressions[0].Value
		eq, err := areEqualJson(expected, result)
		if err != nil {
			panic(err)
		}
		if !eq {
			return fmt.Errorf("expected %v, got %v", expected, result)
		}
	}
	return nil
}

