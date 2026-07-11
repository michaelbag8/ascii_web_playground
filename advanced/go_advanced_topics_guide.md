# Advanced Go: Concurrency, Interfaces, Generics & More

A practical guide — for each topic: **why it exists**, **where you'd use it**, **how to use it**, and **common mistakes**.

---

## 1. Concurrency (Goroutines)

**Why it exists**
Go was built at Google for systems that juggle thousands of network calls, I/O waits, and background jobs. Threads in most languages are heavy (megabytes of stack, expensive context switches). Go's answer is the **goroutine** — a function that runs concurrently, managed by the Go runtime, not the OS. Goroutines start with a tiny ~2KB stack that grows as needed, so you can spin up thousands or millions cheaply.

**Where you'd use it**
- Handling many simultaneous HTTP requests
- Doing independent work in parallel (image processing, fan-out API calls)
- Background tasks (logging, cache refresh, metrics)

**How to use it**

```go
func sayHello() {
    fmt.Println("hello from goroutine")
}

func main() {
    go sayHello() // runs concurrently, doesn't block main
    time.Sleep(100 * time.Millisecond) // give it time to run
}
```

The `go` keyword launches a function as a goroutine. `main()` doesn't wait for goroutines — if `main` returns, the program exits, goroutine or not.

**Common mistakes**
- **Not waiting for goroutines to finish.** `time.Sleep` is a hack; use `sync.WaitGroup` instead.
- **Loop variable capture** (pre-Go 1.22): `for _, v := range items { go func(){ use(v) }() }` — all goroutines could see the same `v`. Go 1.22+ fixed this (each iteration gets its own `v`), but on older Go, pass it as a parameter: `go func(v int){ use(v) }(v)`.
- **Data races**: two goroutines reading/writing the same variable without synchronization. Run `go run -race main.go` to catch these.

```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait() // blocks until all 5 call Done()
```

**🧪 Task**
Write a program that launches 10 goroutines, each computing the square of its index (0–9), and prints all 10 results using `sync.WaitGroup` so `main` doesn't exit early. Then deliberately remove the `WaitGroup` and observe how many results print — confirm your understanding of why they're missing/inconsistent.

---

## 2. Channels

**Why it exists**
Goroutines are cheap, but useless if they can't safely communicate. Go's philosophy: *"Don't communicate by sharing memory; share memory by communicating."* Channels are typed pipes that let one goroutine send a value and another receive it — safely, without manual locks.

**Where you'd use it**
- Passing results from a worker goroutine back to the caller
- Signaling "I'm done" or "start now"
- Coordinating pipelines (stage 1 → channel → stage 2)

**How to use it**

```go
ch := make(chan string) // unbuffered channel of strings

go func() {
    ch <- "hello" // send
}()

msg := <-ch // receive (blocks until a value arrives)
fmt.Println(msg)
```

An **unbuffered channel** synchronizes: the sender blocks until a receiver is ready, and vice versa. This makes it a natural hand-off point / rendezvous.

**Common mistakes**
- **Deadlock**: sending on an unbuffered channel with nobody receiving (and vice versa) freezes forever: `fatal error: all goroutines are asleep - deadlock!`
- Forgetting a channel is a *reference type* — you don't need `*chan int`, just pass `chan int`.
- Using a channel when a simple mutex-protected variable would be clearer. Channels are for *handing off* data/ownership, not for protecting shared state you keep reading in place — that's what mutexes are for.

