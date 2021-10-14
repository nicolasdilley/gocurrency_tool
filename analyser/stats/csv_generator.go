package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/montanaflynn/stats"
)

func GenerateCSVResults(page *PageData, counter *Counter) {
	var normal_relative_filename string = "./results/normal/relative.csv"
	var normal_new_chan_filename string = "./results/normal/new_chan.csv"
	var normal_absolute_filename string = "./results/normal/absolute.csv"
	var normal_all_absolute_filename string = "./results/normal/all_absolute.csv"

	var normal_absolute_wg_filename string = "./results/normal/absolute_wg.csv"
	var normal_new_wg_filename string = "./results/normal/new_wg.csv"
	var normal_relative_wg_filename string = "./results/normal/relative_wg.csv"

	var normal_new_mu_filename string = "./results/normal/new_mu.csv"
	var normal_absolute_mu_filename string = "./results/normal/absolute_mu.csv"
	var normal_relative_mu_filename string = "./results/normal/relative_mu.csv"
	var normal_relative_all_filename string = "./results/normal/relative_all.csv"

	var median_absolute_filename string = "./results/median/absolute.csv"
	var median_all_absolute_filename string = "./results/median/all_absolute.csv"
	var median_absolute_wg_filename string = "./results/median/absolute_wg.csv"
	var median_absolute_mu_filename string = "./results/median/absolute_mu.csv"
	var median_cloc_filename string = "./results/cloc/median.csv"

	var normal_cloc_filename string = "./results/cloc/normal.csv"
	var known_size_chan_filename string = "./results/known_size_chan.csv"
	var non_zero_known_chan_filename string = "./results/non_zero_known_chan.csv"
	var go_in_any_for_filename string = "./results/go_in_any_for.csv"
	var chan_in_any_for_filename string = "./results/chan_in_any_for.csv"
	var go_in_for_filename string = "./results/go_in_unknown_for.csv"
	var chan_in_for_filename string = "./results/chan_in_for.csv"
	var go_in_constant_for_filename string = "./results/go_in_constant_for.csv"
	var go_per_projects_filename string = "./results/go_per_projects.csv"
	var constant_go_in_constant_for_filename string = "./results/constant_go_in_constant_for.csv"
	var num_branch_per_projects_filename string = "./results/num_branch_per_projects.csv"

	os.Mkdir("results", 0755)
	os.Mkdir("results/normal", 0755)
	os.Mkdir("results/median", 0755)
	os.Mkdir("results/cloc", 0755)

	normal_relative, _ := os.Create(normal_relative_filename)
	defer normal_relative.Close()

	normal_new_chan, _ := os.Create(normal_new_chan_filename)
	defer normal_new_chan.Close()

	normal_absolute, _ := os.Create(normal_absolute_filename)
	defer normal_absolute.Close()

	normal_all_absolute, _ := os.Create(normal_all_absolute_filename)
	defer normal_all_absolute.Close()

	normal_absolute_wg, _ := os.Create(normal_absolute_wg_filename)
	defer normal_absolute_wg.Close()

	normal_new_wg, _ := os.Create(normal_new_wg_filename)
	defer normal_new_wg.Close()

	normal_relative_wg, _ := os.Create(normal_relative_wg_filename)
	defer normal_relative_wg.Close()

	normal_new_mu, _ := os.Create(normal_new_mu_filename)
	defer normal_new_mu.Close()

	normal_relative_mu, _ := os.Create(normal_relative_mu_filename)
	defer normal_relative_mu.Close()

	normal_absolute_mu, _ := os.Create(normal_absolute_mu_filename)
	defer normal_absolute_mu.Close()

	normal_relative_all, _ := os.Create(normal_relative_all_filename)
	defer normal_relative_all.Close()

	median_absolute, _ := os.Create(median_absolute_filename)
	defer median_absolute.Close()

	median_all_absolute, _ := os.Create(median_all_absolute_filename)
	defer median_all_absolute.Close()

	median_absolute_wg, _ := os.Create(median_absolute_wg_filename)
	defer median_absolute_wg.Close()

	median_absolute_mu, _ := os.Create(median_absolute_mu_filename)
	defer median_absolute_mu.Close()

	median_cloc, _ := os.Create(median_cloc_filename)
	defer median_cloc.Close()

	normal_cloc, _ := os.Create(normal_cloc_filename)
	defer normal_cloc.Close()

	known_size_chan, _ := os.Create(known_size_chan_filename)
	defer known_size_chan.Close()

	non_zero_known_chan, _ := os.Create(non_zero_known_chan_filename)
	defer non_zero_known_chan.Close()

	go_in_any_for, _ := os.Create(go_in_any_for_filename)
	defer go_in_any_for.Close()

	chan_in_any_for, _ := os.Create(chan_in_any_for_filename)
	defer chan_in_any_for.Close()

	go_in_for, _ := os.Create(go_in_for_filename)
	defer go_in_for.Close()

	chan_in_for, _ := os.Create(chan_in_for_filename)
	defer chan_in_for.Close()

	go_in_constant_for, _ := os.Create(go_in_constant_for_filename)
	defer go_in_constant_for.Close()

	go_per_projects, _ := os.Create(go_per_projects_filename)
	defer go_per_projects.Close()

	constant_go_in_constant_for, _ := os.Create(constant_go_in_constant_for_filename)
	defer constant_go_in_constant_for.Close()

	num_branch_per_projects, _ := os.Create(num_branch_per_projects_filename)
	defer num_branch_per_projects.Close()

	median := calculateMedian(page)
	percent_around_median := 0.3
	num_median := 0 // the number of median projects

	for proj, features := range page.Stats.Features_per_projects {

		num_lines_in_featured_files := float64(page.Stats.Num_lines_in_featured_files.Get(proj).Concurrent_size)

		num_chan := 0.0
		num_select := 0.0
		num_send := 0.0
		num_receive := 0.0
		num_close := 0.0
		num_range_over_chan := 0.0

		num_wg := 0.0
		num_known_add := 0.0
		num_unknown_add := 0.0
		num_done := 0.0
		num_wait := 0.0

		num_mutex := 0.0
		num_unlock := 0.0
		num_lock := 0.0

		for _, feature := range features {
			switch feature.f_type_num {
			case UNKNOWN_CHAN_DEPTH_COUNT:
				num_chan++
			case KNOWN_CHAN_DEPTH_COUNT:
				num_chan++
			case MAKE_CHAN_COUNT:
				num_chan++
			case SELECT_COUNT:
				num_select++
			case DEFAULT_SELECT_COUNT:
				num_select++
			case SEND_COUNT:
				num_send++
			case RECEIVE_COUNT:
				num_receive++
			case RANGE_OVER_CHAN_COUNT:
				num_range_over_chan++
			case CLOSE_CHAN_COUNT:
				num_close++
			case WAITGROUP_COUNT:
				num_wg++
			case KNOWN_ADD_COUNT:
				num_known_add++
			case UNKNOWN_ADD_COUNT:
				num_unknown_add++
			case DONE_COUNT:
				num_done++
			case WAIT_COUNT:
				num_wait++
			case MUTEX_COUNT:
				num_mutex++
			case LOCK_COUNT:
				num_lock++
			case UNLOCK_COUNT:
				num_unlock++
			}
		}
		num_featured_lines := float64(num_lines_in_featured_files)

		if num_featured_lines == 0 {
			fmt.Println(proj, " has 0 featured lines")
		}

		if num_lines_in_featured_files > float64(median)*(1.0-percent_around_median) &&
			num_lines_in_featured_files < float64(median)*(1.0+percent_around_median) {
			num_median++

			median_all_absolute.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_chan,
				num_send,
				num_receive,
				num_select,
				num_close,
				num_range_over_chan,

				num_wg,
				num_known_add+
					num_unknown_add,
				num_done,
				num_wait,

				num_mutex,
				num_lock,
				num_unlock,
			))
		}

		if num_chan > 0 {

			normal_relative.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				(num_chan/num_featured_lines)*1000,
				(num_send/num_featured_lines)*1000,
				(num_receive/num_featured_lines)*1000,
				(num_select/num_featured_lines)*1000,
				(num_close/num_featured_lines)*1000,
				(num_range_over_chan/num_featured_lines)*1000))

			normal_new_chan.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_send/num_chan,
				num_receive/num_chan,
				num_select/num_chan,
				num_close/num_chan,
				num_range_over_chan/num_chan))

			normal_absolute.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_chan,
				num_send,
				num_receive,
				num_select,
				num_close,
				num_range_over_chan,
			))

			normal_cloc.WriteString(fmt.Sprintf("%s,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				float64(num_lines_in_featured_files)/float64(features[0].f_number_of_lines)*100.0,
				float64(features[0].f_featured_packages)/float64(features[0].f_num_package)*100.0,
				float64(features[0].f_num_featured_files)/float64(features[0].f_num_files)*100))

			if num_lines_in_featured_files > float64(median)*(1.0-percent_around_median) &&
				num_lines_in_featured_files < float64(median)*(1.0+percent_around_median) {

				median_absolute.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%f\n",
					strings.TrimSuffix(proj, ".csv"),
					num_chan,
					num_send,
					num_receive,
					num_select,
					num_close,
					num_range_over_chan,
					num_lines_in_featured_files,
				))

				median_cloc.WriteString(fmt.Sprintf("%s,%.2f,%.2f,%.2f\n",
					strings.TrimSuffix(proj, ".csv"),
					float64(num_lines_in_featured_files)/float64(features[0].f_number_of_lines)*100.0,
					float64(features[0].f_featured_packages)/float64(features[0].f_num_package)*100.0,
					float64(features[0].f_num_featured_files)/float64(features[0].f_num_files)*100,
				))
			}
		} else {
			fmt.Println("proj ", proj, " does not contain channels")
		}

		if num_wg > 0 {

			normal_absolute_wg.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_wg,
				num_known_add,
				num_unknown_add,
				num_done,
				num_wait,
			))

			normal_relative_wg.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				(num_wg/num_featured_lines)*1000,
				(num_known_add/num_featured_lines)*1000,
				(num_unknown_add/num_featured_lines)*1000,
				(num_done/num_featured_lines)*1000,
				(num_wait/num_featured_lines)*1000,
			))

			normal_new_wg.WriteString(fmt.Sprintf("%s,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				(num_known_add+num_unknown_add)/num_wg,
				num_done/num_wg,
				num_wait/num_wg,
			))

			if num_lines_in_featured_files > float64(median)*(1.0-percent_around_median) &&
				num_lines_in_featured_files < float64(median)*(1.0+percent_around_median) {

				median_absolute_wg.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f\n",
					strings.TrimSuffix(proj, ".csv"),
					num_wg,
					num_known_add,
					num_unknown_add,
					num_done,
					num_wait,
					num_lines_in_featured_files,
				))
			}
		}

		if num_lock > 0 || num_unlock > 0 || num_mutex > 0 {

			// all projects
			normal_absolute_mu.WriteString(fmt.Sprintf("%s,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_mutex,
				num_lock,
				num_unlock,
			))

			// all projects
			normal_relative_mu.WriteString(fmt.Sprintf("%s,%f,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				(num_mutex/num_lines_in_featured_files)*1000,
				(num_lock/num_lines_in_featured_files)*1000,
				(num_unlock/num_lines_in_featured_files)*1000,
			))
			// all projects
			normal_new_mu.WriteString(fmt.Sprintf("%s,%f,%f\n",
				strings.TrimSuffix(proj, ".csv"),
				num_lock/num_mutex,
				num_unlock/num_mutex,
			))

			// fmt.Println(num_lines_in_featured_files, float64(median)*(1.0-percent_around_median), float64(median)*(1.0+percent_around_median))
			if num_lines_in_featured_files > float64(median)*(1.0-percent_around_median) &&
				num_lines_in_featured_files < float64(median)*(1.0+percent_around_median) {

				median_absolute_mu.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f\n",
					strings.TrimSuffix(proj, ".csv"),
					num_mutex,
					num_lock,
					num_unlock,
					num_lines_in_featured_files,
				))
			}
		}

		normal_all_absolute.WriteString(fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n",
			strings.TrimSuffix(proj, ".csv"),
			num_chan,
			num_send,
			num_receive,
			num_select,
			num_close,
			num_range_over_chan,

			num_wg,
			num_known_add+
				num_unknown_add,
			num_done,
			num_wait,

			num_mutex,
			num_lock,
			num_unlock,
		))

		// if  waitgroup  generate absolute and relative
		// if unlock or lock generate absolute and relative
	}

	for _, size := range counter.known_chan_map {
		known_size_chan.WriteString(fmt.Sprintf("%s,%d\n", size.proj_name, size.chan_size))

		if size.chan_size != 0 {
			non_zero_known_chan.WriteString(fmt.Sprintf("%s,%d\n", size.proj_name, size.chan_size))
		}
	}

	for proj, sizes := range counter.constant_go_in_constant_for_map {
		for _, size := range sizes {
			constant_go_in_constant_for.WriteString(fmt.Sprintf("%s,%d\n", proj, size))
		}
	}
	for proj, sizes := range counter.num_branch_per_projects {
		for _, size := range sizes {
			num_branch_per_projects.WriteString(fmt.Sprintf("%s,%d\n", proj, size))
		}
	}

	for proj, size := range counter.go_in_any_for_map {
		go_in_any_for.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}

	for proj, size := range counter.chan_in_any_for_map {
		chan_in_any_for.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}
	for proj, size := range counter.go_in_for_map {
		go_in_for.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}

	for proj, size := range counter.chan_in_for_map {
		chan_in_for.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}
	for proj, size := range counter.go_in_constant_for_map {
		go_in_constant_for.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}
	for proj, size := range counter.go_per_projects {
		go_per_projects.WriteString(fmt.Sprintf("%s,%f\n", proj, 1000*float64(size)/page.Stats.Features_per_projects[proj+".csv"][0].f_featured_file_line_average))
	}

	page.Stats.Num_median_projects = num_median

}

func calculateIQR(features_per_projects map[string][]*Feature) (int, int, int) {
	var lines stats.Float64Data

	for name, _ := range features_per_projects {
		lines = append(lines, float64(features_per_projects[name][0].f_number_of_lines))
	}
	sort.Float64s(lines)
	quartile, _ := stats.Quartile(lines)

	fmt.Printf("min: %f, q2:%f max: %f\n", quartile.Q1, quartile.Q2, quartile.Q3)
	return int(quartile.Q1), int(quartile.Q2), int(quartile.Q3)
}

func calculateMedian(page *PageData) int {
	var lines stats.Float64Data

	for _, project := range page.Stats.Num_lines_in_featured_files {
		lines = append(lines, float64(project.Concurrent_size))
	}

	sort.Float64s(lines)
	quartile, _ := stats.Quartile(lines)

	fmt.Printf("median: %d\n", int(quartile.Q2))
	return int(quartile.Q2)
}
