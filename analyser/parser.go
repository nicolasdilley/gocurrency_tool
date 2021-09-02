package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// parse a particular dir
func ParseDir(proj_name string, path_to_dir string, path_to_main_dir string) PackageCounter {

	var fileSet *token.FileSet = token.NewFileSet()
	var counter PackageCounter = PackageCounter{
		Counter: Counter{
			Go_count:     0,
			Send_count:   0,
			Rcv_count:    0,
			Chan_count:   0,
			IsPackage:    true,
			Project_name: proj_name},
		File_counters: []*Counter{}}

	f, err := parser.ParseDir(fileSet, path_to_dir, nil, parser.AllErrors)

	if proj_name == "test" {
		ast.Print(fileSet, f)
	}
	if err != nil {
		fmt.Printf("An error was found in package %s : %v", filepath.Base(path_to_dir), err)
	}

	if len(f) == 0 {
		return counter
	}

	for pack_name, pack := range f {

		var package_counter_chan chan Counter = make(chan Counter)
		counter.Counter.Package_name = strings.TrimPrefix(strings.TrimPrefix(path_to_dir, path_to_main_dir)+"/"+pack_name, "/")
		counter.Counter.Package_path = path_to_dir
		// Analyse each file
		for name, file := range pack.Files {
			filename := strings.TrimPrefix(strings.TrimPrefix(path_to_dir, path_to_main_dir)+"/"+filepath.Base(name), "/")
			go AnalyseAst(fileSet, pack_name, filename, file, package_counter_chan, name) // launch a goroutine for each file
		}

		// Receive the results of the analysis of each file
		for range pack.Files {

			var new_counter Counter = <-package_counter_chan

			new_counter.IsPackage = false
			new_counter.Project_name = proj_name
			if len(new_counter.Features) > 0 {
				new_counter.Has_feature = true
			}
			counter.Counter.Go_count += new_counter.Go_count
			counter.Counter.Send_count += new_counter.Send_count
			counter.Counter.Rcv_count += new_counter.Rcv_count
			counter.Counter.Chan_count += new_counter.Chan_count
			counter.Counter.Go_in_for_count += new_counter.Go_in_for_count
			counter.Counter.Range_over_chan_count += new_counter.Range_over_chan_count
			counter.Counter.Go_in_constant_for_count += new_counter.Go_in_constant_for_count
			counter.Counter.Array_of_channels_count += new_counter.Array_of_channels_count
			counter.Counter.Sync_Chan_count += new_counter.Sync_Chan_count
			counter.Counter.Known_chan_depth_count += new_counter.Known_chan_depth_count
			counter.Counter.Unknown_chan_depth_count += new_counter.Unknown_chan_depth_count
			counter.Counter.Make_chan_in_for_count += new_counter.Make_chan_in_for_count
			counter.Counter.Make_chan_in_constant_for_count += new_counter.Make_chan_in_constant_for_count
			counter.Counter.Constant_chan_array_count += new_counter.Constant_chan_array_count
			counter.Counter.Chan_slice_count += new_counter.Chan_slice_count
			counter.Counter.Chan_map_count += new_counter.Chan_map_count
			counter.Counter.Close_chan_count += new_counter.Close_chan_count
			counter.Counter.Select_count += new_counter.Select_count
			counter.Counter.Default_select_count += new_counter.Default_select_count
			counter.Counter.Assign_chan_in_for_count += new_counter.Assign_chan_in_for_count
			counter.Counter.Chan_of_chans_count += new_counter.Chan_of_chans_count
			counter.Counter.Send_chan_count += new_counter.Send_chan_count
			counter.Counter.Receive_chan_count += new_counter.Receive_chan_count
			counter.Counter.Param_chan_count += new_counter.Param_chan_count

			counter.File_counters = append(counter.File_counters, &new_counter)

		}

	}

	return counter
}

func ParseConcurrencyPrimitives(path_to_dir string, counter Counter) Counter {
	package_names := []string{}

	filepath.Walk(path_to_dir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			if file.Name() != "vendor" && file.Name() != "third_party" {
				path, _ = filepath.Abs(path)
				package_names = append(package_names, path)
			} else {
				return filepath.SkipDir
			}
		}
		return nil
	})

	var ast_map map[string]*packages.Package = make(map[string]*packages.Package)
	var cfg *packages.Config = &packages.Config{Mode: packages.LoadAllSyntax, Fset: &token.FileSet{}, Dir: path_to_dir, Tests: true}

	package_names = append([]string{"."}, package_names...)
	lpkgs, err := packages.Load(cfg, package_names...)

	if err != nil {
		fmt.Println("couldn't load ", path_to_dir)
	}

	for _, pack := range lpkgs {
		ast_map[pack.Name] = pack
	}

	for pack_name, node := range ast_map {
		// Analyse each file

		// make sure the package doesnt contain any global concurrency primitives

		for _, file := range node.Syntax {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.FuncDecl:
					// Analyse each function decleration
					if decl.Body != nil {
						counter = AnalyseConcurrencyPrimitives(pack_name, decl, counter, cfg.Fset, ast_map)
					}
				}
			}
		}

	}

	return counter
}
