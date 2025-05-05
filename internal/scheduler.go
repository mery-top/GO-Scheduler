package internal

import(
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)


type Task func()


//Goroutines
type G struct{
	task Task
	id int
	state string
}

//Processors
type P struct{
	id int
	runQueue chan *G
}

//OS Kernel threads
type M struct{
	id int
	running bool
	boundP *P
}


//Scheduler

type Scheduler struct{
	Ps []*P
	Ms []*M
	globalQueue chan *G
	mu sync.Mutex
	networkPoller chan *G
	blockedG chan *G //Syscalls
	gIDCounter int
}

func NewScheduler(numPs, numMs int) *Scheduler{
	s:= &Scheduler{
		Ps: make([]*P, numPs),
		Ms: make([]*M, numMs),
		globalQueue: make(chan *G,10),
		networkPoller: make(chan *G, 10),
		blockedG: make(chan *G, 10),
	}

	for i:=0; i<numPs; i++{
		s.Ps[i] = &P{
			id: i,
			runQueue: make(chan *G, 10),
		}
	}

	for i:=0; i<numMs; i++{
		s.Ms[i] = &M{
			id: i,
			running: false,
		}
	}

	return s
}

func (s *Scheduler) Go(task Task){
	s.mu.Lock()
	g:= &G{
		task: task,
		id: s.gIDCounter,
		state: "runnable",
	}
	s.gIDCounter++
	s.mu.Unlock()
	s.globalQueue <- g
}

