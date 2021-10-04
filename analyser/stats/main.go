package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	csv_folder  = "../results/csv"
	html_folder = "../results/html"

	GOROUTINE_COUNT                 int = 1
	RECEIVE_COUNT                   int = 2
	SEND_COUNT                      int = 3
	MAKE_CHAN_COUNT                 int = 4
	GO_IN_FOR_COUNT                 int = 5
	RANGE_OVER_CHAN_COUNT           int = 6
	GO_IN_CONSTANT_FOR_COUNT        int = 7
	KNOWN_CHAN_DEPTH_COUNT          int = 8
	UNKNOWN_CHAN_DEPTH_COUNT        int = 9
	MAKE_CHAN_IN_FOR_COUNT          int = 10
	MAKE_CHAN_IN_CONSTANT_FOR_COUNT int = 23
	ARRAY_OF_CHANNELS_COUNT         int = 11
	CONSTANT_CHAN_ARRAY_COUNT       int = 12
	CHAN_SLICE_COUNT                int = 13
	CHAN_MAP_COUNT                  int = 14
	CLOSE_CHAN_COUNT                int = 15
	SELECT_COUNT                    int = 16
	DEFAULT_SELECT_COUNT            int = 17
	ASSIGN_CHAN_IN_FOR_COUNT        int = 18
	CHAN_OF_CHANS_COUNT             int = 19
	RECEIVE_CHAN_COUNT              int = 20
	SEND_CHAN_COUNT                 int = 21
	PARAM_CHAN_COUNT                int = 22
	WAITGROUP_COUNT                 int = 30
	KNOWN_ADD_COUNT                 int = 24
	UNKNOWN_ADD_COUNT               int = 25
	DONE_COUNT                      int = 26
	WAIT_COUNT                      int = 31
	MUTEX_COUNT                     int = 27
	UNLOCK_COUNT                    int = 28
	LOCK_COUNT                      int = 29

	constant = 100000.0
)

type Feature struct {
	f_type                       string // type
	f_type_num                   int    // type num
	f_filename                   string
	f_package_name               string
	f_line_num                   int
	f_number                     string // A number used to report additional info about a feature
	f_number_of_lines            int
	f_featured_file_line_average float64 // the number of lines on average in featured files
	f_featured_packages          int
	f_num_package                int
	f_num_files                  int
	f_num_featured_files         int
}

type ChanSizeInfo struct {
	proj_name string
	chan_size int
}

type Counter struct {
	channels                          int
	Zero_chans                        int
	Non_zero_know_chans               int
	Known_chans                       int // num of non-zero and zero chans
	Unknown_chans                     int
	goroutines                        float64
	go_in_for                         float64
	dynamic_structures                float64
	defined_chans                     float64
	chan_size_map                     map[int]int
	go_map                            map[int][]int
	chan_map                          map[int][]int
	go_int_for_over_go                int
	chan_in_for                       int
	make_chan                         int
	known_chan_map                    []ChanSizeInfo
	constant_go_in_constant_for_map   map[string][]int
	go_in_any_for_map                 map[string]int
	chan_in_any_for_map               map[string]int
	go_in_constant_for_map            map[string]int
	go_in_for_map                     map[string]int
	chan_in_for_map                   map[string]int
	go_per_projects                   map[string]int
	num_branch_per_projects           map[string][]int
	Project_with_interesting_features int

	// Waitgroup
	wg          map[string]int
	known_add   map[string]int
	unknown_add map[string]int
	done        map[string]int

	// Mutex
	mutex  map[string]int
	unlock map[string]int
	lock   map[string]int
}

type ProjectCounter struct {
	go_in_for          int
	chan_in_for        int
	assign_chan_in_for int
	params_chan        int
	contains_no_chan   int // how many project without chans
}

// statsInfo is a struct that holds combined projects stats
type StatsInfo struct {
	Num_projects                          int            // the number of projects analyzed
	Num_featured_projects                 int            // the number of projects with at least one feature
	Num_projects_without_chans            int            // the number of projects analyzed without a new channel
	Num_of_package_with_features          int            // how many packages with features
	Num_of_packages_from_featured_project int            // num of packages only in featured project
	Num_packages                          int            // the total number of packages
	Num_features                          int            // the total number of features
	Num_files                             int            // total number of files in the project
	Num_featured_files                    int            // total number of featured files in the project
	Num_lines_in_featured_files           map[string]int // total number of featured files in the project
	Project_with_interesting_features     int
	Features_per_projects                 map[string][]*Feature // a map of all the features per project name
	Num_assign_chan_in_for                int
	Average                               AverageInfo
}

