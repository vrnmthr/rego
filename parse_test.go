package eval

import (
	"testing"
	"github.com/open-policy-agent/opa/ast"
	"bytes"
	"io/ioutil"
	"os"
)

func testSerialization(def string, t *testing.T) {
	// compile and parse module
	mod, err:= ParseBytes("test", []byte(def))
	if err != nil {
		t.Fatalf(err.Error())
	}
	cmp := NewCompiler()
	err = Compile(cmp, map[string]*ast.Module{"test" : mod})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// serialize and deserialize
	data, err := SerializeModuleJson(mod)
	if err != nil {
		t.Fatalf(err.Error())
	}
	ioutil.WriteFile("test.json", data, os.ModePerm)

	result, err := DeserializeModuleJson(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// make sure policy is compilable
	c := NewCompiler()
	err = Compile(c, map[string]*ast.Module{"test" : result})
	if err != nil {
		t.Fatalf("uncompilable: " + err.Error())
	}

	//if !reflect.DeepEqual(*result, *mod) {
	//	t.Fatalf("not deeply equal")
	//}

	// two are equal if the pretty printed representations of the two are equal
	ast1 := new(bytes.Buffer)
	ast.Pretty(ast1, mod)
	ast2 := new(bytes.Buffer)
	ast.Pretty(ast2, mod)

	if string(ast1.Bytes()) != string(ast2.Bytes()) {
		t.Logf("Different parse trees produced:\n")
		t.Logf("Input:\n")
		t.Logf(string(ast1.Bytes()))
		t.Logf("\nOutput:\n")
		t.Fatalf(string(ast2.Bytes()))
	}
}

func TestModuleSerializationSimple(t *testing.T) {
	def := `
	package test
	default eval = true
	`
	testSerialization(def, t)
}

func TestModuleSerializationComplex(t *testing.T) {
	def := `
	package test
	import data.otherpackage
	
	v1 = input.arg
	v2 = "hello"
	v3 = 123.1324
	v4 = [1,2,3]
	v5 = [f | f = v4[_]]
	v6 = plus(6, 4)
	v7 = false
	# comment
	
	default boolrule = false
	boolrule = true {
		obj := v4
		plus(6,5) == 11
	} {
		v2
		v6
	}


	setrule[result] {
		result = v5[_]
	}
	`
	testSerialization(def, t)
}