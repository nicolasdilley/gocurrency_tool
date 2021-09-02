package main

import (
	"sync"

	"github.com/nicolasdilley/gocurrency_tool/analyser/tests/stru"
)

const (
	add_bound = 10
)

var (
	decl_wg *sync.WaitGroup
	decl_mu *sync.Mutex
)

func main() {
	s1 := stru.Stru1{}
	var wg_pt *sync.WaitGroup
	var mu_pt *sync.Mutex
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg1 := &sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	mu1 := &sync.Mutex{}
	mu2 := sync.Mutex{}
	a := 10
	go func() {
		wg1.Add(1)
		wg2.Add(1)
		wg.Add(a)
		wg_pt.Add(3)
		decl_wg.Add(a)
		mu.Lock()
		mu.Unlock()
		mu.Lock()
		mu1.Unlock()
		mu1.Lock()
		mu2.Unlock()
		mu_pt.Lock()
		mu_pt.Unlock()
		decl_mu.Lock()
		decl_mu.Unlock()
		s1.Mu.Unlock()
	}()
	wg1.Done()
	wg2.Done()
	wg.Done()
	wg_pt.Done()
	decl_wg.Done()
	s1.Wg.Done()
}
