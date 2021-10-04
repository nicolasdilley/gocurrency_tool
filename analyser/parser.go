package main

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/nicolasdilley/gocurrency_tool/analyser/analyse"
	"github.com/nicolasdilley/gocurrency_tool/analyser/output"
	"golang.org/x/tools/go/packages"
)

// parse a particular dir
func ParseDir(proj_name string, path_to_main_dir string, fileSet *token.FileSet, ast_map map[string]*packages.Package) []*output.PackageCounter {
	var packages []*output.PackageCounter

	for pack_name, pack := range ast_map {
		var counter output.PackageCounter = output.PackageCounter{
			Counter: analyse.Counter{
				Go_count:     0,
				Send_count:   0,
				Rcv_count:    0,
				Chan_count:   0,
				IsPackage:    true,
				Project_name: proj_name},
			File_counters: []*analyse.Counter{}}
		fmt.Print("Parsing package ", pack_name, " ", len(pack.Syntax), " files ", " : ")

		var package_counter_chan chan analyse.Counter = make(chan analyse.Counter)
		full_pack_name := ""
		if strings.Contains(pack.PkgPath, "projects-gocurrency/") {
			full_pack_name = strings.Join(strings.Split(strings.Split(pack.PkgPath, "projects-gocurrency/")[1], "/")[1:], "/")

		} else {
			if strings.Contains(pack.PkgPath, "github.com") {
				full_pack_name = strings.Join(strings.Split(pack.PkgPath, "/")[3:], "/")
			} else {
				// bizarre path like kubernetes "k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
				if len(pack.GoFiles) > 0 {
					if strings.Contains(pack.GoFiles[0], "projects-gocurrency/") {
						replaced_path := strings.Split(strings.Split(pack.GoFiles[0], "projects-gocurrency/")[1], "/")
						full_pack_name = strings.Join(replaced_path[1:len(replaced_path)-1], "/")
					}
				}
			}

		}
		counter.Counter.Package_name = full_pack_name

		// the package path is used later to count number of lines of all go files in the folder

		// remove github.com/projectname/
		pack_path := path_to_main_dir + "/" + full_pack_name
		// append result to path_to_main_dir
		counter.Counter.Package_path = pack_path
		// Analyse each file

		if strings.Contains(full_pack_name, ".test") {
			continue
		}
		for _, file := range pack.Syntax {
			filename := fileSet.Position(file.Pos()).Filename
			go analyse.AnalyseAst(fileSet, pack_name, filename, file, package_counter_chan, filepath.Base(fileSet.Position(file.Pos()).Filename), ast_map) // launch a goroutine for each file
		}

		// Receive the results of the analysis of each file
		for range pack.Syntax {
			fmt.Print("#")
			var new_counter analyse.Counter = <-package_counter_chan

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

		packages = append(packages, &counter)

		fmt.Println()

	}

	return packages
}