type AverageInfo struct {
	Channels                        int
	Zero_chans                      int
	Non_zero_know_chans             int
	Known_chans                     int // num of non-zero and zero chans
	Unknown_chans                   int
	Go_in_for                       float64
	Dynamic_structures              float64
	Goroutines                      float64
	Go_in_for_over_go               float64
	Channels_in_for_over_chan       float64
	Featured_packages_over_packages float64
	Project_with_go_in_for          float64
	Project_with_chan_in_for        float64
	Features_over_featured_packages float64
	Defined_over_params             float64
}

type PageData struct {
	Stats                  StatsInfo
	Constant               float64
	Channel_size_graph     GraphData
	Go_graph               GraphData
	Chan_graph             GraphData
	Features_lines_of_code GraphData
}

type GraphData struct {
	Dataset template.JS
	Colors  template.JS
	Labels  template.JS
}

func main() {

	var page PageData = PageData{Constant: constant}
	page.Stats = generatesFeaturesList()

	f := createStatFile()                            // generate the file for the .html file
	counter, projectCounter := getStats(&page.Stats) // generate the statsInfo object

	page.Stats.Num_projects_without_chans = projectCounter.contains_no_chan
	GenerateCSVResults(&page, &counter)
	tmpl := template.Must(template.ParseFiles("layout.html"))
	tmpl.Execute(f, page) // write the data to the file
	return
}

