package models

import (
	"sync"
)

//FOR YOUR REFERENCE


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


