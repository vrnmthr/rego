package rego

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/storage"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"context"
	"github.com/pkg/errors"
)

// Query returns a ResultSet for the given query run on the given compiler
func Query(cmp *ast.Compiler, query string, inputs map[string]interface{}, store *storage.Store) (rego.ResultSet, error) {
	args := []func(r *rego.Rego){
		rego.Query(query),
		rego.Compiler(cmp),
		rego.Input(inputs),
	}

	if store != nil {
		args = append(args, rego.Store(*store))
	}

	rg := rego.New(args...)

	// will return rego_unsafe_var if junk in query
	rs, err := rg.Eval(context.Background())
	if err != nil {
		return nil, NewEvalError(query + ": " + err.Error())
	}

	return rs, nil
}

// QueryRule makes a query and returns a *single* value of any type that is produced by evaluation. If multiple objects
// are produced upon evaluation or no object is produced, error != nil.
func QueryRule(cmp *ast.Compiler, pkg, rule string, inputs map[string]interface{}, store *storage.Store) (interface{}, error) {
	q := fmt.Sprintf("data.%v.%v", pkg, rule)
	rs, err := Query(cmp, q, inputs, store)
	if err != nil {
		return nil, err
	}

	if len(rs) == 0 {
		msg := fmt.Sprintf("%v: query undefined", rule)
		return nil, NewUndefinedError(msg)
	}

	if len(rs) > 1 {
		err = errors.Wrap(fmt.Errorf("multiple results produced"), rule)
	}

	if len(rs[0].Expressions) == 0 {
		err = errors.Wrap(fmt.Errorf("no values produced by this rule"), rule)
	}

	return rs[0].Expressions[0].Value, err
}

func IsUndefined(err error) bool {
	_, ok := err.(*UndefinedErr)
	return ok
}

func IsEvalErr(err error) bool {
	_, ok := err.(*EvalErr)
	return ok
}

