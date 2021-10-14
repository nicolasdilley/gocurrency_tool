package output

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nicolasdilley/gocurrency_tool/analyser/analyse"
)

type PackageCounter struct {
	Counter           analyse.Counter    // The overall counter of the package an
	File_counters     []*analyse.Counter // the counters of each of the file in the package
	Featured_packages int
	Featured_files    int
	Num_files         int
}
type PackageInfo struct {
	Name      string
	num_lines int
}

const Csv_result_dir = "results/csv"

func OutputCounters(project_name string, package_counters []*PackageCounter, path_to_dir string, project_counter analyse.Counter) {

	var filename string = "./" + Csv_result_dir + "/" + strings.Replace(project_name, "/", "-", -1) + ".csv"
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	num_featured_files := 0
	num_files := 0
	// num_featured_packages := 0
	for _, counter := range package_counters {
		num_featured_files += counter.Featured_files
		num_files += counter.Num_files
	}
	// Read number of lines in all packages in go projects

	num_packages := findNumOfPackages(path_to_dir)
	total_num_lines := 0

	packageInfos := []PackageInfo{}

	if len(package_counters) > 0 {
		filepath.Walk(path_to_dir, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {

				files, _ := ioutil.ReadDir(path)

				go_files := []string{}
				contains_go_files := false

				full_pack_name := strings.Join(strings.Split(strings.Split(path, "projects-gocurrency/")[1], "/")[1:], "/")
				for _, file := range files {
					if strings.Contains(file.Name(), ".go") {
						contains_go_files = true
						go_files = append(go_files, path+"/"+file.Name()+"\n")
					}
				}

				if contains_go_files {

					// Remove the ../project-gocurrency/author++name/ from the path
					num_lines := ReadNumberOfLines(go_files)
					total_num_lines += num_lines
					packageInfos = append(packageInfos, PackageInfo{Name: full_pack_name, num_lines: num_lines})
				}
			}

			return nil
		})

		f.WriteString(fmt.Sprintf("Line num,%d\n", total_num_lines))
		f.WriteString(fmt.Sprintf("packages num,%d,%d\n", project_counter.Num_of_packages_with_features, num_packages))
		f.WriteString(fmt.Sprintf("files num,%d,%d\n", num_featured_files, num_files))
		average_line_in_featured_file, total_line_featured_file := readNumberOfLinesPerFeaturedFile(package_counters)
		f.WriteString(fmt.Sprintf("Average Line num per featured file and total concurrent file,%d,%d\n", average_line_in_featured_file, total_line_featured_file))

		for _, pack := range packageInfos {

			f.WriteString(fmt.Sprintf("%s,%d\n", pack.Name, pack.num_lines))
		}

		f.WriteString("Filename, #,Concurrent Type, Line number , Number\n")
		for _, counter := range package_counters {
			for _, file := range counter.File_counters {
				for _, feature := range file.Features {
					f.WriteString(fmt.Sprintf("%s,%d,%s,%d,%s\n",
						feature.F_filename,
						feature.F_type_num,
						feature.F_type,
						feature.F_line_num,
						feature.F_number,
					))
				}
			}
		}
	}
}

func readNumberOfLinesPerFeaturedFile(package_counters []*PackageCounter) (int, int) {
	var git_out bytes.Buffer
	var xargs_out bytes.Buffer

	var filenames []string

	for _, package_counter := range package_counters {
		for _, counter := range package_counter.File_counters {
			if len(counter.Features) > 0 {
				// check if there is a none mutex feature
				filenames = append(filenames, counter.Filename)
			}
		}
	}

	for _, filename := range filenames {
		if filename != "" {
			git_out.WriteString(filename + "\n")
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
			return 0, 0
		}

		return num / len(filenames), num

	} else {
		return 0, 0
	}

}

func ContainsOtherFeatureThanMutex(features []*analyse.Feature) bool {
	for _, feature := range features {
		if feature.F_type_num != analyse.MUTEX_COUNT {
			return true
		}
	}

	return false
}

// Go through the path and output number of lines in the path
// returns the number of package
func findNumOfPackages(path_to_dir string) int {

	num_packages := 0

	filepath.Walk(path_to_dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {

			files, _ := ioutil.ReadDir(path)

			contains_go_files := false

			for _, file := range files {
				if strings.Contains(file.Name(), ".go") {
					contains_go_files = true
				}
			}

			if contains_go_files {
				num_packages++
			}
		}

		return nil
	})

	return num_packages
}
