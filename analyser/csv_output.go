package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const csv_result_dir = "results/csv"

func OutputCounters(project_name string, package_counters []*PackageCounter, path_to_dir string, project_counter Counter) {

	var filename string = "./" + csv_result_dir + "/" + strings.Replace(project_name, "/", "-", -1) + ".csv"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	var num_of_packages int = 0

	num_featured_files := 0
	num_files := 0
	// num_featured_packages := 0
	for _, counter := range package_counters {
		if len(counter.File_counters) > 0 {
			num_of_packages++
		}
		num_featured_files += counter.Featured_files
		num_files += counter.Num_files
	}
	if len(package_counters) > 0 {
		f.WriteString(fmt.Sprintf("Line num,%d\n", project_counter.Line_number))
		f.WriteString(fmt.Sprintf("packages num,%d,%d\n", project_counter.Num_of_packages_with_features, num_of_packages))
		f.WriteString(fmt.Sprintf("files num,%d,%d\n", num_featured_files, num_files))
		f.WriteString(fmt.Sprintf("Average Line num per featured file,%.d\n", readNumberOfLinesPerFeaturedFile(package_counters)))
	}

	for _, counter := range package_counters {
		if len(counter.File_counters) > 0 {
			f.WriteString(fmt.Sprintf("%s,%d\n", counter.Counter.Package_name, ReadNumberOfLines(GeneratePackageListFiles(counter.Counter.Package_path))))
		}
	}

	f.WriteString("Filename, #,Concurrent Type, Line number , Number\n")

	for _, feature := range project_counter.Features {
		f.WriteString(fmt.Sprintf("%s,%d,%s,%d,%s\n",
			feature.F_filename,
			feature.F_type_num,
			feature.F_type,
			feature.F_line_num,
			feature.F_number,
		))
	}
}

func readNumberOfLinesPerFeaturedFile(package_counters []*PackageCounter) int {
	var git_out bytes.Buffer
	var xargs_out bytes.Buffer

	var filenames []string

	for _, package_counter := range package_counters {
		for _, counter := range package_counter.File_counters {
			if len(counter.Features) > 0 {
				filenames = append(filenames, counter.filename)
			}
		}
	}
	for _, filename := range filenames {
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
