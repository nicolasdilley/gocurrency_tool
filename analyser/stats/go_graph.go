package main

import (
	"fmt"
	"html/template"
)

func GenerateGoGraph(page *PageData, counter *Counter) {
	go_map := generateGoMap(counter)

	page.Go_graph = goMapToString(go_map)
}

func generateGoMap(counter *Counter) map[int]int {
	m2 := make(map[int]int)

	for line_number, gos := range counter.go_map {
		total_go := 0.0
		for _, num_go := range gos {
			total_go += float64(num_go)
		}
		m2[line_number] = int(total_go / float64(len(gos)))
	}
	return m2
}

func goMapToString(m map[int]int) GraphData {

	dataset := "["

	keys, values := parseMap(m)

	for index, val := range keys {
		dataset += fmt.Sprintf("{x:%s,y:%s}", val, values[index])

		if index != len(keys)-1 {
			dataset += ","
		}
	}

	dataset += "]"
	var data GraphData
	data.Dataset = template.JS(dataset)
	return data
}
