# Think Like a Go Engineer: A Complete Course

This course has two halves that work together. **Part 0** is the meta-skill — how a
senior engineer approaches *any* unfamiliar problem, in any language. **Part 1** is
the Go curriculum itself — pointers, errors, and concurrency — taught lesson by
lesson, with a mini task at the end of every lesson so you actually use what you
just read instead of just recognizing it.

Do the mini tasks. Recognizing code is not the same skill as writing it.

---

# PART 0 — How to Think Like an Engineer

Before you touch a keyboard, a senior engineer runs a mental process that looks
almost boring from the outside. It's not cleverness — it's discipline. Four habits
carry almost all of the weight:

1. Break the unknown problem into known pieces
2. Separate what you actually know from what you're assuming
3. Sketch the *shape* of the solution before writing a line of code
4. Build in small verified steps instead of one big leap

## 0.1 Breaking down a problem you've never seen

The instinct of a beginner facing a new problem is to start typing and hope the
solution emerges. The instinct of an engineer is to **not touch the keyboard yet**
and instead ask: *what smaller, already-solved problems is this made of?*

A technique that works almost every time: restate the problem as a pipeline of
verbs. Take "build a service that alerts me when a website goes down."

- **Check** if a URL responds (a known, solvable sub-problem: HTTP request + timeout)
- **Repeat** that check on a schedule (a known sub-problem: loops, tickers)
- **Decide** what "down" means (non-200 status? timeout? N consecutive failures?)
- **Notify** someone (a known sub-problem: send an email/Slack message)
- **Avoid** spamming on every single failure (a known sub-problem: state — was it
  already down last time we checked?)

Notice none of these five verbs is hard on its own. The "new, scary problem" was
just five familiar problems wearing a trench coat. This is the core move: **you are
never actually solving a novel problem, you're composing familiar ones.** If a
piece genuinely has no familiar shape, that piece — and only that piece — is where
you need to go research, prototype, or ask for help. Everything else you already
know how to do.

**Practice habit:** when a task feels overwhelming, write one line per verb before
writing any code. If you can't name the verbs, you don't understand the problem
yet — go find out more, don't start coding.

## 0.2 Identifying what you know vs. what you're assuming

Bugs and wasted afternoons come disproportionately from one source: treating an
assumption as a fact. Engineers get good at literally separating these into two
lists before starting:

**Known facts** — things you have verified: read in documentation, tested in a
REPL, confirmed with a teammate, seen in the actual data.

**Assumptions** — things you believe but have not checked: "the API probably
returns JSON," "this list is probably never empty," "this function is probably
thread-safe."

The habit: for every assumption, ask *"what is the cheapest way to check this
right now?"* Usually the cheapest check is absurdly fast — read one line of docs,
run one command, print one value. Do that check before you build five hours of
code on top of the assumption.

A concrete Go example: you assume `os.ReadFile` returns an error if the file
doesn't exist, so you plan to check `err != nil`. That's correct — but did you
verify whether it returns a partial slice on error, or `nil`? Thirty seconds with
`go doc os.ReadFile` (or writing a 4-line test program) turns the assumption into
a fact, and now your error-handling code is built on solid ground instead of a
guess.

**Practice habit:** before implementing, write two short lists — "I know" and "I'm
assuming." Convert every assumption you can into a known fact with the cheapest
possible check.

## 0.3 Finding the shape of a problem before writing code

"Shape" means: what are the moving parts, what are their types, and how does data
flow between them? You find the shape with a pen, a whiteboard, or plain comments
— never a code editor first.

A reliable method — **data-first design**:

1. What data comes in? What's its type/shape?
2. What data needs to come out? What's its type/shape?
3. What transformation turns the first into the second?
4. Where can that transformation fail, and what should happen when it does?
5. Does anything need to happen *concurrently*, or is this strictly sequential?

Apply it to "count word frequency in a text file":

- In: raw bytes from a file → string
- Out: a mapping from word → count
- Transform: split into words → normalize case/punctuation → tally into a map
- Failure points: file doesn't exist, file is huge (memory), empty file
- Concurrency: not needed for a single file; would matter for 10,000 files

