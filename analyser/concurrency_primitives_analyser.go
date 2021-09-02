package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/packages"
)

const MAX_STRUCT_DEPTH int = 6 // The maximum depthness at which we analyse structs (needs to have a bound cause it could be infinite)

// Parses a function declaration "decl" and update counter to reflects what "decl" uses in terms of concurrency primitives
func AnalyseConcurrencyPrimitives(pack_name string, decl *ast.FuncDecl, counter Counter, fileset *token.FileSet, ast_map map[string]*packages.Package) Counter {
	ast.Inspect(decl.Body, func(stmt ast.Node) bool {
		switch stmt := stmt.(type) {
		case *ast.AssignStmt:
			if stmt.Tok == token.DEFINE {
				for _, e := range stmt.Lhs {
					counter = analyseLhs(pack_name, e, counter, fileset, ast_map)
				}
			}
		case *ast.GenDecl:
			for _, spec := range stmt.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						counter = analyseLhs(pack_name, name, counter, fileset, ast_map)
					}
				}
			}
		case *ast.CallExpr:
			counter = analyseCallExpr(pack_name, stmt, counter, fileset, ast_map)
		}
		return true
	})

	return counter
}

func analyseLhs(pack_name string, expr ast.Expr, counter Counter, fileset *token.FileSet, ast_map map[string]*packages.Package) Counter {
	switch typ := removePointer(ast_map[pack_name].TypesInfo.TypeOf(expr)).(type) {
	case nil:
		fmt.Println("couldn't find type of ", expr)
	case *types.Named:
		feature := Feature{
			F_filename:     fileset.Position(expr.Pos()).Filename,
			F_package_name: pack_name,
			F_line_num:     fileset.Position(expr.Pos()).Line}
		if typ.String() == "sync.Mutex" {
			feature.F_type = MUTEX
			counter.Features = append(counter.Features, &feature)
			counter.Mutex_count = counter.Mutex_count + 1
		} else if typ.String() == "sync.WaitGroup" {
			feature.F_type = WAITGROUP
			counter.Features = append(counter.Features, &feature)
			counter.Waitgroup_count = counter.Waitgroup_count + 1
		} else {
			// analyse if the underlyings of the struct contains one
			counter = analyseUnderlying(pack_name, expr, typ.Underlying(), MAX_STRUCT_DEPTH, counter, fileset, ast_map)
		}
	}

	return counter
}
func analyseUnderlying(pack_name string, expr ast.Expr, typ types.Type, depth int, counter Counter, fileset *token.FileSet, ast_map map[string]*packages.Package) Counter {
	if depth > 0 {
		switch typ := removePointer(typ).(type) {
		case nil:
			fmt.Println("couldn't find type of ", typ)
		case *types.Named:
			feature := Feature{
				F_filename:     fileset.Position(expr.Pos()).Filename,
				F_package_name: pack_name,
				F_line_num:     fileset.Position(expr.Pos()).Line}
			if typ.String() == "sync.Mutex" {
				feature.F_type = MUTEX
				counter.Features = append(counter.Features, &feature)
				counter.Mutex_count = counter.Mutex_count + 1
			} else if typ.String() == "sync.WaitGroup" {
				feature.F_type = WAITGROUP
				counter.Features = append(counter.Features, &feature)
				counter.Waitgroup_count = counter.Waitgroup_count + 1
			} else {

				// analyse if the underlyings of the struct contains one
				counter = analyseUnderlying(pack_name, expr, typ.Underlying(), depth-1, counter, fileset, ast_map)
			}

		case *types.Struct:

			for i := 0; i < typ.NumFields(); i++ {
				counter = analyseUnderlying(pack_name, expr, typ.Field(i).Type(), depth-1, counter, fileset, ast_map)
			}
		}
	}

	return counter
}