**🧪 Task**
Write a function `squareAsync(n int) <-chan int` that returns a channel, spawns a goroutine to compute `n*n`, sends the result on the channel, and returns immediately (don't wait inside the function). In `main`, call it for 5 numbers and receive+print each result.

---

## 3. Buffered Channels

**Why it exists**
Sometimes you don't want the sender to block until someone's ready to receive — you want a small queue. A buffered channel holds a fixed number of values before blocking the sender.

**How to use it**

```go
ch := make(chan int, 3) // buffer size 3

ch <- 1 // doesn't block, buffer has room
ch <- 2
ch <- 3
// ch <- 4 // this WOULD block — buffer is full

fmt.Println(<-ch) // 1
```

**Where you'd use it**
- Rate-limiting / worker pools (buffer = max in-flight jobs)
- Decoupling producer speed from consumer speed when some slack is fine
- Semaphores: `sem := make(chan struct{}, maxConcurrent)`

```go
sem := make(chan struct{}, 3) // max 3 concurrent
for _, job := range jobs {
    sem <- struct{}{}        // acquire
    go func(j Job) {
        defer func() { <-sem }() // release
        process(j)
    }(job)
}
```

**Common mistakes**
- Treating buffer size as a magic performance fix — a buffer just delays blocking, it doesn't remove backpressure entirely. If consumers are permanently slower than producers, the buffer fills up and you're back to blocking.
- Assuming buffered channels are "safer" than unbuffered — a full buffer still blocks, and unclear buffer sizing can hide bugs that only show under load.

**🧪 Task**
Build a mini "job queue": a buffered channel of capacity 3 holding `int` job IDs. Launch 5 "worker" goroutines that each `<-` a job, sleep for a random short duration, and print `"worker done with job X"`. Push 10 jobs into the channel from `main`. Watch how the buffer limits how many jobs queue up versus how many workers process concurrently.

---

## 4. Closing Channels

**Why it exists**
A channel doesn't know when the sender is "done." Closing a channel is how the sender broadcasts "no more values are coming" to all receivers, cleanly and to any number of listeners at once.

**How to use it**

```go
ch := make(chan int)

go func() {
    for i := 0; i < 3; i++ {
        ch <- i
    }
    close(ch) // signal: done sending
}()

for {
    v, ok := <-ch
    if !ok {
        fmt.Println("channel closed")
        break
    }
    fmt.Println(v)
}
```

`v, ok := <-ch` — `ok` is `false` once the channel is closed *and* drained. Reading a closed channel never blocks; you get the zero value + `ok == false`.

**Common mistakes**
- **Only the sender should close a channel.** Closing from the receiving side, or from multiple goroutines, causes a panic.
- **Sending on a closed channel panics** — `panic: send on closed channel`. Make sure only one goroutine "owns" the close decision.
- **Closing a channel twice panics.** If multiple goroutines might try to close, use `sync.Once`.
- You don't *have* to close every channel — only close when receivers need to know "no more values," e.g., to stop a `range` loop. Channels are garbage collected like anything else once unreferenced.

**🧪 Task**
Write a producer goroutine that sends the numbers 1–20 on a channel, then closes it. In `main`, read from the channel using the `v, ok := <-ch` form in a `for` loop, and print `"channel closed, exiting"` once `ok` is `false`. Then try sending one extra value *after* `close(ch)` and confirm you get a panic.

---

## 5. Range over Channels

**Why it exists**
Manually checking `ok` every time is repetitive. `range` over a channel does it for you — it receives until the channel is closed, then stops automatically.

**How to use it**

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i * i
    }
    close(ch)
}()