Notice you now know the core data structure (`map[string]int`) and the two risky
edges (missing file, huge file) *before writing a single line of Go*. That is the
entire point of "finding the shape" — it turns coding from an exploratory act into
a transcription act. You're no longer discovering the solution while typing;
you're typing out a solution you already understand.

**Practice habit:** for any nontrivial task, answer the five data-first questions
above in comments or a notebook before opening the editor.

## 0.4 Building incrementally and verifying at every step

The biggest structural mistake junior engineers make is writing the *entire*
solution before running any of it. When it inevitably fails, they now have to
debug a large surface area with many possible fault points at once.

The professional habit is the opposite: make the smallest possible change that
could be wrong, run it, confirm it behaves as expected, then make the next small
change. This is sometimes called "walking skeleton" development — get a thin,
ugly, end-to-end version working first, then flesh it out.

Concretely, for a Go program with goroutines, channels, and error handling, a
sane build order looks like:

1. Write the function signature and a stub that returns zero values. Compile it.
2. Implement the logic with no concurrency and no error handling — a single
   synchronous happy path. Run it against one known input, confirm the output.
3. Add error handling for the *one* failure mode you already know about (from
   §0.2/0.3). Run it against a bad input, confirm it fails the way you expect.
4. Only now introduce goroutines/channels, if the problem actually needs them.
   Test with a tiny, controlled number of goroutines (2, not 200) so you can
   reason about the output by hand.
5. Scale up (more goroutines, bigger input, add timeouts) once the small case is
   provably correct.

Each step above should end with you running the code and looking at real output,
not just reading it and assuming it's right. "It compiles" is not verification.
"I ran it and the output matched what I predicted on paper" is verification.

**Practice habit:** never write more than ~10–15 lines without compiling and
running them. If you can't verify a piece in isolation, that's a sign it's too
big — break it down further (back to §0.1).

---

You'll use all four of these habits explicitly in the mini tasks below — each one
will ask you to name the sub-problems, list your assumptions, sketch the data
shape, and build in steps, *before* asking you to write the final code.

---

# PART 1 — The Go Curriculum

## Lesson 1 — Pointers

### The idea

Every variable lives at some address in memory. Normally you only ever touch the
*value* stored there. A pointer is a variable that stores an *address* instead of
a value — it lets you refer to "the box" rather than "what's currently in the
box."

Two operators do all the work:

- `&x` — "give me the address of x" (produces a pointer)
- `*p` — "give me the value stored at the address p points to" (dereferences a
  pointer)

```go
package main

import "fmt"

func main() {
    age := 30          // a normal int variable
    ptr := &age        // ptr holds the ADDRESS of age, type is *int

    fmt.Println("value:", age)     // 30
    fmt.Println("address:", ptr)   // something like 0xc0000140a0
    fmt.Println("deref:", *ptr)    // 30 — following the pointer back to the value

    *ptr = 31           // write THROUGH the pointer
    fmt.Println("age is now:", age) // 31 — the original variable changed!
}
```

That last line is the entire reason pointers matter: `*ptr = 31` reached back and
changed `age` itself, even though we only touched `ptr`. Two variables, `age` and
`ptr`, but only one memory location involved.

### Why bother — the mental model

Think of a variable as a labelled locker, and its value as what's inside. A
pointer isn't a second locker with a copy of the contents — it's a note that says
"the thing you want is in locker #4291." Handing someone that note is cheap (an
address is just a number), and if they scribble on what's in locker #4291,
*everyone* holding that note sees the change, because there's still only one
locker.

Compare that to handing someone the value directly: `age` — that copies the
number into their hands. They can scribble on their copy all they like; your
locker is untouched.

### A second pointer type example

```go
package main

import "fmt"

func main() {
    var name string = "Ada"
    var p *string        // zero value of a pointer is nil — points at nothing yet
    fmt.Println(p)        // <nil>

    p = &name
    fmt.Println(*p)       // Ada
}
```

Always be careful with a `nil` pointer: dereferencing it (`*p` when `p` is `nil`)
crashes your program with a runtime panic. A `nil` pointer is a promise of an
address, not an address — check for `nil` before dereferencing if there's any
chance the pointer wasn't set.

### Mini Task 1

Apply Part 0's habits explicitly, in writing, before coding:

