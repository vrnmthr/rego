package eval

import (
	"github.com/open-policy-agent/opa/ast"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/opa-policy"
	"io/ioutil"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"bytes"
)

// Parses the file specified by data
func ParseBytes(fname string, data []byte) (*ast.Module, error) {
	return ast.ParseModule(fname, string(data))
}

// Parses the file specified by fpath and returns a module
func ParseFile(fpath string) (*ast.Module, error) {
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return ast.ParseModule(fpath, string(file))
}

// Parses all the files in fpaths and returns a map[string]*ast.DeserializeModuleJson where the filenames are the keys
func ParseFiles(fpaths []string) (map[string]*ast.Module, error) {
	errs := new(policy.Errors)
	modules := make(map[string]*ast.Module)
	for _, fpath := range fpaths {
		module, err := ParseFile(fpath)
		errs.Add(err)
		modules[fpath] = module
	}
	return modules, errs.NilIfEmpty()
}

// SerializeModuleJson converts into a JSON document. The document should always be pre-compiled and checked for correctness
// as the JSON document does not store the locations of any of the elements in the AST. This will make error messages
// during compilation very difficult.
func SerializeModuleJson(module *ast.Module) ([]byte, error) {
	return json.Marshal(*module)
}

// Reads a module from a Json byte array. The AST produced has no location fields.
func DeserializeModuleJson(data []byte) (*ast.Module, error) {
	var module = &ast.Module{}
	err := json.Unmarshal(data, module)
	if err != nil {
		return nil, err
	}
	addModuleToRules(module)
	return module, nil
}

// SerializeModuleGob uses Gob to convert into a byte array. The disadvantage of this method is that it is less space
// efficient as it stores all the location fields
func SerializeModuleGob(module *ast.Module) ([]byte, error) {
	buf := new(bytes.Buffer)
	removeModuleFromRules(module)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(*module)
	return buf.Bytes(), err
}

// DeserializeModuleGob uses Gob to deserialize.
func DeserializeModuleGob(data []byte) (*ast.Module, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var mod ast.Module
	err := dec.Decode(&mod)
	addModuleToRules(&mod)
	return &mod, err
}

// Removes the circular module reference from each rule
func removeModuleFromRules(module *ast.Module) {
	for _, rule := range module.Rules {
		rule.Module = nil
	}
}

// Adds the module reference to each rule
func addModuleToRules(module *ast.Module) {
	for _, rule := range module.Rules {
		rule.Module = module
	}
}

// reflect based conversion of all fields to nil
// Used for testing
func convertLocationToNil(value reflect.Value) {
	switch value.Kind(){
	case reflect.Ptr:
		convertLocationToNil(value.Elem())
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			convertLocationToNil(value.Index(i))
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			convertLocationToNil(value.MapIndex(key))
		}
	case reflect.Struct:
		loc := value.FieldByName("Location")
		if loc.IsValid() {
			if loc.CanSet() {
				loc.Set(reflect.Zero(loc.Type()))
			}
		}
		for i := 0; i < value.NumField(); i++ {
			convertLocationToNil(value.Field(i))
		}
	}
}

// removeLocations sets all the location tags of the struct to nil
func removeLocation(mod *ast.Module) {
	for _, rule := range mod.Rules {
		rule.Module = nil
	}
	convertLocationToNil(reflect.ValueOf(mod))
	for _, rule := range mod.Rules {
		rule.Module = mod
	}
}
//
////func removeLocationFromTerm(term *ast.Term) {
////	term.Location = nil
//	switch val := term.Value.(type) {
//	case ast.Array:
//		for _, i := range val {
//			removeLocationFromTerm(i)
//		}
//	case ast.Call:
//		for _, i := range val {
//			removeLocationFromTerm(i)
//		}
//	case ast.Ref:
//		for _, i := range val {
//			removeLocationFromTerm(i)
//		}
//	case *ast.ArrayComprehension:
//		removeLocationFromTerm(val.Term)
//		for _, i := range val.Body {
//			removeLocationFromExpr(i)
//		}
//	case ast.Set:
//	case *ast.SetComprehension:
//	case ast.Boolean, ast.Number, ast.Null, ast.Var, ast.String:
//		return
//	default:
//		panic("unrecognized term")
//	}
//}
//
//func removeLocationFromExpr(expr *ast.Expr) {
//	expr.Location = nil
//
//	for _, w := range expr.With {
//		w.Location = nil
//		removeLocationFromTerm(w.Value)
//		removeLocationFromTerm(w.Target)
//	}
//}

func init() {
	//gob.Register(ast.Var(""))
	//gob.Register(ast.String(""))
	//gob.Register(ast.Ref{})
	//gob.Register(ast.Array{})
	//gob.Register(ast.ArrayComprehension{})
	//gob.Register(ast.Boolean(true))
	//gob.Register(ast.Builtin{})
	//gob.Register(ast.Call{})
	//gob.Register(ast.Comment{})
	//gob.Register(ast.Null{})
	////gob.Register(ast.Object(nil))
	//gob.Register(ast.ObjectComprehension{})
	////gob.Register(ast.Set(nil))
	//gob.Register(ast.SetComprehension{})
	//gob.Register(ast.Module{})
	//gob.Register(ast.Error{})
	//gob.Register(ast.Errors{})
	//gob.Register(ast.Body{})
	//gob.Register(ast.Number(json.Number("5")))
	//gob.Register(ast.Term{})
	//gob.Register(ast.Head{})
	//gob.Register(ast.Rule{})
	//gob.Register(ast.RuleSet{})
	//gob.Register(ast.With{})
	//gob.Register(ast.Expr{})
	//gob.Register(ast.Location{})
	//gob.Register(ast.Args{})
	////gob.Register(ast.DocKind())
	//gob.Register(ast.Import{})
	//gob.Register(ast.Package{})
	//gob.Register(ast.GenericTransformer{})
	//gob.Register(ast.GenericVisitor{})
	//gob.Register(ast.ArgErrDetail{})
	//gob.Register(ast.RefErrInvalidDetail{})
	//gob.Register(ast.RefErrUnsupportedDetail{})
	//gob.Register(ast.UnificationErrDetail{})
}
