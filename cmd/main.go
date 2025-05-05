package main

import (
	"fmt"
	"go_Scheduler/internal"
	"runtime"
	"time"
)

func main(){
	runtime.GOMAXPROCS(1) //simulate single P at hardware level
	sched:= internal.NewScheduler(2,3)

	sched.Start()

	for i:=0; i<10 ;i++{
		i:=i
		sched.Go(func() {
			fmt.Printf("GOROUTINE %d running...\n", i)
			time.Sleep(50 * time.Millisecond)
		})
	}
}