1. **Shape it**: write down the "in / out / transform" for this task: *given two
   int variables `a` and `b`, swap their values using pointers, not a third
   variable at the caller's call site.*
2. **Assumptions**: what do you know for certain about how Go pointers behave
   when you write `*ptr = newValue`? What are you merely assuming? Verify one
   assumption by writing a 5-line test program.
3. **Build incrementally**: first write a function `swap(a, b *int)` that does
   nothing but print the two dereferenced values (step 1 — compile, run,
   verify). Then add the actual swap logic (step 2 — run, verify with
   `fmt.Println` before and after).

---

## Lesson 2 — Pointers and Functions

### Passing a pointer as an argument

Go is a **pass-by-value** language: when you call a function, every argument is
*copied*. If you pass an `int`, the function gets its own copy and any change it
makes is invisible to the caller. If you pass a *pointer*, the function gets its
own copy of the address — but that address still points at the caller's original
data, so writes through it are visible outside the function.

```go
package main

import "fmt"

func double(n *int) {
    *n = *n * 2   // write through the pointer
}

func main() {
    value := 21
    double(&value)
    fmt.Println(value) // 42
}
```

Contrast this with the value-passing version, which cannot affect the caller:

```go
func doubleByValue(n int) {
    n = n * 2       // only changes the local copy
}

func main() {
    value := 21
    doubleByValue(value)
    fmt.Println(value) // still 21
}
```

This distinction — "call by value" vs. "call by reference" (using a pointer) — is
one of the most important mental models in Go. Slices, maps, and channels already
behave somewhat reference-like internally, but plain values (`int`, `string`,
`struct`, arrays) are fully copied unless you explicitly pass a pointer.

### Returning a pointer from a function

You can also hand a pointer *back* to the caller. This is safe in Go (unlike C)
because Go's garbage collector keeps the underlying memory alive for as long as
anyone still holds a pointer to it — even after the function that created it has
returned.

```go
package main

import "fmt"

func newCounter() *int {
    count := 0
    return &count   // perfectly safe in Go
}

func main() {
    c := newCounter()
    *c++
    *c++
    fmt.Println(*c) // 2
}
```

### When to actually use pointer arguments

A common beginner mistake is pointer-ifying everything "to be safe." A better
rule of thumb:

- Use a pointer argument when the function needs to **mutate** the caller's data,
  or when the value is **large** and copying it would be wasteful (e.g. a struct
  with many fields, copied millions of times in a hot loop).
- Skip the pointer for small, read-only values (`int`, `bool`, small structs) —
  copying them is cheap and the code is easier to reason about, because you know
  the function *can't* have side effects on your data.

### Mini Task 2

1. **Break it down**: the task is *"write a function `applyDiscount` that takes a
   price and a percentage and reduces the price in place."* Name the sub-problem
   this reduces to (you already solved a nearly identical one in Lesson 1/2).
2. **Assumptions**: will you use a pointer parameter or a return value to hand
   the new price back? Write one sentence justifying your choice using the "when
   to use pointer arguments" rule above.
3. **Build it in steps**: (a) write the signature and a stub, compile; (b) hard
   -code a single test call and print the before/after price; (c) implement the
   real math.

---

## Lesson 3 — Pointers to Structs

### The idea

Structs bundle related fields together. Like any other value, a struct is copied
in full when passed around — which matters a lot once a struct has many fields.
Pointers solve this the same way they did for plain variables.

```go
package main

import "fmt"

type Account struct {
    Owner   string
    Balance float64
}

func main() {
    acc := Account{Owner: "Maria", Balance: 100.0}

    ptr := &acc
    fmt.Println(ptr)        // &{Maria 100}
    fmt.Println(ptr.Owner)  // Maria — Go auto-dereferences for field access
}
```

Notice `ptr.Owner`, not `(*ptr).Owner`. Go lets you use the dot operator directly
on a struct pointer as a convenience — it's automatically shorthand for
dereferencing first. Both forms work; everyone uses the short one.

### Mutating a struct through a pointer

```go
package main

import "fmt"

type Account struct {
    Owner   string
    Balance float64
}

func deposit(acc *Account, amount float64) {
    acc.Balance += amount   // mutates the caller's struct
}

func main() {
    acc := Account{Owner: "Maria", Balance: 100.0}
    deposit(&acc, 50)
    fmt.Println(acc.Balance) // 150
}
```

