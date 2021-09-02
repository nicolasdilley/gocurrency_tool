package stru

import "sync"

type Stru1 struct {
	Wg *sync.WaitGroup
	Mu *sync.Mutex
}

type Stru2 struct {
	stru Stru1
}

type Stru3 struct {
	*sync.WaitGroup
	Mu *sync.Mutex
}
type Stru4 struct {
	*sync.Mutex
}
