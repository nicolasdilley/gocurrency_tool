package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Counter struct {
	Go_count                          int     // Count how many time the term "go" appear in source code
	Send_count                        int     // Count how many time a send  "chan <- val" appear in the source code
	Rcv_count                         int     // Count how many time a rcv "val <- chan" appear in the source code
	Chan_count                        int     // the number of channel overall
	Go_in_for_count                   int     // Count how many times.a goroutine is started in a for loop
	Range_over_chan_count             int     // Count the number of range over a chan
	Go_in_constant_for_count          int     // Goroutine launched in a for loop where the looping is controled by a constant
	Array_of_channels_count           int     // How many unknown length arrays are made chan of
	Sync_Chan_count                   int     // Count how many chan are created in the source code "make(chan type)"
	Known_chan_depth_count            int     // How many make(chan int, n) where n is either a constant or a hard coded number
	Unknown_chan_depth_count          int     // How many make(chan int, n) where n is completely dynamic
	Make_chan_in_for_count            int     // How many time a channel is created in a for loop
	Make_chan_in_constant_for_count   int     // How many time a channel is created in a constant for loop
	Constant_chan_array_count         int     // How many array of channels of constant size
	Chan_slice_count                  int     // How many dynamic array of channels
	Chan_map_count                    int     // how many map of channels
	Close_chan_count                  int     // How many close(chan)
	Select_count                      int     // how many select
	Default_select_count              int     // how many select with a default
	Assign_chan_in_for_count          int     // How many chan are assigned another chan in a for loop
	Assign_chan_in_constant_for_count int     // How many chan are assigned another chan in a for loop
	Chan_of_chans_count               int     // How many channel of channels
	Receive_chan_count                int     // how many receive chan
	Send_chan_count                   int     // how many send only chan
	Param_chan_count                  int     // How many times a chan is used as a param without specifying receives only or write only
	Waitgroup_count                   int     // How many waitgroup declaration are contained
	Known_add_count                   int     // How many known bound of add(n) where n is a constant
	Unknown_add_count                 int     // How many unknown bound of add(n) where n is not a constant
	Done_count                        int     // How many wg.Done()
	Mutex_count                       int     // How many mutex declaration where found
	Unlock_count                      int     // How many unlock in the code
	Lock_count                        int     // How many lock in the code
	IsPackage                         bool    // Return if the counter represent the counter for just a file or the whole package
	Package_name                      string  // The name of the package
	Package_path                      string  // path of the package
	Project_name                      string  // The name of the whole project
	Line_number                       int     // The number of lines in the counter
	Num_of_packages_with_features     int     // The number of package that contains at least one feature
	Has_feature                       bool    // Is there any features in this package ?
	Undefined_over_defined_chans      float64 // percent of undefined chan over defined (chan / chan<-, <-chan)
	Known_over_unknown_chan           float64 // percent of known chan size over unknown
	Features                          []*Feature
	filename                          string // the name of the file
}

type PackageCounter struct {
	Counter           Counter    // The overall counter of the package an
	File_counters     []*Counter // the counters of each of the file in the package
	Featured_packages int
	Featured_files    int
	Num_files         int
}

func main() {

	os.Mkdir("results", 0755)
	os.Mkdir(csv_result_dir, 0755)
	os.Mkdir(html_results_dir, 0755)

	if os.Args[1] == "test" {
		var new_counter PackageCounter = ParseDir("test", "tests", "")
		var test_counter Counter = HtmlOutputCounters([]*PackageCounter{&new_counter}, "test", "test", nil, "")

		fmt.Println(len(test_counter.Features))
		test_counter = ParseConcurrencyPrimitives("tests", test_counter) // analyses occurences of Waitgroup,mutexes and operations on them
		fmt.Println(len(test_counter.Features))

		OutputCounters("tests", []*PackageCounter{&new_counter}, "", test_counter)
		return
	}

	data, e := ioutil.ReadFile(os.Args[1])

	if e != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", os.Args[1], e)
		return
	}
	proj_listings := strings.Split(string(data), "\n")
	// var project_counters []Counter

	var index_data *IndexFileData = &IndexFileData{Indexes: []*IndexData{}}

	for _, project_name := range proj_listings {

		if project_name != "" {
			proj_name := filepath.Base(string(project_name))
			var path_to_dir string
			var commit_hash string
			path_to_dir, commit_hash = CloneRepo(string(project_name))

			_, err1 := os.Stat(path_to_dir)
			if os.IsNotExist(err1) {
				continue
			}
			var packages []*PackageCounter

			err := filepath.Walk(path_to_dir, func(path string, info os.FileInfo, err error) error {

				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path_to_dir, err)
					return err
				}
				if info.IsDir() {
					if info.Name() == "vendor" || info.Name() == "tests" || info.Name() == "test" {
						fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
						return filepath.SkipDir
					}
					var new_counter PackageCounter = ParseDir(proj_name, path, path_to_dir)
					packages = append(packages, &new_counter)
					return nil
				}
				return nil
			})

			if err != nil {
				fmt.Printf("error walking the path %q: %v\n", path_to_dir, err)
			}
			var project_counter Counter = HtmlOutputCounters(packages, commit_hash, project_name, index_data, path_to_dir) // html

			project_counter = ParseConcurrencyPrimitives(path_to_dir, project_counter) // analyses occurences of Waitgroup,mutexes and operations on them

			OutputCounters(project_name, packages, path_to_dir, project_counter) // csvs
			// project_counters = append(project_counters, project_counter)
		}
	}
	createIndexFile(index_data) // index html

}
func createIndexFile(index_data *IndexFileData) {
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.ParseFiles("../analyser/index_layout.html"))
	tmpl.Execute(f, index_data) // write the index page
}
