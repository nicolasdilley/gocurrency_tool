package main

import (
	"fmt"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicolasdilley/gocurrency_tool/analyser/analyse"
	"github.com/nicolasdilley/gocurrency_tool/analyser/output"
	"golang.org/x/tools/go/packages"
)

func main() {

	os.Mkdir("results", 0755)
	os.Mkdir(output.Csv_result_dir, 0755)
	os.Mkdir(output.Html_results_dir, 0755)

	if os.Args[1] == "test" {
		ast_map, fileSet := generateAstMap("./tests")
		var packages []*output.PackageCounter = ParseDir("test", "", fileSet, ast_map)
		var test_counter analyse.Counter = output.HtmlOutputCounters(packages, "test", "test", nil, "")

		output.OutputCounters("tests", packages, "", test_counter)
		return
	}

	data, e := ioutil.ReadFile(os.Args[1])

	if e != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", os.Args[1], e)
		return
	}
	proj_listings := strings.Split(string(data), "\n")

	var index_data *output.IndexFileData = &output.IndexFileData{Indexes: []*output.IndexData{}}

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
			var packages []*output.PackageCounter

			// Generating the type info of the project

			ast_map, fileSet := generateAstMap(path_to_dir)
			packages = ParseDir(proj_name, path_to_dir, fileSet, ast_map)

			var project_counter analyse.Counter = output.HtmlOutputCounters(packages, commit_hash, project_name, index_data, path_to_dir) // html
			output.OutputCounters(project_name, packages, path_to_dir, project_counter)                                                   // csvs
		}
	}
	createIndexFile(index_data) // index html

}
func createIndexFile(index_data *output.IndexFileData) {
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.ParseFiles("../analyser/index_layout.html"))
	tmpl.Execute(f, index_data) // write the index page
}

func generateAstMap(path_to_dir string) (map[string]*packages.Package, *token.FileSet) {
	package_names := []string{}

	filepath.Walk(path_to_dir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			if file.Name() != "vendor" || file.Name() != "third_party" || file.Name() != "tests" || file.Name() != "test" {
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

	return ast_map, cfg.Fset
}
