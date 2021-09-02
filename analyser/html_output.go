package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	html_results_dir = "results/html"
)

type PageData struct {
	Counter                       Counter
	Project_counter               Counter
	Num_of_packages_with_features int
	Num_of_packages               int
	Num_of_features               int
	Package_counters              []*PackageCounter
	Full_project_name             string
	Line_number                   int
	Commit                        string
}

type IndexFileData struct {
	Indexes []*IndexData
}
type IndexData struct {
	Filename                      string
	Project_name                  string
	Num_of_features               int
	Line_number                   int
	Num_of_packages_with_features int
	Num_of_packages               int
}

func GenerateProjectCounter(project_counter Counter) Counter {

	project_counter.Undefined_over_defined_chans = 0.0
	if project_counter.Receive_chan_count+project_counter.Send_chan_count+
		project_counter.Param_chan_count != 0 {
		project_counter.Undefined_over_defined_chans = (float64(project_counter.Receive_chan_count+project_counter.Send_chan_count) / float64(project_counter.Receive_chan_count+project_counter.Send_chan_count+
			project_counter.Param_chan_count)) * 100
	}

	project_counter.Known_over_unknown_chan = 0.0
	project_counter.Chan_count = project_counter.Known_chan_depth_count + project_counter.Unknown_chan_depth_count + project_counter.Sync_Chan_count
	if project_counter.Chan_count != 0 {
		project_counter.Known_over_unknown_chan = (float64(project_counter.Known_chan_depth_count) /
			(float64(project_counter.Sync_Chan_count) + float64(project_counter.Unknown_chan_depth_count) + float64(project_counter.Known_chan_depth_count))) * 100
	}
	return project_counter
}

func HtmlOutputCounters(package_counters []*PackageCounter, commit string, project_name string, index_data *IndexFileData, path_to_dir string) Counter {
	var filename string = "./" + html_results_dir + "/" + strings.Replace(project_name, "/", "-", -1) + ".html"
	var project_counter Counter = Counter{Project_name: project_name}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// project_counter.Line_number = ReadNumberOfLines(GenerateListFiles(path_to_dir))
	var page PageData = PageData{
		Full_project_name: project_name,
		Commit:            commit,
		Package_counters:  package_counters,
		Line_number:       project_counter.Line_number}

	for _, counter := range package_counters {
		counter.Featured_packages = 0
		counter.Featured_files = 0
		counter.Num_files = len(counter.File_counters)
		if counter.Num_files > 0 {
			page.Num_of_packages++
		}
		for _, file := range counter.File_counters {
			counter.Featured_packages += len(file.Features)
			file.Chan_count = file.Sync_Chan_count + file.Unknown_chan_depth_count + file.Known_chan_depth_count
			if len(file.Features) != 0 {
				file.Has_feature = true
				counter.Featured_files++
				for _, feature := range file.Features {
					feature.F_commit = commit
					feature.F_project_name = project_name

					project_counter.Features = append(project_counter.Features, feature)
				}
			}
		}
		if counter.Featured_packages > 0 {
			page.Num_of_packages_with_features++
			counter.Counter.Has_feature = true
			page.Num_of_features += counter.Featured_packages
		}
		project_counter.Go_count += counter.Counter.Go_count
		project_counter.Send_count += counter.Counter.Send_count
		project_counter.Rcv_count += counter.Counter.Rcv_count
		project_counter.Chan_count += counter.Counter.Chan_count
		project_counter.Range_over_chan_count += counter.Counter.Range_over_chan_count
		project_counter.Go_in_for_count += counter.Counter.Go_in_for_count
		project_counter.Go_in_constant_for_count += counter.Counter.Go_in_constant_for_count
		project_counter.Array_of_channels_count += counter.Counter.Array_of_channels_count
		project_counter.Sync_Chan_count += counter.Counter.Sync_Chan_count
		project_counter.Known_chan_depth_count += counter.Counter.Known_chan_depth_count
		project_counter.Unknown_chan_depth_count += counter.Counter.Unknown_chan_depth_count
		project_counter.Make_chan_in_for_count += counter.Counter.Make_chan_in_for_count
		project_counter.Make_chan_in_constant_for_count += counter.Counter.Make_chan_in_constant_for_count
		project_counter.Constant_chan_array_count += counter.Counter.Constant_chan_array_count
		project_counter.Chan_slice_count += counter.Counter.Chan_slice_count
		project_counter.Chan_map_count += counter.Counter.Chan_map_count
		project_counter.Close_chan_count += counter.Counter.Close_chan_count
		project_counter.Select_count += counter.Counter.Select_count
		project_counter.Default_select_count += counter.Counter.Default_select_count
		project_counter.Assign_chan_in_for_count += counter.Counter.Assign_chan_in_for_count
		project_counter.Chan_of_chans_count += counter.Counter.Chan_of_chans_count
		project_counter.Send_chan_count += counter.Counter.Send_chan_count
		project_counter.Receive_chan_count += counter.Counter.Receive_chan_count
		project_counter.Param_chan_count += counter.Counter.Param_chan_count
	}
	project_counter.Num_of_packages_with_features = page.Num_of_packages_with_features

	if index_data != nil {

		var index *IndexData = &IndexData{}
		index.Filename = filename
		index.Project_name = project_name
		index.Num_of_features = page.Num_of_features
		index.Num_of_packages_with_features = page.Num_of_packages_with_features
		index.Num_of_packages = len(package_counters)
		index.Line_number = project_counter.Line_number

		index_data.Indexes = append(index_data.Indexes, index)
	}

	page.Project_counter = GenerateProjectCounter(project_counter)
	tmpl := template.Must(template.ParseFiles("../analyser/html_layout.html"))
	tmpl.Execute(f, page) // write the data to the file

	return project_counter
}