func analyseCallExpr(pack_name string, call_expr *ast.CallExpr, counter Counter, fileset *token.FileSet, ast_map map[string]*packages.Package) Counter {
	switch expr := call_expr.Fun.(type) {
	case *ast.SelectorExpr:
		if expr.Sel.Name == "Unlock" || expr.Sel.Name == "Lock" {
			switch typ := removePointer(ast_map[pack_name].TypesInfo.TypeOf(expr.X)).(type) {
			case *types.Named:
				if typ.String() == "sync.Mutex" {
					feature := Feature{
						F_filename:     fileset.Position(call_expr.Pos()).Filename,
						F_package_name: pack_name,
						F_line_num:     fileset.Position(call_expr.Pos()).Line}
					if expr.Sel.Name == "Unlock" {
						feature.F_type = UNLOCK
						counter.Features = append(counter.Features, &feature)
						counter.Unlock_count = counter.Unlock_count + 1
					} else {
						feature.F_type = LOCK
						counter.Features = append(counter.Features, &feature)
						counter.Lock_count = counter.Lock_count + 1
					}
				}
			}
		}

		if expr.Sel.Name == "Add" && len(call_expr.Args) == 1 {
			switch typ := removePointer(ast_map[pack_name].TypesInfo.TypeOf(expr.X)).(type) {
			case *types.Named:
				if typ.String() == "sync.WaitGroup" {

					// Look at right hand side if it is a const or not
					if isConstant, val := isConst(call_expr.Args[0], ast_map[pack_name]); isConstant {

						feature := Feature{
							F_filename:     fileset.Position(call_expr.Pos()).Filename,
							F_package_name: pack_name,
							F_line_num:     fileset.Position(call_expr.Pos()).Line}
						feature.F_type = KNOWN_ADD
						feature.F_number = fmt.Sprint(call_expr.Args[0]) + " val is : " + strconv.Itoa(val)
						counter.Known_add_count = counter.Known_add_count + 1
						counter.Features = append(counter.Features, &feature)
					} else {
						feature := Feature{
							F_filename:     fileset.Position(call_expr.Pos()).Filename,
							F_package_name: pack_name,
							F_line_num:     fileset.Position(call_expr.Pos()).Line}
						feature.F_type = UNKNOWN_ADD
						counter.Features = append(counter.Features, &feature)
						counter.Unknown_add_count = counter.Unknown_add_count + 1
					}
				}
			}
		}

		if expr.Sel.Name == "Done" {
			switch typ := removePointer(ast_map[pack_name].TypesInfo.TypeOf(expr.X)).(type) {
			case *types.Named:
				if typ.String() == "sync.WaitGroup" {
					feature := Feature{
						F_filename:     fileset.Position(call_expr.Pos()).Filename,
						F_package_name: pack_name,
						F_line_num:     fileset.Position(call_expr.Pos()).Line}
					feature.F_type = DONE
					counter.Features = append(counter.Features, &feature)
					counter.Done_count = counter.Done_count + 1
				}
			}
		}
	}

	return counter
}

func isConst(expr ast.Expr, pack *packages.Package) (found bool, val int) {
	switch expr := expr.(type) {
	case *ast.Ident:
		obj := expr.Obj
		if obj != nil {
			if obj.Kind == ast.Con {
				switch value_spec := obj.Decl.(type) {
				case *ast.ValueSpec:
					if value_spec.Values != nil && len(value_spec.Values) > 0 {
						switch val := value_spec.Values[0].(type) {
						case *ast.BasicLit:
							v, err := strconv.Atoi(val.Value)
							if err == nil {
								return true, v
							}
						case *ast.Ident:
							return isConst(val, pack)
						}
					}
				}
			}
		}
	case *ast.SelectorExpr:
		obj := expr.Sel.Obj
		if obj != nil {
			if obj.Kind == ast.Con {
				switch value_spec := obj.Decl.(type) {
				case *ast.ValueSpec:
					if value_spec.Values != nil && len(value_spec.Values) > 0 {
						switch val := value_spec.Values[0].(type) {
						case *ast.BasicLit:
							v, err := strconv.Atoi(val.Value)
							if err == nil {
								return true, v
							}
						case *ast.Ident:
							return isConst(val, pack)
						}
					}
				}
			}
		}
	case *ast.BasicLit:
		if expr.Kind == token.INT {
			val, err := strconv.Atoi(expr.Value)
			if err == nil {
				return true, val
			}
		}
	}
	return false, -1
}

func removePointer(typ types.Type) types.Type {
	switch typ := typ.(type) {
	case *types.Pointer:
		return removePointer(typ.Elem())
	default:
		return typ
	}
}