If `deposit` had taken `acc Account` (no pointer), it would have received a full
copy, changed the copy's balance, and the caller's `acc.Balance` would still be
100. This is such a common source of confusion that it's worth internalizing
fully: **methods and functions that need to modify a struct must take a pointer
receiver/parameter.**

### Pointer receivers on methods

This same idea extends to methods:

```go
package main

import "fmt"

type Counter struct {
    Count int
}

// pointer receiver — mutates the real struct
func (c *Counter) Increment() {
    c.Count++
}

// value receiver — operates on a copy, cannot mutate
func (c Counter) Report() string {
    return fmt.Sprintf("count is %d", c.Count)
}

func main() {
    c := Counter{}
    c.Increment()
    c.Increment()
    fmt.Println(c.Report()) // count is 2
}
```

Go quietly takes the address of `c` for you when you call `c.Increment()`, even
though `c` itself isn't a pointer — this works as long as `c` is an addressable
variable. Rule of thumb: if *any* method on a type needs a pointer receiver
(to mutate, or to avoid copying a large struct), make *all* its methods pointer
receivers, for consistency.

### Mini Task 3

1. **Shape it**: design a `Task` struct with `Title string` and `Done bool`.
   Write, in comments, the in/out/transform for a method `MarkDone()` that flips
   `Done` to `true`.
2. **Assumptions**: will `MarkDone` need a pointer receiver or a value receiver?
   Justify it in one sentence before coding.
3. **Build incrementally**: implement the struct and stub method first (compile),
   then implement the logic, then write a `main()` that creates a `Task`, calls
   `MarkDone()`, and prints the struct to *verify* the mutation actually
   happened — don't just assume it worked because it compiled.

---

## Lesson 4 — Go Errors

### The idea

Go has no exceptions and no `try`/`catch`. Instead, functions that can fail
simply **return an error value as their last return value**, and the caller is
expected to check it immediately. This is a deliberate design choice: it forces
error handling to be visible in the code, right next to the call that might fail,
instead of hidden away in a distant catch block.

The built-in `error` type is just an interface:

```go
type error interface {
    Error() string
}
```

Anything with an `Error() string` method satisfies it. `nil` means "no error."

### Creating simple errors

Two standard ways to produce an error value:

**`errors.New`** — a fixed string message:

```go
package main

import (
    "errors"
    "fmt"
)

func parseAge(input int) (int, error) {
    if input < 0 {
        return 0, errors.New("age cannot be negative")
    }
    return input, nil
}

func main() {
    age, err := parseAge(-5)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println("age:", age)
}
```

**`fmt.Errorf`** — a formatted message, letting you inject values:

```go
package main

import "fmt"

func parseAge(input int) (int, error) {
    if input < 0 {
        return 0, fmt.Errorf("invalid age %d: cannot be negative", input)
    }
    return input, nil
}

func main() {
    if _, err := parseAge(-5); err != nil {
        fmt.Println(err) // invalid age -5: cannot be negative
    }
}
```

### The standard pattern: check immediately

The idiom you will type hundreds of times in Go:

```go
result, err := someFunction()
if err != nil {
    // handle it: log it, return it, wrap it — but don't ignore it
    return err
}
// safe to use result here
```

