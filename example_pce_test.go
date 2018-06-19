package rego

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/storage/inmem"
	"fmt"
)

var PolicyDoc = `
	package example
	
	default rule = false
	rule = x {
		x := input.a + input.b + data.a + data.b
	}
`

// Indicates how to use this package to parse, save the parsed representation to a file,
// compile the file, and query the final compiled representation.
func Example()  {

	// Parse
	mod, _ := ParseBytes("example", []byte(PolicyDoc))

	// CHECK THAT MODULE CAN BE COMPILED BEFORE SERIALIZATION
	cmp := NewCompiler()
	if Compile(cmp, map[string]*ast.Module {"example": mod}) != nil {
		panic("compilation failed")
	}

	bytes, _ := SerializeModuleJson(mod)
	// ... do whatever here (store in database, maybe?)
	mod, _ = DeserializeModuleJson(bytes)

	// Compile the module
	cmp = NewCompiler()
	modules := map[string]*ast.Module {
		"example": mod,
	}
	Compile(cmp, modules)

	// Query
	inputs := map[string]interface{} {
		"a": 1,
		"b": 2,
	}

	data := map[string]interface{} {
		"a": 3,
		"b": 4,
	}

	store := inmem.NewFromObject(data)
	res, _ := QueryRule(cmp, "example", "rule", inputs, &store)
	fmt.Println(res)

	// Output:
	// 10
}
