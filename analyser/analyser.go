package main

import (
	"go/ast"
	"go/token"
	"strconv"
)

func AnalyseAst(fileset *token.FileSet, package_name string, filename string, node ast.Node, channel chan Counter, name string) {
	var counter Counter = Counter{Go_count: 0, Send_count: 0, Rcv_count: 0, Chan_count: 0, filename: name}
	var env []string = []string{}

	switch file := node.(type) {
	case *ast.File:
		addGlobalVarToEnv(file, &env)
		for _, decl := range file.Decls {
			fresh_env := env
			ast.Inspect(decl, func(decl ast.Node) bool {
				analyseNode(fileset, package_name, filename, decl, &counter, &fresh_env)
				return true
			})
		}
	}

	setFeaturesNumber(&counter)
	channel <- counter
}

func analyseNode(fileset *token.FileSet, package_name string, filename string, node ast.Node, counter *Counter, env *[]string) {

	var feature Feature = Feature{
		F_filename:     filename,
		F_package_name: package_name,
		F_type:         NONE}

	switch x := node.(type) {
	// Add General decleration
	case *ast.GenDecl:
		if x.Tok == token.VAR {
			for _, spec := range x.Specs {
				switch value_spec := spec.(type) {
				case *ast.ValueSpec:
					for index, value := range value_spec.Values {
						switch call_expr := value.(type) {
						case *ast.CallExpr:
							switch ident := call_expr.Fun.(type) {
							case *ast.Ident:
								if ident.Name == "make" {
									if len(call_expr.Args) > 0 {
										switch call_expr.Args[0].(type) {
										case *ast.ChanType:
											ident1 := value_spec.Names[index]
											*env = append(*env, ident1.Name)
											checkDepthChan(call_expr, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
										}
									}
								}
							}
						case *ast.CompositeLit:
							switch array_type := call_expr.Type.(type) {
							case *ast.Ident:
								// Possible assignment of a struct struct = Struct{bla:0, bla1}
								for _, elt := range call_expr.Elts {
									switch valueExp := elt.(type) {
									case *ast.KeyValueExpr:
										switch ident := valueExp.Key.(type) {
										case *ast.Ident:
											switch call := valueExp.Value.(type) {
											case *ast.CallExpr:
												ident1 := value_spec.Names[index]
												checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
											}
										}
									}
								}

							case *ast.ArrayType:
								checkArrayType(array_type, counter, feature, fileset, 1)
							case *ast.MapType:
								chan_in_map := false
								// we have a declaration of a map
								switch array_type.Key.(type) {
								case *ast.ChanType:
									chan_in_map = true
								}
								switch array_type.Value.(type) {
								case *ast.ChanType:
									chan_in_map = true
								}

								if chan_in_map {
									chan_feature := feature
									chan_feature.F_line_num = fileset.Position(x.Pos()).Line
									chan_feature.F_type = CHAN_MAP
									counter.Chan_map_count++
									counter.Features = append(counter.Features, &chan_feature)
								}
							}

						case *ast.UnaryExpr:
							switch expr := call_expr.X.(type) {
							case *ast.CompositeLit:
								switch array_type := expr.Type.(type) {
								case *ast.Ident:
									// Possible assignment of a struct struct = Struct{bla:0, bla1}
									for _, elt := range expr.Elts {
										switch valueExp := elt.(type) {
										case *ast.KeyValueExpr:
											switch ident := valueExp.Key.(type) {
											case *ast.Ident:
												switch call := valueExp.Value.(type) {
												case *ast.CallExpr:
													ident1 := value_spec.Names[index]
													checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
												}
											}
										}
									}

								case *ast.ArrayType:
									checkArrayType(array_type, counter, feature, fileset, 1)
								case *ast.MapType:
									chan_in_map := false
									// we have a declaration of a map
									switch array_type.Key.(type) {
									case *ast.ChanType:
										chan_in_map = true
									}
									switch array_type.Value.(type) {
									case *ast.ChanType:
										chan_in_map = true
									}

									if chan_in_map {
										chan_feature := feature
										chan_feature.F_line_num = fileset.Position(x.Pos()).Line
										chan_feature.F_type = CHAN_MAP
										counter.Chan_map_count++
										counter.Features = append(counter.Features, &chan_feature)
									}
								}
							}
						}
					}
				}
			}
		}
	case *ast.GoStmt:
		go_feature := Feature{
			F_filename:     feature.F_filename,
			F_package_name: feature.F_package_name,
			F_line_num:     fileset.Position(x.Pos()).Line}
		go_feature.F_type = GOROUTINE
		counter.Go_count++
		counter.Features = append(counter.Features, &go_feature)
	case *ast.SendStmt:
		send_feature := Feature{
			F_filename:     feature.F_filename,
			F_package_name: feature.F_package_name,
			F_line_num:     fileset.Position(x.Pos()).Line}
		send_feature.F_type = SEND
		counter.Send_count++
		counter.Features = append(counter.Features, &send_feature)
	case *ast.UnaryExpr:
		if x.Op.String() == "<-" {
			send_feature := Feature{
				F_filename:     feature.F_filename,
				F_package_name: feature.F_package_name,
				F_line_num:     fileset.Position(x.Pos()).Line}
			send_feature.F_type = RECEIVE
			counter.Rcv_count++
			counter.Features = append(counter.Features, &send_feature)
		}
	case *ast.AssignStmt:
		// look for a make(chan X) or a make(chan X,n)
		for index, rh := range x.Rhs {
			switch call_expr := rh.(type) {
			case *ast.CallExpr:
				switch ident := x.Lhs[index].(type) {
				case *ast.Ident:
					checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, true)
				case *ast.SelectorExpr:
					if ident.X != nil && ident.Sel != nil {
						switch name := ident.X.(type) {
						case *ast.Ident:
							checkDepthChan(call_expr, feature, env, counter, ident.Sel.Name+"."+name.Name, fileset, true)
						}
					}
				}

			case *ast.CompositeLit:
				switch array_type := call_expr.Type.(type) {
				case *ast.Ident:
					// Possible assignment of a struct struct = Struct{bla:0, bla1}
					for _, elt := range call_expr.Elts {
						switch valueExp := elt.(type) {
						case *ast.KeyValueExpr:
							switch ident := valueExp.Key.(type) {
							case *ast.Ident:
								switch call := valueExp.Value.(type) {
								case *ast.CallExpr:
									switch ident1 := x.Lhs[index].(type) {
									case *ast.Ident:
										checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
									}
								}
							}
						}
					}

				case *ast.ArrayType:
					checkArrayType(array_type, counter, feature, fileset, 1)
				case *ast.MapType:
					chan_in_map := false
					// we have a declaration of a map
					switch array_type.Key.(type) {
					case *ast.ChanType:
						chan_in_map = true
					}
					switch array_type.Value.(type) {
					case *ast.ChanType:
						chan_in_map = true
					}

					if chan_in_map {
						chan_feature := feature
						chan_feature.F_line_num = fileset.Position(x.Pos()).Line
						chan_feature.F_type = CHAN_MAP
						counter.Chan_map_count++
						counter.Features = append(counter.Features, &chan_feature)
					}
				}

			case *ast.UnaryExpr:
				switch call_expr := call_expr.X.(type) {
				case *ast.CompositeLit:
					switch array_type := call_expr.Type.(type) {
					case *ast.Ident:
						// Possible assignment of a struct struct = Struct{bla:0, bla1}
						for _, elt := range call_expr.Elts {
							switch valueExp := elt.(type) {
							case *ast.KeyValueExpr:
								switch ident := valueExp.Key.(type) {
								case *ast.Ident:
									switch call := valueExp.Value.(type) {
									case *ast.CallExpr:
										switch ident1 := x.Lhs[index].(type) {
										case *ast.Ident:
											checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
										}
									}
								}
							}
						}

					case *ast.ArrayType:
						checkArrayType(array_type, counter, feature, fileset, 1)
					case *ast.MapType:
						chan_in_map := false
						// we have a declaration of a map
						switch array_type.Key.(type) {
						case *ast.ChanType:
							chan_in_map = true
						}
						switch array_type.Value.(type) {
						case *ast.ChanType:
							chan_in_map = true
						}

						if chan_in_map {
							chan_feature := feature
							chan_feature.F_line_num = fileset.Position(x.Pos()).Line
							chan_feature.F_type = CHAN_MAP
							counter.Chan_map_count++
							counter.Features = append(counter.Features, &chan_feature)
						}
					}
				}
			}
		}
	case *ast.DeclStmt:
		// look for a make(chan X) or a make(chan X,n)
		switch decl := x.Decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.VAR {
				for _, spec := range decl.Specs {
					switch value := spec.(type) {
					case *ast.ValueSpec:
						switch value_type := value.Type.(type) {
						case *ast.Ident:
							// 	looking for a declaration of a struct
							for index, exp := range value.Values {
								switch composite := exp.(type) {
								case *ast.CompositeLit:
									for _, elt := range composite.Elts {
										switch valueExp := elt.(type) {
										case *ast.KeyValueExpr:
											switch ident := valueExp.Key.(type) {
											case *ast.Ident:
												switch call := valueExp.Value.(type) {
												case *ast.CallExpr:
													checkDepthChan(call, feature, env, counter, value.Names[index].Name+"."+ident.Name, fileset, true)
												}
											}
										}
									}
								}
							}
						case *ast.ArrayType:
							// we have a declaration of an array
							num_of_arrays := len(value.Names)
							checkArrayType(value_type, counter, feature, fileset, num_of_arrays)

						case *ast.MapType:
							chan_in_map := false
							// we have a declaration of a map
							switch value_type.Key.(type) {
							case *ast.ChanType:
								chan_in_map = true
							}
							switch value_type.Value.(type) {
							case *ast.ChanType:
								chan_in_map = true
							}

							if chan_in_map {
								chan_feature := feature
								chan_feature.F_type = CHAN_MAP
								counter.Chan_map_count++
								counter.Features = append(counter.Features, &chan_feature)
							}
						}

						for index, val := range value.Values {
							switch call_expr := val.(type) {
							case *ast.CallExpr:
								checkDepthChan(call_expr, feature, env, counter, value.Names[index].Name, fileset, true)
							}
						}
					}
				}
			}
		}
	case *ast.ForStmt:
		makeChanInFor(x, feature, env, counter, fileset)
		// look in the block and see if goroutine are created in a for loop
		for _, stmt := range x.Body.List {
			switch x_node := stmt.(type) {
			case *ast.GoStmt:
				go_feature := feature
				go_feature.F_type = GO_IN_FOR
				counter.Go_in_for_count++
				go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
				counter.Features = append(counter.Features, &go_feature)
				switch bin_expr := x.Cond.(type) {
				case *ast.BinaryExpr:
					if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
						// check if the right hand side is a constant
						val, isCons := isConstant(bin_expr.Y)
						if isCons {
							go_feature := feature
							go_feature.F_type = GO_IN_CONSTANT_FOR
							go_feature.F_number = strconv.Itoa(val)
							go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
							counter.Go_in_constant_for_count++
							counter.Features = append(counter.Features, &go_feature)
						}
					} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
						// check if the initialisation is a constant
						switch assign := x.Init.(type) {
						case *ast.AssignStmt:
							for _, rh := range assign.Rhs {
								val, isCons := isConstant(rh)
								if isCons {
									go_feature := feature
									go_feature.F_type = GO_IN_CONSTANT_FOR
									go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
									go_feature.F_number = strconv.Itoa(val)
									counter.Go_in_constant_for_count++
									counter.Features = append(counter.Features, &go_feature)
								}
							}
						}
					}
				}
			}
		}

	case *ast.RangeStmt:
		// check if the stmt is a range over a channel

		if x.Key != nil {
			switch ident1 := x.Key.(type) {

			case *ast.Ident:
				if ident1.Obj != nil {
					switch assign := ident1.Obj.Decl.(type) {
					case *ast.AssignStmt:

						for _, rh := range assign.Rhs {
							switch unary := rh.(type) {
							case *ast.UnaryExpr:
								if unary.Op == token.RANGE {

									switch chan_type := unary.X.(type) {
									case *ast.Ident:
										// trying to range over a channel
										if chan_type.Obj != nil {
											found, _ := isChan(unary.X, env)
											if found {
												range_feature := feature
												range_feature.F_type = RANGE_OVER_CHAN
												range_feature.F_line_num = fileset.Position(unary.Pos()).Line
												counter.Range_over_chan_count++
												counter.Features = append(counter.Features, &range_feature)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			switch ident1 := x.X.(type) {
			case *ast.Ident:
				if ident1.Obj != nil {
					found, _ := isChan(ident1, env)
					if found {
						range_feature := feature
						range_feature.F_type = RANGE_OVER_CHAN
						range_feature.F_line_num = fileset.Position(ident1.Pos()).Line
						counter.Range_over_chan_count++
						counter.Features = append(counter.Features, &range_feature)
					}
				}
			}
		}
		if x.Body != nil {
			for _, stmt := range x.Body.List {
				switch assign_stmt := stmt.(type) {
				case *ast.GoStmt:
					go_in_for := Feature{
						F_filename:     filename,
						F_package_name: package_name,
						F_type:         GO_IN_FOR}
					counter.Go_in_for_count++
					go_in_for.F_line_num = fileset.Position(assign_stmt.Pos()).Line
					counter.Features = append(counter.Features, &go_in_for)

				case *ast.AssignStmt:
					for _, expr := range assign_stmt.Rhs {
						found, chan_name := isChan(expr, env)
						if found {
							assign_chan_in_for := Feature{
								F_type:         ASSIGN_CHAN_IN_FOR,
								F_filename:     filename,
								F_package_name: package_name}
							counter.Assign_chan_in_for_count++
							assign_chan_in_for.F_number = chan_name
							assign_chan_in_for.F_line_num = fileset.Position(expr.Pos()).Line
							counter.Features = append(counter.Features, &assign_chan_in_for)
						}
					}
				}
			}
		}
		makeChanInRange(x, feature, env, counter, fileset)

	case *ast.ExprStmt:
		// looking for a close
		switch call_expr := x.X.(type) {
		case *ast.CallExpr:
			switch ident := call_expr.Fun.(type) {
			case *ast.Ident:
				if ident.Name == "close" && len(call_expr.Args) == 1 {
					// we have a close
					found, _ := isChan(call_expr.Args[0], env)
					if found {
						// we have a close on a chan
						close_feature := feature
						counter.Close_chan_count++
						close_feature.F_type = CLOSE_CHAN
						close_feature.F_line_num = fileset.Position(ident.Pos()).Line
						counter.Features = append(counter.Features, &close_feature)
					}
				}
			}
		}
	case *ast.SelectStmt:
		if x.Body != nil {
			var with_default bool = false
			for _, stmt := range x.Body.List {
				switch comm := stmt.(type) {
				case *ast.CommClause:
					if comm.Comm == nil {
						// we have a select with a default

						with_default = true
					}
				}
			}
			select_feature := feature
			if with_default {
				select_feature.F_type = DEFAULT_SELECT
				counter.Default_select_count++
			} else {
				select_feature.F_type = SELECT
				counter.Select_count++
			}
			select_feature.F_number = strconv.Itoa(len(x.Body.List))
			select_feature.F_line_num = fileset.Position(x.Pos()).Line
			counter.Features = append(counter.Features, &select_feature)
		}
	case *ast.DeferStmt:
		if x.Call != nil {
			call_expr := x.Call
			switch ident := call_expr.Fun.(type) {
			case *ast.Ident:
				if ident.Name == "close" && len(call_expr.Args) == 1 {
					found, _ := isChan(call_expr.Args[0], env)
					if found {
						// we have a close on a chan
						close_feature := feature
						counter.Close_chan_count++
						close_feature.F_type = CLOSE_CHAN
						close_feature.F_line_num = fileset.Position(call_expr.Pos()).Line
						counter.Features = append(counter.Features, &close_feature)
					}
				}
			}
		}
	case *ast.FuncDecl:
		// look for a <-chan, chan<- or chan as function fields
		for _, field := range x.Type.Params.List {
			switch chan_type := field.Type.(type) {
			case *ast.ChanType:
				switch chan_type.Dir {
				case ast.RECV:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = RECEIVE_CHAN
					counter.Receive_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				case ast.SEND:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = SEND_CHAN
					counter.Send_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				default:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = PARAM_CHAN
					counter.Param_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				}
			}
		}
	}
}

func makeChanInFor(forStmt *ast.ForStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet) {
	for _, block := range forStmt.Body.List {
		switch stmt := block.(type) {
		case *ast.AssignStmt:
			// chan in for
			for index, rh := range stmt.Rhs {
				switch call_expr := rh.(type) {
				case *ast.CallExpr:
					switch ident := stmt.Lhs[index].(type) {
					case *ast.Ident:
						if checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, false) {
							switch bin_expr := forStmt.Cond.(type) {
							case *ast.BinaryExpr:
								if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
									// check if the right hand side is a constant
									val, isCons := isConstant(bin_expr.Y)
									if isCons {
										make_chan_in_for := Feature{}
										make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
										if stmt.Tok == token.DEFINE {
											make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
										}
										make_chan_in_for.F_filename = feature.F_filename
										make_chan_in_for.F_package_name = feature.F_package_name
										make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
										make_chan_in_for.F_number = strconv.Itoa(val)
										counter.Make_chan_in_constant_for_count++
										counter.Features = append(counter.Features, &make_chan_in_for)
									}
									// }
								} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
									// check if the initialisation is a constant
									switch assign := forStmt.Init.(type) {
									case *ast.AssignStmt:
										for _, rh := range assign.Rhs {
											val, isCons := isConstant(rh)
											if isCons {
												make_chan_in_for := Feature{}
												make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
												if stmt.Tok == token.DEFINE {
													make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
												}
												make_chan_in_for.F_filename = feature.F_filename
												make_chan_in_for.F_package_name = feature.F_package_name
												make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
												make_chan_in_for.F_number = strconv.Itoa(val)
												counter.Make_chan_in_constant_for_count++
												counter.Features = append(counter.Features, &make_chan_in_for)
											}
										}
									}
								} else {
									make_chan_in_for := Feature{}
									make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
									if stmt.Tok == token.DEFINE {
										make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
									}
									make_chan_in_for.F_filename = feature.F_filename
									make_chan_in_for.F_package_name = feature.F_package_name
									make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
									counter.Make_chan_in_constant_for_count++
									counter.Features = append(counter.Features, &make_chan_in_for)
								}
							default:
								make_chan_in_for := Feature{}
								make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
								if stmt.Tok == token.DEFINE {
									make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
								}
								make_chan_in_for.F_filename = feature.F_filename
								make_chan_in_for.F_package_name = feature.F_package_name
								make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
								counter.Make_chan_in_constant_for_count++
								counter.Features = append(counter.Features, &make_chan_in_for)
							}
						}
					}

				}
			}

		case *ast.DeclStmt: // is the declaration in a constant or not for loop ?
			if chanDecleration(stmt, feature, env, counter, fileset, false) {
				switch bin_expr := forStmt.Cond.(type) {
				case *ast.BinaryExpr:
					if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
						// check if the right hand side is a constant
						val, isCons := isConstant(bin_expr.Y)
						if isCons {
							make_chan_in_for := Feature{}
							make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
							make_chan_in_for.F_filename = feature.F_filename
							make_chan_in_for.F_package_name = feature.F_package_name
							make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
							make_chan_in_for.F_number = strconv.Itoa(val)
							counter.Make_chan_in_constant_for_count++
							counter.Features = append(counter.Features, &make_chan_in_for)
						}
					} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
						// check if the initialisation is a constant
						switch assign := forStmt.Init.(type) {
						case *ast.AssignStmt:
							for _, rh := range assign.Rhs {
								val, isCons := isConstant(rh)
								if isCons {
									make_chan_in_for := Feature{}
									make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
									make_chan_in_for.F_filename = feature.F_filename
									make_chan_in_for.F_package_name = feature.F_package_name
									make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
									make_chan_in_for.F_number = strconv.Itoa(val)
									counter.Make_chan_in_constant_for_count++
									counter.Features = append(counter.Features, &make_chan_in_for)
								}
							}
						}
					} else {
						make_chan_in_for := Feature{}
						make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
						make_chan_in_for.F_filename = feature.F_filename
						make_chan_in_for.F_package_name = feature.F_package_name
						make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
						counter.Make_chan_in_constant_for_count++
						counter.Features = append(counter.Features, &make_chan_in_for)
					}
				default:
					make_chan_in_for := Feature{}
					make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
					make_chan_in_for.F_filename = feature.F_filename
					make_chan_in_for.F_package_name = feature.F_package_name
					make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
					counter.Make_chan_in_constant_for_count++
					counter.Features = append(counter.Features, &make_chan_in_for)
				}
			}
		}
	}

	for _, stmt := range forStmt.Body.List {
		switch x_node := stmt.(type) {

		case *ast.AssignStmt:
			for _, expr := range x_node.Rhs {
				found, chan_name := isChan(expr, env)
				chan_feature := feature
				if found {
					chan_feature.F_type = ASSIGN_CHAN_IN_FOR
					if x_node.Tok == token.DEFINE {
						chan_feature.F_type = MAKE_CHAN_IN_CONSTANT_FOR
					}
					chan_feature.F_line_num = fileset.Position(expr.Pos()).Line
					chan_feature.F_number = chan_name
					counter.Assign_chan_in_for_count++
					counter.Features = append(counter.Features, &chan_feature)
				}
			}
		}
	}
}

func makeChanInRange(rangeStmt *ast.RangeStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet) {

	for _, block := range rangeStmt.Body.List {
		switch stmt := block.(type) {
		case *ast.AssignStmt:
			// chan in for
			for index, rh := range stmt.Rhs {
				switch call_expr := rh.(type) {
				case *ast.CallExpr:
					switch ident := stmt.Lhs[index].(type) {
					case *ast.Ident:
						if checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, false) {
							make_chan_in_for := Feature{}
							make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
							make_chan_in_for.F_filename = feature.F_filename
							make_chan_in_for.F_package_name = feature.F_package_name
							make_chan_in_for.F_line_num = fileset.Position(call_expr.Pos()).Line
							counter.Make_chan_in_constant_for_count++
							counter.Features = append(counter.Features, &make_chan_in_for)
						}
					}
				}
			}

		case *ast.DeclStmt:

			if chanDecleration(stmt, feature, env, counter, fileset, false) {
				make_chan_in_for := Feature{}
				make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
				make_chan_in_for.F_filename = feature.F_filename
				make_chan_in_for.F_package_name = feature.F_package_name
				make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
				counter.Make_chan_in_constant_for_count++
				counter.Features = append(counter.Features, &make_chan_in_for)
			}
		}
	}
}

func chanDecleration(stmt *ast.DeclStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet, add bool) bool {
	var found_decl bool = false
	switch decl := stmt.Decl.(type) {
	case *ast.GenDecl:
		if decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				switch value := spec.(type) {
				case *ast.ValueSpec:
					switch value.Type.(type) {
					case *ast.ChanType:
						// we have a var x chan X
						if value.Values != nil {
							if len(value.Values) == len(value.Names) {
								for index, val := range value.Values {
									switch call_expr := val.(type) {
									case *ast.CallExpr:
										found_decl = checkDepthChan(call_expr, feature, env, counter, value.Names[index].Name, fileset, add)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return found_decl
}

func checkArrayType(array_type *ast.ArrayType, counter *Counter, feature Feature, fileset *token.FileSet, num_of_arrays int) {
	switch chan_type := array_type.Elt.(type) {
	case *ast.ChanType:
		//we have an array of chan
		if array_type.Len != nil {
			// check if constant
			val, isCons := isConstant(array_type.Len)
			if isCons {
				for i := 0; i < num_of_arrays; i++ {
					array_feature := feature
					array_feature.F_type = CONSTANT_CHAN_ARRAY
					array_feature.F_number = strconv.Itoa(val)
					array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					counter.Constant_chan_array_count += num_of_arrays
					counter.Features = append(counter.Features, &array_feature)
				}
			} else {
				for i := 0; i < num_of_arrays; i++ {
					array_feature := feature
					array_feature.F_type = ARRAY_OF_CHANNELS
					counter.Array_of_channels_count += num_of_arrays
					array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					counter.Features = append(counter.Features, &array_feature)
				}
			}
		} else {
			for i := 0; i < num_of_arrays; i++ {
				array_feature := feature
				array_feature.F_type = CHAN_SLICE
				counter.Chan_slice_count += num_of_arrays
				array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
				counter.Features = append(counter.Features, &array_feature)
			}
		}
	}
}

func checkDepthChan(call_expr *ast.CallExpr, feature Feature, env *[]string, counter *Counter, chan_name string, fileset *token.FileSet, add bool) bool {
	var chan_found bool = false
	switch ident := call_expr.Fun.(type) {
	case *ast.Ident:
		if ident.Name == "make" {
			if len(call_expr.Args) > 0 {
				switch chan_type := call_expr.Args[0].(type) {
				case *ast.ChanType:
					chan_found = true
					*env = append(*env, chan_name)
					switch chan_type.Value.(type) {

					case *ast.ChanType:
						chan_feature := Feature{
							F_filename:     feature.F_filename,
							F_package_name: feature.F_package_name,
							F_line_num:     fileset.Position(call_expr.Pos()).Line}
						chan_feature.F_type = CHAN_OF_CHANS
						chan_feature.F_number = chan_name
						counter.Chan_of_chans_count++
						counter.Features = append(counter.Features, &chan_feature)
					default:
						if len(call_expr.Args) > 1 {
							val, isCons := isConstant(call_expr.Args[1])
							if isCons {
								if add {
									if val != 0 {
										chan_feature := Feature{
											F_filename:     feature.F_filename,
											F_package_name: feature.F_package_name,
											F_line_num:     fileset.Position(call_expr.Pos()).Line}
										chan_feature.F_type = KNOWN_CHAN_DEPTH
										chan_feature.F_number = strconv.Itoa(val)
										counter.Known_chan_depth_count++
										counter.Features = append(counter.Features, &chan_feature)
									} else {
										chan_feature := Feature{
											F_filename:     feature.F_filename,
											F_package_name: feature.F_package_name,
											F_line_num:     fileset.Position(call_expr.Pos()).Line}
										chan_feature.F_type = MAKE_CHAN
										counter.Sync_Chan_count++
										counter.Features = append(counter.Features, &chan_feature)
									}
								}
							} else {
								if add {
									chan_feature := Feature{
										F_filename:     feature.F_filename,
										F_package_name: feature.F_package_name,
										F_line_num:     fileset.Position(call_expr.Pos()).Line}
									chan_feature.F_type = UNKNOWN_CHAN_DEPTH //unknown depth
									counter.Unknown_chan_depth_count++
									counter.Features = append(counter.Features, &chan_feature)
								}
							}
						} else {
							if add {
								chan_feature := Feature{
									F_filename:     feature.F_filename,
									F_package_name: feature.F_package_name,
									F_line_num:     fileset.Position(call_expr.Pos()).Line}
								chan_feature.F_type = MAKE_CHAN
								counter.Sync_Chan_count++
								counter.Features = append(counter.Features, &chan_feature)
							}
						}
					}
				}
			}
		}
	}

	return chan_found
}

func isConstant(node ast.Node) (int, bool) {
	var isCons bool = false
	var value int = 0
	switch ident := node.(type) {
	case *ast.Ident:
		if ident.Obj != nil {
			if ident.Obj.Kind == ast.Con {
				switch value_spec := ident.Obj.Decl.(type) {
				case *ast.ValueSpec:

					if value_spec.Values != nil && len(value_spec.Values) > 0 {
						switch val := value_spec.Values[0].(type) {
						case *ast.BasicLit:
							parsed_int, _ := strconv.Atoi(val.Value)
							value = int(parsed_int)
							isCons = true
						case *ast.Ident:
							value, isCons = isConstant(val)
						}
					}
				}
			}
		}
	case *ast.BasicLit:
		if ident.Kind == token.INT {
			isCons = true
			parsed_int, _ := strconv.Atoi(ident.Value)
			value = int(parsed_int)
		}
	default:
		isCons = false
	}

	return value, isCons
}

func isChan(node interface{}, env *[]string) (bool, string) {

	chan_name := ""
	switch make_chan := node.(type) {
	case *ast.AssignStmt:
		var chan_found bool = false
		ast.Inspect(make_chan, func(x_node ast.Node) bool {
			switch x_node.(type) {
			case *ast.ChanType:
				chan_found = true

				return false
			}
			return true
		})

		if !chan_found {
			for _, rh := range make_chan.Rhs {
				switch ident := rh.(type) {
				case *ast.Ident:
					for _, name := range *env {
						if name == ident.Name {
							chan_found = true
							chan_name = name
							break
						}
					}
				}
			}
		}
		return chan_found, chan_name
	case *ast.Ident:
		for _, name := range *env {
			if name == make_chan.Name {
				chan_name = name
				return true, chan_name
			}
		}
	}

	return false, chan_name
}

func addGlobalVarToEnv(file *ast.File, env *[]string) {
	for _, decl := range file.Decls {
		switch genDecl := decl.(type) {
		case *ast.GenDecl:
			if genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					switch value_spec := spec.(type) {
					case *ast.ValueSpec:
						for index, value := range value_spec.Values {
							switch call_expr := value.(type) {
							case *ast.CallExpr:
								switch ident := call_expr.Fun.(type) {
								case *ast.Ident:
									if ident.Name == "make" {
										if len(call_expr.Args) > 0 {
											switch call_expr.Args[0].(type) {
											case *ast.ChanType:
												*env = append(*env, value_spec.Names[index].Name)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
