package internal

import(
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)


//Give me some ideas to import structs from models can't use models.P again & again....

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

func (s *Scheduler) Start(){
	for _, m:= range s.Ms{
		go s.RunMachine(m) 
		// Start a goroutine for each M to simulate OS threads. :)
	}
	go s.PollNetwork()
	go s.HandleSysCalls()
}

func(s *Scheduler) PollNetwork(){
	for{
		time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
		s.mu.Lock()
		g:= &G{
			id: s.gIDCounter,
			task: func() {
				fmt.Println("[NetPoll]: Handling network Event")
			},
			state: "runnable",
		}
		s.gIDCounter++
		s.mu.Unlock()
		s.networkPoller <- g
	}
}


func (s *Scheduler) HandleSysCalls(){
	for g:= range s.blockedG{
		time.Sleep(200 * time.Millisecond) //Mimic the syscall delay time 
		g.state = "runnable"
		fmt.Printf("[SyscallReturn]: g[%d] returning from Syscall", g.id)
		s.globalQueue <- g
	}
}

func (s *Scheduler) RunMachine( m *M){
	m.running = true //assign a kernel thread
	m.boundP = s.Ps[m.id% len(s.Ps)] //static round-robin
	p:= m.boundP //m.boundP = *P

	fmt.Printf("M[%d] BOUND to P[%d]", m.id, p.id)

	//ASSIGN GO-ROUTINES

	for{
		var g *G
		select{
			case g = <- p.runQueue:
			case g = <- s.globalQueue:
			default:
				//Steal
				for _, otherP:= range s.Ps{
					if otherP.id != p.id{
						select{
						case g = <- otherP.runQueue:
							
						}
					}
				}
		
		}
	}
}