package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"sort"
)

func GenerateChannelSizeGraph(page *PageData, counter *Counter) {
	page.Channel_size_graph = mapToString(counter.chan_size_map)
}

func mapToString(m map[int]int) GraphData {

	dataset := "["
	labels := "["
	colors := "["

	keys, values := parseMap(m)

	for index, val := range keys {
		dataset += values[index]
		labels += "'=+" + val + "'"
		colors += "'rgba(" + fmt.Sprint(int(rand.Float64()*255.0)) + "," + fmt.Sprint(int(rand.Float64()*255.0)) + "," + fmt.Sprint(int(rand.Float64()*255.0)) + ",0.3)'"

		if index != len(keys)-1 {
			dataset += ","
			labels += ","
			colors += ","
		}
	}

	dataset += "]"
	labels += "]"
	colors += "]"
	var data GraphData
	data.Colors = template.JS(colors)
	data.Labels = template.JS(labels)
	data.Dataset = template.JS(dataset)
	return data
}

func parseMap(m map[int]int) ([]string, []string) {

	keys := []string{}
	values := []string{}

	num_keys := []int{}
	for key, _ := range m {
		num_keys = append(num_keys, key)
	}

	sort.Ints(num_keys)

	for _, k := range num_keys {
		key := fmt.Sprint(k)
		value := fmt.Sprint(m[k])
		keys = append(keys, key)
		values = append(values, value)
	}

	return keys, values
}
