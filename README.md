# 🌀 Go Scheduler Simulator

> 🚀 This project emulates Go's runtime scheduler and models how Go manages goroutines using its M:N scheduler.

## 🎯 Focus Areas

This project simulates and demonstrates:

- ⚙️ **M:N Scheduling Model** – Mapping many goroutines to fewer OS threads
- 🔄 **Goroutine Lifecycle Management** – From creation to termination
- 🔀 **Work-Stealing Mechanisms** – Idle processors steal work from others
- ⏱️ **Preemptive Scheduling** – Prevents long-running goroutines from hogging execution
- 🔧 **System Call Handling** – How goroutines pause and resume after syscalls
- 🌐 **Network Polling Simulation** – Simulated I/O wait using network pollers

## Go Source Code

You can view the full `proc.go` source file from the Go runtime here:  
[Go runtime/proc.go](https://github.com/golang/go/blob/master/src/runtime/proc.go)

----

  ```
                            ┌────────────┐
                            │ Scheduler  │
                            └────┬───────┘
                                 │
                      ┌──────────┼────────────┐
                      ▼          ▼            ▼
                   ┌─────┐    ┌─────┐      ┌─────┐
                   │  P0 │    │  P1 │  ... │  Pn │   <- Processors (with local runQueue)
                   └─┬───┘    └─┬───┘      └─┬───┘
                     │          │            │
                     ▼          ▼            ▼
                 [G, G, G]   [G, G]        [G]
  
  Each Processor P is bound to a Machine M (Kernel Threads):
  
      M0 ──> P0
      M1 ──> P1
      M2 ──> P0  (shared binding via modulo)

  Mutex for New Go Routines and Network Pollers
  
  Each Machine M pulls tasks from:
      1. Local run queue (P)
      2. Global queue
      3. Other P’s run queues (stealing)
      4. Network poller
      5. Premptive Scheduling