func generatesFeaturesList() StatsInfo {
	var stats StatsInfo = StatsInfo{
		Features_per_projects:       map[string][]*Feature{},
		Num_lines_in_featured_files: map[string]int{},
	}

	num_line_overall := 0

	err := filepath.Walk(csv_folder, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {

			if strings.HasSuffix(info.Name(), ".csv") { // look only for .csv files

				csv_data, _ := ioutil.ReadFile(path)
				csv_strings := strings.Split(string(csv_data), "\n")
				stats.Num_projects++
				line_number, _ := strconv.Atoi(strings.Split(csv_strings[0], ",")[1])
				num_line_overall += line_number
				packages_line := strings.Split(csv_strings[1], ",")

				if len(csv_strings) <= 2 {
					fmt.Println("No go code in ", info.Name())
					return nil
				}
				featured_files, _ := strconv.Atoi(strings.Split(csv_strings[2], ",")[1])
				num_files, _ := strconv.Atoi(strings.Split(csv_strings[2], ",")[2])

				if strings.Split(csv_strings[3], ",")[1] == "" {
					fmt.Println("no features ", info.Name())
				}
				average_line_number_of_featured_file, _ := strconv.ParseFloat(strings.Split(csv_strings[3], ",")[1], 0)
				line_number_of_featured_file, _ := strconv.Atoi(strings.Split(csv_strings[3], ",")[2])

				stats.Num_lines_in_featured_files[info.Name()] = line_number_of_featured_file

				feat_pack, _ := strconv.Atoi(packages_line[1])
				stats.Num_of_package_with_features += feat_pack
				pack, _ := strconv.Atoi(packages_line[2])

				if feat_pack > 0 {
					stats.Num_of_packages_from_featured_project += pack
				}

				stats.Num_packages += pack
				stats.Num_files = num_files
				stats.Num_featured_files = featured_files

				var s_features []string = csv_strings[pack+5:] // Features start at line 5

				for _, s_feature := range s_features {
					if s_feature != "" {
						s := strings.Split(s_feature, ",")
						fmt.Println("ici ", s)
						type_num, _ := strconv.Atoi(s[1])
						line_num, _ := strconv.Atoi(s[3])

						stats.Features_per_projects[info.Name()] = append(stats.Features_per_projects[info.Name()], &Feature{
							f_filename:                   s[0],
							f_type_num:                   type_num,
							f_type:                       s[2],
							f_line_num:                   line_num,
							f_number:                     s[4],
							f_number_of_lines:            line_number,
							f_featured_file_line_average: average_line_number_of_featured_file,
							f_num_package:                pack,
							f_num_featured_files:         featured_files,
							f_featured_packages:          feat_pack,
							f_num_files:                  num_files,
						})
						stats.Num_features++

					}
				}
				return nil
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	stats.Num_featured_projects = len(stats.Features_per_projects)
	fmt.Printf("num of loc overall : %d\n", num_line_overall)
	return stats
}

func getStats(stats *StatsInfo) (Counter, ProjectCounter) {

	// read and process csv
	counter, project_counter := getCounters(stats.Features_per_projects)

	// calculate averages based on counter generated
	stats.Average.Channels = counter.channels
	stats.Average.Known_chans = counter.Known_chans
	stats.Average.Zero_chans = counter.Zero_chans
	stats.Average.Non_zero_know_chans = counter.Non_zero_know_chans
	stats.Average.Unknown_chans = counter.Unknown_chans
	stats.Average.Go_in_for = (counter.go_in_for / float64(stats.Num_projects))
	stats.Average.Dynamic_structures = (counter.dynamic_structures / float64(stats.Num_projects))
	stats.Average.Goroutines = (counter.goroutines / float64(stats.Num_projects))
	stats.Average.Go_in_for_over_go = (float64(counter.go_in_for) / float64(counter.goroutines)) * 100
	stats.Average.Channels_in_for_over_chan = (float64(counter.chan_in_for) / float64(counter.make_chan)) * 100
	stats.Average.Featured_packages_over_packages = (float64(stats.Num_of_package_with_features) / float64(stats.Num_of_packages_from_featured_project)) * 100
	stats.Average.Project_with_go_in_for = float64(project_counter.go_in_for) / float64(stats.Num_projects) * 100
	stats.Average.Project_with_chan_in_for = float64(project_counter.chan_in_for) / float64(stats.Num_projects) * 100
	stats.Average.Features_over_featured_packages = float64(stats.Num_features) / float64(stats.Num_of_package_with_features)
	stats.Average.Defined_over_params = counter.defined_chans / float64(project_counter.params_chan)
	stats.Num_assign_chan_in_for = project_counter.assign_chan_in_for
	stats.Project_with_interesting_features = counter.Project_with_interesting_features

	return counter, project_counter
}

func getCounters(features_per_projects map[string][]*Feature) (Counter, ProjectCounter) {
	var counter Counter
	var project_counter ProjectCounter

	counter.chan_size_map = make(map[int]int)
	counter.go_map = make(map[int][]int)
	counter.chan_map = make(map[int][]int)
	counter.known_chan_map = []ChanSizeInfo{}
	counter.go_in_any_for_map = make(map[string]int)
	counter.chan_in_any_for_map = make(map[string]int)
	counter.go_in_for_map = make(map[string]int)
	counter.chan_in_for_map = make(map[string]int)
	counter.go_in_constant_for_map = make(map[string]int)
	counter.go_per_projects = make(map[string]int)
	counter.mutex = make(map[string]int)
	counter.wg = make(map[string]int)
	counter.known_add = make(map[string]int)
	counter.unknown_add = make(map[string]int)
	counter.done = make(map[string]int)
	counter.lock = make(map[string]int)
	counter.unlock = make(map[string]int)
	counter.constant_go_in_constant_for_map = make(map[string][]int)
	counter.num_branch_per_projects = make(map[string][]int)

	for proj, features := range features_per_projects {
		containsGoInFor := false
		containsChanInFor := false
		containsAssignChanInFor := false
		proj = strings.TrimSuffix(proj, ".csv")

		num_go := 0
		num_chan := 0
		num_select := 0
		num_receive := 0
		num_range := 0
		num_send := 0
		defined_chans := 0
		undefined_chans := 0
		num_wg := 0
		num_mutex := 0

		for _, feature := range features {
			switch feature.f_type_num {
			case RECEIVE_COUNT:
				num_receive++
			case SELECT_COUNT:
				num_select++
				num_branch, _ := strconv.Atoi(feature.f_number)
				counter.num_branch_per_projects[proj] = append(counter.num_branch_per_projects[proj], num_branch)
			case DEFAULT_SELECT_COUNT:
				num_select++
				num_branch, _ := strconv.Atoi(feature.f_number)
				counter.num_branch_per_projects[proj] = append(counter.num_branch_per_projects[proj], num_branch)
			case RANGE_OVER_CHAN_COUNT:
				num_range++
			case SEND_COUNT:
				num_send++
			case UNKNOWN_CHAN_DEPTH_COUNT:
				counter.Unknown_chans++
				counter.make_chan++
				num_chan++
			case GOROUTINE_COUNT:
				counter.goroutines += 1 / (float64(feature.f_number_of_lines) / constant)
				num_go++
			case GO_IN_FOR_COUNT:
				counter.go_in_any_for_map[proj]++
				containsGoInFor = true
				counter.go_in_for += 1 / (float64(feature.f_number_of_lines) / constant)
			case GO_IN_CONSTANT_FOR_COUNT:
				counter.go_in_constant_for_map[proj]++
				num, _ := strconv.Atoi(feature.f_number)
				counter.constant_go_in_constant_for_map[proj] = append(counter.constant_go_in_constant_for_map[proj], num)
			case CHAN_MAP_COUNT, CHAN_SLICE_COUNT:
				counter.dynamic_structures += 1 / (float64(feature.f_number_of_lines) / constant)
			case SEND_CHAN_COUNT, RECEIVE_CHAN_COUNT:
				defined_chans++
			case PARAM_CHAN_COUNT:
				undefined_chans++
			case KNOWN_CHAN_DEPTH_COUNT:
				counter.Non_zero_know_chans++
				counter.Known_chans++
				num, _ := strconv.Atoi(feature.f_number)
				if num > 1000 {
					counter.chan_size_map[1000]++
				} else if num > 100 {
					counter.chan_size_map[100]++
				} else if num > 10 {
					counter.chan_size_map[10]++
				} else if num > 3 {
					counter.chan_size_map[3]++
				} else {
					counter.chan_size_map[num]++
				}
				counter.make_chan++
				counter.known_chan_map = append(counter.known_chan_map, ChanSizeInfo{chan_size: num, proj_name: proj})
				num_chan++
			case MAKE_CHAN_COUNT:
				counter.Zero_chans++
				counter.Known_chans++
				counter.chan_size_map[0]++
				counter.make_chan++
				counter.known_chan_map = append(counter.known_chan_map, ChanSizeInfo{chan_size: 0, proj_name: proj})
				num_chan++
			case MAKE_CHAN_IN_FOR_COUNT:
				counter.chan_in_for++
				containsChanInFor = true
				counter.chan_in_any_for_map[proj]++
				counter.chan_in_for_map[proj]++
			case MAKE_CHAN_IN_CONSTANT_FOR_COUNT:
				counter.chan_in_any_for_map[proj]++
			case ASSIGN_CHAN_IN_FOR_COUNT:
				containsAssignChanInFor = true
			case WAITGROUP_COUNT:
				counter.wg[proj]++
				num_wg++
			case KNOWN_ADD_COUNT:
				counter.known_add[proj]++
			case UNKNOWN_ADD_COUNT:
				counter.unknown_add[proj]++
			case DONE_COUNT:
				counter.done[proj]++
			case MUTEX_COUNT:
				counter.mutex[proj]++
				num_mutex++
			case UNLOCK_COUNT:
				counter.unlock[proj]++
			case LOCK_COUNT:
				counter.lock[proj]++
			}
		}
		if containsChanInFor {
			project_counter.chan_in_for++
		}
		if containsGoInFor {
			project_counter.go_in_for++
		}
		if containsAssignChanInFor {
			project_counter.assign_chan_in_for++
		}
		if defined_chans+undefined_chans != 0 {
			counter.defined_chans += (float64(defined_chans) / float64(defined_chans+undefined_chans)) * 100
			project_counter.params_chan++
		}

		if counter.go_in_any_for_map[proj]-counter.go_in_constant_for_map[proj] > 0 {
			counter.go_in_for_map[proj] = counter.go_in_any_for_map[proj] - counter.go_in_constant_for_map[proj]
		}

		if len(features) > 0 {
			line_num := features[0].f_number_of_lines
			counter.go_map[line_num] = append(counter.go_map[line_num], num_go)
			counter.chan_map[line_num] = append(counter.chan_map[line_num], num_chan)
		}
		if num_go > 0 {
			counter.go_per_projects[proj] += num_go
		}

		if num_chan == 0 {
			project_counter.contains_no_chan++
		}

		if num_chan > 0 || counter.lock[proj] > 0 || counter.unlock[proj] > 0 || counter.wg[proj] > 0 {
			counter.Project_with_interesting_features++
		}
		counter.channels += num_chan
	}

	return counter, project_counter
}

func createStatFile() *os.File {
	f, _ := os.Create("stats.html")
	return f
}