Never bury error checks far away from the call, and never silently discard an
`error` with `_` unless you have a specific, documented reason (e.g. you already
know it can't fail in this exact context).

### Custom error types

For richer errors — ones that carry structured data, not just a string — define
your own type and implement `Error() string` on it:

```go
package main

import "fmt"

type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Reason)
}

func validateUsername(name string) error {
    if len(name) < 3 {
        return &ValidationError{Field: "username", Reason: "too short"}
    }
    return nil
}

func main() {
    err := validateUsername("al")
    if err != nil {
        fmt.Println(err)
    }
}
```

Because `*ValidationError` implements `Error() string`, it satisfies the `error`
interface and can be returned wherever `error` is expected. Callers who need the
extra structured fields can type-assert it back:

```go
if verr, ok := err.(*ValidationError); ok {
    fmt.Println("bad field:", verr.Field)
}
```

### Mini Task 4

1. **Break it down**: the task is *"write a function `divide(a, b float64)
   (float64, error)` that returns an error instead of crashing on
   division by zero."* Name the two branches this splits into.
2. **Assumptions vs. facts**: what happens in Go if you divide a float by 0.0
   (not int)? Is it a panic, or a special value like `+Inf`? Don't guess — write
   a 3-line test program to find out before deciding whether you even need an
   error check for the float case.
3. **Build incrementally**: implement the happy path first and verify with one
   normal call; then add the error branch and verify with a call you expect to
   fail, checking that `err != nil` and the message is what you intended.

---

## Lesson 5 — defer, panic, and recover

### `defer` — delay a call until the function returns

`defer` schedules a function call to run right before the surrounding function
returns, no matter how it returns (normal return, or even after a panic). It's
most commonly used for cleanup: closing files, unlocking mutexes, closing
database connections.

```go
package main

import "fmt"

func main() {
    defer fmt.Println("this runs last")
    fmt.Println("this runs first")
    fmt.Println("this runs second")
}
// Output:
// this runs first
// this runs second
// this runs last
```

### Multiple defers run LIFO

If you stack several `defer` statements, they unwind in **last-in, first-out**
order — like a stack of plates:

```go
package main

import "fmt"

func main() {
    defer fmt.Println("closed database")
    defer fmt.Println("closed file")
    defer fmt.Println("released lock")
    fmt.Println("doing work...")
}
// Output:
// doing work...
// released lock
// closed file
// closed database
```

This ordering matters: resources are released in the reverse order they were
acquired, which is exactly the order you usually want (release the most
recently-acquired thing first).

A realistic pattern:

```go
func readConfig(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close() // guaranteed to run even if code below panics or returns early

    // ... read from f ...
    return nil
}
```

### `panic` — stop everything, right now

`panic` immediately halts normal execution of the current function, runs any
deferred calls on its way out, and then propagates up the call stack, crashing
the program if nothing intercepts it. Reserve it for truly unrecoverable
situations — a corrupted invariant, a programmer error — not for ordinary,
expected failures (those should be `error` values, per Lesson 4).

```go
package main

import "fmt"

func mustPositive(n int) int {
    if n < 0 {
        panic(fmt.Sprintf("mustPositive: got negative value %d", n))
    }
    return n
}

func main() {
    fmt.Println(mustPositive(5))
    fmt.Println(mustPositive(-1)) // program crashes here
    fmt.Println("never reached")
}
```

### `recover` — catch a panic and keep going

`recover` only does anything useful when called directly inside a deferred
function. It stops the panic from propagating further and returns the value that
was passed to `panic`.

```go
package main

import "fmt"

func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()

    result = a / b // panics if b == 0 (integer division)
    return result, nil
}

func main() {
    r, err := safeDivide(10, 2)
    fmt.Println(r, err) // 5 <nil>

    r, err = safeDivide(10, 0)
    fmt.Println(r, err) // 0 recovered from panic: runtime error: integer divide by zero
}
```

Notice the named return values (`result int, err error`) — the deferred function
assigns to `err` directly, and that's what the caller sees after recovery.

**Design guidance**: `recover` is mostly useful at boundaries — e.g. a web server
recovering per-request so that one bad request panic doesn't crash the whole
server, while still returning a normal `error` from your own APIs. Don't reach
for `panic`/`recover` as a general substitute for `error` — it obscures control
flow and most Go style guides treat routine `panic` use as a smell.

### Mini Task 5

1. **Shape it**: design a function `processItems(items []int) (processed int,
   err error)` that panics internally if it hits a negative number (simulate a
   "should never happen" invariant violation), but the function should recover
   and report a normal error to its caller instead of crashing the whole
   program.
2. **Assumptions**: will `recover()` see anything if you call it *outside* of a
   deferred function? Verify this by trying it and observing the result, rather
   than trusting your memory.
3. **Build incrementally**: (a) write the loop with no panic/recover, verify it
   works on all-positive input; (b) add the deliberate panic on negative input,
   verify the crash happens; (c) add `defer`+`recover`, verify the function now
   returns an error instead of crashing.

---

## Lesson 6 — Goroutines

### The idea

A goroutine is a function that runs concurrently with the rest of your program.
You start one by putting `go` in front of a function call. Goroutines are
extremely cheap compared to OS threads — a Go program can comfortably run tens of
thousands of them.

```go
package main