func GeneratePackageListFiles(path_to_dir string) string {
	git_cmd := exec.Command("ls")
	git_cmd.Dir = path_to_dir
	var git_out bytes.Buffer
	git_cmd.Stdout = &git_out
	err := git_cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while running git ls-files : ", err)
	}
	filenames := ""
	for _, name := range strings.Split(git_out.String(), "\n") {
		if strings.HasSuffix(name, ".go") {
			filenames += path_to_dir + "/" + name + "\n"
		}
	}

	return filenames
}

func GenerateListFiles(path_to_dir string) string {
	git_cmd := exec.Command("git", "ls-files")
	git_cmd.Dir = path_to_dir
	var git_out bytes.Buffer
	git_cmd.Stdout = &git_out
	err := git_cmd.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while running git ls-files : ", err)
	}
	filenames := ""
	for _, name := range strings.Split(git_out.String(), "\n") {
		if name != "" && strings.HasSuffix(name, ".go") &&
			!strings.HasPrefix(name, "vendors") &&
			!strings.HasPrefix(name, "./vendors") &&
			!strings.HasPrefix(name, "vendor") &&
			!strings.HasPrefix(name, "./vendor") &&
			!strings.HasPrefix(name, "./test") &&
			!strings.HasPrefix(name, "./tests") &&
			!strings.HasPrefix(name, "test") &&
			!strings.HasPrefix(name, "tests") {
			filenames += path_to_dir + "/" + name + "\n"
		}
	}
	return filenames
}
func ReadNumberOfLines(list_filenames string) int {
	var xargs_out bytes.Buffer
	var git_out bytes.Buffer
	filenames := strings.Split(list_filenames, "\n")
	git_out.Reset()
	for _, filename := range strings.Split(list_filenames, "\n") {
		if filename != "" {
			git_out.WriteString("\"" + filename + "\"\n")
		}
	}
	xargs_cmd := exec.Command("xargs", "cat")
	xargs_cmd.Stdin = &git_out
	xargs_cmd.Stdout = &xargs_out
	xargs_cmd.Run()

	f, _ := os.Create("temp.go")
	f.Write(xargs_out.Bytes())
	var wc_out bytes.Buffer
	wc_cmd := exec.Command("cloc", "temp.go", "--csv")
	// wc_cmd.Stdin = &xargs_out
	wc_cmd.Stdout = &wc_out
	err3 := wc_cmd.Run()
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Error while running wc : ", err3)
	}
	os.Remove("temp.go")
	word_count := strings.Split(strings.TrimSpace(wc_out.String()), "\n")
	cloc_infos := strings.Split(strings.TrimSpace(word_count[len(word_count)-1]), ",")

	if len(cloc_infos) >= 5 {
		num, _ := strconv.Atoi(cloc_infos[4])

		if len(filenames) == 0 {
			return 0
		}

		return num
	} else {
		return 0.0
	}
}
