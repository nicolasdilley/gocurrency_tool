package main

import (
	"sync"
)

type Test struct {
	i  int
	wg *sync.WaitGroup
}

func main() {
	t := Test{i: 0}

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)
	t.wg.Add(10)
	t.wg.Done()
	mu.Lock()
}