import "fmt"

func greet(name string) {
    fmt.Println("hello,", name)
}

func main() {
    go greet("goroutine")  // scheduled to run concurrently
    greet("main")          // runs directly

    // main() may exit before the goroutine above even gets scheduled!
}
```

### The trap: `main()` doesn't wait

The single most common beginner mistake with goroutines: `main()` is itself the
program's initial goroutine, and when it returns, the *entire program* exits —
even if other goroutines haven't finished, or haven't even started yet.

```go
package main

import "fmt"

func main() {
    go fmt.Println("might never print!")
    // main() returns almost immediately — the goroutine above
    // frequently doesn't get a chance to run at all.
}
```

A quick (but crude) fix is `time.Sleep`, just to *see* the concurrency happening:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go fmt.Println("now it usually prints")
    time.Sleep(100 * time.Millisecond) // give the goroutine a chance to run
}
```

`time.Sleep` is a teaching tool, not a real solution — real Go code uses
synchronization primitives (`sync.WaitGroup`, or channels — see Lesson 7) to wait
for goroutines deterministically instead of guessing at a sleep duration.

### The correct fix: `sync.WaitGroup`

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // signal completion when this function returns
    fmt.Println("worker", id, "done")
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)         // register one goroutine we're waiting on
        go worker(i, &wg)
    }

    wg.Wait() // blocks until all three call wg.Done()
    fmt.Println("all workers finished")
}
```

`wg.Add(1)` before starting each goroutine, `wg.Done()` (usually via `defer`)
inside the goroutine when it finishes, `wg.Wait()` in the caller to block until
the count returns to zero. This is the idiomatic replacement for "sleep and
hope."

### Order is not guaranteed

When multiple goroutines run, the Go scheduler interleaves them in an order that
is *not* predictable and can vary between runs. Never write code that depends on
goroutines finishing in a particular sequence unless you've explicitly
synchronized that ordering yourself.

### Mini Task 6

1. **Break it down**: the task is *"launch 5 goroutines, each printing its own
   ID, and make sure the program doesn't exit before all 5 have printed."* Name
   the two sub-problems (starting N goroutines; waiting for all N to finish).
2. **Assumptions**: before running anything, predict on paper whether the
   printed order will be 1,2,3,4,5. Then run it several times and check whether
   your prediction held. What does that tell you about relying on goroutine
   order?
3. **Build incrementally**: (a) write the loop starting goroutines with no
   synchronization, run it a few times, observe the flaky/missing output; (b)
   add a `sync.WaitGroup` properly; (c) run it 5 times in a row and confirm all
   5 IDs print every time.

---

## Lesson 7 — Channels

### The idea

Channels are typed pipes that goroutines use to send and receive values safely,
without you managing locks by hand. Create one with `make`:

```go
ch := make(chan int) // an unbuffered channel of ints
```

Send with `ch <- value`, receive with `<-ch`.

```go
package main

import "fmt"

func produce(ch chan int) {
    ch <- 42 // send 42 into the channel
}

func main() {
    ch := make(chan int)
    go produce(ch)

    value := <-ch // receive from the channel
    fmt.Println(value) // 42
}
```

### Unbuffered channels block — and that's the point

An unbuffered channel has no internal storage: a send blocks until some goroutine
is ready to receive, and a receive blocks until some goroutine is ready to send.
This makes channels a natural synchronization tool, not just a data pipe — the
handoff itself guarantees both sides were "present" at the same moment.

```go
package main

import "fmt"

func sender(ch chan string) {
    ch <- "message sent"           // blocks here until received
    fmt.Println("sender: handoff complete")
}

func main() {
    ch := make(chan string)
    go sender(ch)

    fmt.Println("main: waiting for message")
    msg := <-ch
    fmt.Println("main: received:", msg)
}
```

### Buffered channels: some slack in the pipe

`make(chan int, 3)` creates a channel with room for 3 values before a send
blocks. This is useful when you want a producer to get a small head start
without waiting on a consumer for every single item.

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    ch <- 1 // doesn't block — buffer has room
    ch <- 2 // doesn't block — buffer now full
    // ch <- 3 would block here until something is received

    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
}
```