for v := range ch {
    fmt.Println(v) // 0, 1, 4, 9, 16
}
// loop exits automatically when ch is closed
```

**Common mistakes**
- If the channel is **never closed**, `range` blocks forever — a silent deadlock/goroutine leak.
- Forgetting that `range` gives you only the value, not the `ok` flag — you can't distinguish "zero value sent" from "channel closed" inside the loop (you just won't get an iteration for the close).

**🧪 Task**
Rewrite your Task from section 4 (producer sending 1–20) to use `range` instead of manually checking `ok`. Then write a second version where you forget to `close()` the channel — run it and confirm the program hangs, and explain in a comment why.

---

## 6. Select

**Why it exists**
Sometimes a goroutine needs to wait on *multiple* channels at once and act on whichever is ready first — like a `switch` for channel operations. Without `select`, you'd need to poll or pick one channel arbitrarily.

**How to use it**

```go
select {
case msg1 := <-ch1:
    fmt.Println("from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("from ch2:", msg2)
case ch3 <- "ping":
    fmt.Println("sent to ch3")
}
```

`select` blocks until *one* of its cases can proceed. If multiple are ready simultaneously, it picks **one at random** (this is intentional — it prevents starvation and stops you from relying on ordering).

**Where you'd use it**
- Timeouts: `select { case res := <-ch: ...; case <-time.After(2*time.Second): ... }`
- Merging multiple input streams into one
- Cancellation via `context.Done()`

```go
select {
case res := <-resultCh:
    fmt.Println("got result:", res)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
}
```

**Common mistakes**
- Expecting deterministic/ordered case selection — it's random among ready cases, by design.
- An empty `select{}` blocks forever (occasionally used intentionally to park `main`, but easy to write by accident).
- Forgetting `time.After` creates a new timer each call — in a loop, this can leak timers; prefer `time.NewTimer` + `Stop()`/`Reset()` in hot loops.

**🧪 Task**
Simulate two "data sources" as goroutines sending to `ch1` and `ch2` at random intervals (e.g., 100–500ms). In `main`, use a `select` in a loop to print whichever arrives first, tagging each print with which channel it came from. Add a `time.After(3 * time.Second)` case that prints `"timed out, stopping"` and exits the loop.

---

## 7. Select with Default (Non-blocking Channels)

**Why it exists**
Normally `select` blocks until a case is ready. Adding `default` makes it **non-blocking**: if no channel is ready *right now*, run `default` instead of waiting.

**How to use it**

```go
select {
case msg := <-ch:
    fmt.Println("received:", msg)
default:
    fmt.Println("no message ready, moving on")
}
```

**Where you'd use it**
- Polling a channel without blocking the rest of your logic
- "Try to send, but don't wait" patterns (e.g., dropping a metric if the buffer's full instead of blocking)

```go
select {
case jobQueue <- job:
    // enqueued
default:
    fmt.Println("queue full, dropping job")
}
```

**Common mistakes**
- Overusing non-blocking selects in a busy loop = **busy-waiting**, burning CPU. If you're polling constantly, you probably want a blocking `select` or a ticker instead.
- Silently dropping data via `default` without realizing it — make sure that's actually the behavior you want, not a bug hiding a full/unready channel.

**🧪 Task**
Create a buffered channel of capacity 2 acting as an "event queue." Write a loop that tries to send 5 events into it using `select` + `default`, printing `"enqueued event N"` on success or `"queue full, dropping event N"` when it falls to `default`. Nothing should ever block.

---

## 8. Channels Review — Mental Model

| Concept | Blocking behavior |
|---|---|
| Unbuffered `chan` | send blocks until receive, and vice versa |
| Buffered `chan` (cap N) | send blocks only when buffer full; receive blocks only when empty |
| `close(ch)` | further sends panic; receives drain remaining values then return zero value, `ok=false` |
| `range ch` | receives until closed, then exits automatically |
| `select` | waits on multiple channel ops, picks a ready one at random |
| `select` + `default` | never blocks — falls to `default` if nothing's ready |

**Rule of thumb:** channels are for *ownership transfer* and *signaling*, not general-purpose shared state. If you just need to protect a variable multiple goroutines read/write in place, reach for a mutex (see below) — it's usually simpler and faster.

**🧪 Task**
Without writing code, go back through your Task solutions for sections 2–7 and, for each, write a one-line note: "unbuffered / buffered / closed / select" — which mental-model row it demonstrates. This is a review exercise to make sure the model is sticking before combining these into bigger programs.

---

## 9. Ping Pong (Classic Channel Pattern)

**Why it exists as an exercise**
"Ping pong" is the canonical toy problem for internalizing how two goroutines can hand control back and forth using channels — a building block for more complex coordination (turn-based systems, pipelines, request/response over channels).

**How to use it**

```go
func main() {
    ping := make(chan struct{})
    pong := make(chan struct{})
    done := make(chan struct{})

    go func() { // "ping" player
        for i := 0; i < 3; i++ {
            <-pong        // wait for pong's turn signal (except first time)
            fmt.Println("ping")
            ping <- struct{}{} // hand off to pong
        }
        close(done)
    }()

    go func() { // "pong" player
        for i := 0; i < 3; i++ {
            pong <- struct{}{} // kick things off / hand off to ping
            <-ping
            fmt.Println("pong")
        }
    }()

    <-done
}
```

The key idea: each goroutine **blocks on receive** until the other sends, creating a strict alternation — a hand-off, not shared state.

**Common mistakes**
- Using buffered channels here defeats the purpose — buffering lets both sides race ahead instead of strictly alternating.
- Forgetting a `done`/exit condition — ping-pong loops with no termination run forever (fine for demos, a leak in real code).

**🧪 Task**
Extend the ping-pong example to 3 goroutines ("ping", "pong", "peng") passing a token around in a ring — each one only prints and passes the token when it's their turn — for 5 full rounds, then all exit cleanly. Use a single `chan struct{}` per "hand-off" (3 channels total, or a ring of channels), not buffered ones.

---

## 10. Interfaces in Go

**Why it exists**
Go doesn't have classical inheritance. Instead, it uses **structural typing**: a type satisfies an interface automatically just by implementing its methods — no `implements` keyword needed. This keeps coupling low: you can write code against behavior ("can this Speak()?") rather than concrete types, and existing types (even ones you didn't write) can satisfy an interface you define later.

**How to use it**

```go
type Speaker interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof" }

type Robot struct{}
func (r Robot) Speak() string { return "Beep" }

func announce(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    announce(Dog{})   // works — Dog has Speak()
    announce(Robot{}) // works — Robot has Speak() too
}
```

Neither `Dog` nor `Robot` declares "I implement Speaker" — they just have the right method, and that's enough.

**Where you'd use it**
- Dependency injection / testing (pass a fake `Speaker` in tests instead of a real one)
- Standard library patterns: `io.Reader`, `io.Writer`, `error`, `sort.Interface`
- Decoupling packages — define the interface where it's *used*, not where it's implemented (Go convention: small interfaces near the consumer)

**Common mistakes**
- **Designing huge interfaces.** Go favors small ones — `io.Reader` is one method. Prefer several small interfaces over one big one; compose if needed (`io.ReadWriter`).
- **Returning interfaces, accepting interfaces** is often backwards advice applied wrong — the actual Go idiom is "accept interfaces, return concrete types" so callers get full information and can still pass your struct wherever an interface is needed.
- Comparing interface values with `==` when the underlying concrete type isn't comparable (e.g., contains a slice) — panics at runtime.
- **The nil interface gotcha:**

```go
type MyError struct{}
func (e *MyError) Error() string { return "boom" }

func doSomething() error {
    var err *MyError = nil
    return err // returns a NON-nil error interface!
}

func main() {
    if err := doSomething(); err != nil {
        fmt.Println("got error!") // this prints, surprisingly
    }
}
```
An interface value is `nil` only if **both** its type and value are nil. Here the interface holds type `*MyError` (non-nil type, nil value), so it's *not* `== nil`. Return `nil` literally, not a typed nil pointer, when you mean "no error."

**🧪 Task**
Define a `Shape` interface with `Area() float64` and `Perimeter() float64`. Implement it for `Rectangle` and `Circle` structs. Write a function `describe(s Shape)` that prints both values. Then write a `TotalArea(shapes []Shape) float64` that sums the area across a mixed slice of rectangles and circles — this is the payoff moment where interfaces let you treat different types uniformly.

---

## 11. Mutexes in Go

**Why it exists**
When goroutines *must* share memory (not just hand it off via channels) — like a shared counter or cache — you need to guarantee only one goroutine touches it at a time. That's what `sync.Mutex` (**mut**ual **ex**clusion lock) provides.

**Why it's called "mutex"**
The word is a contraction of **"mutual exclusion"** — the core guarantee it provides: only one goroutine can hold the lock at a time, so critical sections of code execute *mutually exclusively* from each other. It comes from classic concurrency theory (Dijkstra-era OS/concurrency research), long predating Go — Go's `sync.Mutex` is just an implementation of that decades-old concept.

**How to use it**

```go
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

For read-heavy workloads, `sync.RWMutex` lets multiple readers in simultaneously, but writers still get exclusive access:

```go
var mu sync.RWMutex
mu.RLock()   // many goroutines can hold this at once
defer mu.RUnlock()
// ... read shared data ...
```

**Where you'd use it**
- Shared counters, caches, maps accessed by multiple goroutines (Go maps are **not** safe for concurrent read/write)
- Protecting any in-place mutable state multiple goroutines touch

**Common mistakes**
- **Forgetting to Unlock** — always `defer mu.Unlock()` right after `Lock()` so it fires even on early return/panic.
- **Copying a struct containing a `sync.Mutex`** — this copies the lock state too, breaking mutual exclusion. `go vet` catches this. Pass structs with mutexes by pointer.
- **Locking, then calling a function that also tries to lock the same (non-reentrant) mutex** → deadlock. Go's mutexes are **not reentrant** — the same goroutine locking twice will block forever.
- Holding a lock across a slow operation (network call, channel send) — shrink the critical section to just the shared-state access.
- Reaching for a mutex when a channel would express the intent more clearly (ownership hand-off), or vice versa — reaching for channels for simple protected state when a mutex is simpler and faster.

**🧪 Task**
Build a thread-safe `Counter` struct with `Increment()` and `Value()` methods, protected by a `sync.Mutex`. Launch 100 goroutines that each call `Increment()` once, `Wait()` for them all, then print `Value()` — it should reliably be 100. Run it with `go run -race` to confirm no race. Then comment out the lock/unlock calls and run with `-race` again to see the detector catch it.

---

## 12. Generics in Go

**Why it exists**
Before Go 1.18, writing a function that worked for `int`, `float64`, `string`, etc. meant either duplicating code per type or using `interface{}` and losing type safety (plus needing runtime type assertions). Generics let you write one function/type parameterized over a type, with compile-time type checking.

**How to use it**

```go
// T must support ordering (< > etc.) — constraint from the "constraints"-like built-in
func Max[T int | float64](a, b T) T {
    if a > b {
        return a
    }
    return b
}

fmt.Println(Max(3, 7))       // T inferred as int
fmt.Println(Max(2.5, 1.1))   // T inferred as float64
```

A generic type:

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

var s Stack[string]
s.Push("a")
s.Push("b")
v, _ := s.Pop() // "b"
```

The standard library's `constraints` ideas were folded into the builtin `cmp` and `slices`/`maps` packages (Go 1.21+): `slices.Sort`, `maps.Keys`, etc. are generic under the hood.

**Where you'd use it**
- Generic data structures: stacks, sets, linked lists, trees
- Utility functions: `Map`, `Filter`, `Reduce` over slices of any type
- Reusable, type-safe algorithms (Min/Max, sorting comparators)

**Common mistakes**
- **Overusing generics where an interface would do.** If you just need "something that can `Speak()`," that's an interface problem, not a generics problem. Reach for generics when you need the *same type* preserved through the operation (e.g., a `Stack[T]` should stay a stack of `T`, not `any`).
- Forgetting the constraint — `any` (alias for `interface{}`) means "no constraint," so you lose access to operators like `<`, `+`, etc. unless your constraint interface explicitly permits them.
- Writing overly generic, hard-to-read APIs "just because you can." Generics are a tool for real duplication problems, not a default style.

**🧪 Task**
Write generic `Map[T, U any](s []T, f func(T) U) []U` and `Filter[T any](s []T, keep func(T) bool) []T` functions. Use `Map` to turn `[]int{1,2,3,4}` into their string representations, and `Filter` to keep only even numbers from `[]int{1..10}`. Then build the generic `Stack[T]` from this guide and use it to reverse a slice of strings.

---

## 13. Lack of Enums in Go

**Why Go doesn't have them**
Go deliberately has no dedicated `enum` keyword — it's part of the language's minimalism philosophy (fewer constructs, less magic). Instead, the idiomatic approach is to define a named type based on an integer (or string) and a set of typed constants.

**How to fake it well**

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusClosed
)

func (s Status) String() string {
    switch s {
    case StatusPending:
        return "Pending"
    case StatusActive:
        return "Active"
    case StatusClosed:
        return "Closed"
    default:
        return "Unknown"
    }
}
```

This gives you type safety (`Status` is distinct from plain `int`, so you can't accidentally pass a raw number where a `Status` is expected without a conversion) plus a `String()` method so `fmt.Println` prints something readable instead of `0`, `1`, `2`.

**Where you'd use it**
- Any fixed set of related constants: statuses, directions, log levels, HTTP methods

**Common mistakes**
- **Go enums are not "closed" the way enums are in Java/Rust/Swift** — nothing stops someone from writing `Status(99)`, an out-of-range value. Always handle a `default` case in switches over your "enum" type.
- Forgetting to give a `String()` method — without it, printing/logging shows raw integers, which is a common debugging headache.
- Using plain untyped `int` constants instead of a named type — you lose the compile-time safety that makes this pattern worthwhile.

**🧪 Task**
Model traffic-light states (`Red`, `Yellow`, `Green`) as a named `Light` type with a `String()` method. Write a `Next() Light` method that returns the next state in the cycle `Red → Green → Yellow → Red`. Write a `switch` over `Light` with a `default: panic("unknown light state")` to show you're handling the "not really closed" nature of Go enums.

---

## 14. iota

**Why it exists**
`iota` is Go's mechanism for auto-incrementing constants inside a `const` block, so you don't hand-write `0, 1, 2, 3...` and risk typos or renumbering pain when you insert a new value.

**How to use it**

```go
const (
    A = iota // 0
    B        // 1 (implicitly repeats "= iota")
    C        // 2
    D        // 3
)
```

`iota` resets to `0` at the start of each `const` block and increments by 1 for every constant *spec* (line), whether or not that line explicitly uses `iota`.

**Common patterns**

```go
// Skipping a value
const (
    _  = iota // skip 0
    KB = 1 << (10 * iota) // 1 << 10
    MB                     // 1 << 20
    GB                     // 1 << 30
)

// Bit flags
type Perm int
const (
    Read Perm = 1 << iota // 1
    Write                 // 2
    Execute                // 4
)
// combine: Read | Write == 3
```

**Where you'd use it**
- Enum-like constant sets (see section 13)
- Bit flag / permission systems
- Byte-size constants (KB/MB/GB) as shown above

**Common mistakes**
- **Reordering or inserting a line in the middle of an `iota` block silently changes every value after it** — if these constants are ever serialized/stored (e.g., in a database), this is a real footgun. For persisted values, consider explicit numbers instead of `iota`, or append new values only at the end.
- Forgetting `iota` increments per **line**, not per usage — a line with multiple comma-separated names still only consumes one `iota` step:
```go
const (
    A, B = iota, iota + 10 // A=0, B=10
    C, D                    // C=1, D=11
)
```
- Assuming `iota` works outside `const` blocks — it doesn't; it's meaningless in `var` or regular expressions.

**🧪 Task**
Define byte-size constants `KB`, `MB`, `GB`, `TB` using `iota` and bit-shifting (as shown above). Write a function `Humanize(bytes int64) string` that picks the largest unit that fits and formats accordingly (e.g., `2500000` → `"2.38 MB"`). Then define a `Permission` bit-flag type (`Read`, `Write`, `Execute`) with `iota`, and write a `Has(p Permission, flag Permission) bool` helper using bitwise `&`.

---

## Capstone Project: Put It All Together

Once you've done the per-topic tasks, build this single program — it forces every concept above to work together:

**Build a concurrent task processor:**

1. Define a `Priority` enum (`Low`, `Medium`, `High`) using `iota` + `String()`.
2. Define a `Task` struct: `ID int`, `Priority Priority`, `Payload string`.
3. Define a `Processor` interface with `Process(Task) error`.
4. Write two implementations: `PrintProcessor` (just logs the task) and `FailingProcessor` (fails tasks with odd IDs, for testing error handling).
5. Build a generic `Queue[T]` (push/pop) and use `Queue[Task]` to hold incoming tasks.
6. Spin up a **worker pool** of N goroutines (`N` configurable) pulling from a **buffered channel** of `Task`. Use a `sync.WaitGroup` to know when all workers are done.
7. Protect a shared results map (`map[int]error`, keyed by task ID) with a `sync.Mutex` (or `sync.RWMutex` if you also read it concurrently while processing).
8. Use `select` with a timeout (`time.After`) so the whole pipeline gives up and exits cleanly if it takes longer than N seconds.
9. Feed in 50 tasks with random priorities, close the input channel when done, and print a final summary: how many succeeded, how many failed, grouped by `Priority`.
10. Run the whole thing with `go run -race` and fix anything it flags before considering it done.

This single project touches goroutines, channels (buffered + closing + range + select), interfaces, mutexes, generics, enums, and `iota` — exactly the list you asked to learn, working together the way they would in real code.
