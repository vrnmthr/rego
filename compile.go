package eval

import "github.com/open-policy-agent/opa/ast"

// Creates a new Compiler
func NewCompiler() (*ast.Compiler) {
	return ast.NewCompiler()
}

// Compile the modules with the specified compiler
func Compile(cmp *ast.Compiler, modules map[string]*ast.Module) (error) {
	cmp.Compile(modules)
	if cmp.Failed() {
		return cmp.Errors
	}
	return nil
}