### Closing a channel and ranging over it

A sender can `close(ch)` to signal "no more values are coming." Receivers can
detect this, and `for range` over a channel automatically stops when it's closed
and drained:

```go
package main

import "fmt"

func generate(ch chan int) {
    for i := 1; i <= 3; i++ {
        ch <- i
    }
    close(ch) // signal completion
}

func main() {
    ch := make(chan int)
    go generate(ch)

    for v := range ch { // stops automatically once ch is closed and empty
        fmt.Println("got:", v)
    }
    fmt.Println("channel closed, loop ended cleanly")
}
```

Only the *sender* should close a channel, and never send on a closed channel — 
doing so panics.

### Mini Task 7

1. **Shape it**: design a pipeline where one goroutine computes the squares of
   1..5 and sends each result on a channel, and `main()` receives and prints
   each one. Write the in/out/transform before coding: what type is the
   channel, what closes it, and when does `main()` know to stop reading?
2. **Assumptions**: what happens if you forget to `close(ch)` and use `for range`
   to consume it? Predict the behavior, then intentionally remove the `close`
   call and run it to confirm (you should see it hang — that's the lesson).
3. **Build incrementally**: (a) get one hardcoded value flowing end-to-end
   through the channel and printed, verify; (b) generalize to the loop of 5
   values; (c) add `close` and switch `main()` to `for range`, verify clean
   termination.

---

## Lesson 8 — Select

### The idea

`select` lets a goroutine wait on multiple channel operations at once and
proceeds with whichever one becomes ready first — structurally similar to a
`switch`, but each `case` is a channel send or receive rather than a value
comparison.

```go
package main

import "fmt"

func main() {
    numCh := make(chan int)
    strCh := make(chan string)

    go func() { numCh <- 7 }()
    go func() { strCh <- "seven" }()

    select {
    case n := <-numCh:
        fmt.Println("got a number:", n)
    case s := <-strCh:
        fmt.Println("got a string:", s)
    }
}
```

If both channels happen to be ready at roughly the same moment, `select` picks
one of the ready cases at random — you cannot assume ordering between cases.

### `select` naturally handles "whichever is ready first"

This makes `select` the natural tool for timeouts and for consuming from
whichever of several producers responds first:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    result := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        result <- "slow operation finished"
    }()

    select {
    case r := <-result:
        fmt.Println(r)
    case <-time.After(1 * time.Second): // fires after 1s if result isn't ready
        fmt.Println("timed out waiting for result")
    }
}
```

`time.After` returns a channel that receives a value once the given duration has
elapsed — a clean, idiomatic way to implement "wait up to N seconds, then give
up."

### The `default` case: don't block at all

Adding a `default` case makes `select` non-blocking: if no other case is
immediately ready, `default` runs instantly instead of waiting.

```go
package main

import "fmt"

func main() {
    ch := make(chan int)

    select {
    case v := <-ch:
        fmt.Println("received", v)
    default:
        fmt.Println("nothing ready right now, moving on")
    }
}
```

This is useful for polling a channel without stalling the rest of your program
while you wait.

### Mini Task 8

1. **Break it down**: the task is *"query two mock services concurrently and use
   whichever answers first, but give up after 500ms if neither answers."* Name
   the three channel operations this decomposes into (service A, service B, the
   timeout).
2. **Assumptions**: if both services would respond in exactly the same
   microsecond (contrived, but imagine it), what does `select` do? Look it up or
   test it rather than guessing — this matters for reasoning about race
   conditions later.
3. **Build incrementally**: (a) implement just service A responding, with
   `select` having only one real case plus `default`, verify; (b) add service B
   as a second case, verify either can "win" across multiple runs; (c) add the
   `time.After` timeout case and verify it fires when you make both services
   artificially slow.

---

# Where to go from here

You now have both halves: the engineering-thinking habits from Part 0, and eight
concrete Go lessons to apply them to. A good next exercise is to combine several
lessons into one small program without a prompt telling you which concepts to
use — e.g., "build a worker pool of 3 goroutines that pull jobs off a channel,
process them, return errors on a results channel, and time out if a job takes too
long." Before writing any code, run through the four Part 0 habits explicitly on
paper. That's the actual skill this course is trying to build.
