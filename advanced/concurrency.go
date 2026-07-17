package main

import (
	"fmt"
	"sync"
)

func printLetters(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 'A'; i <= 'E'; i++ {
		fmt.Println("Letters: ", string(i))
	}

}

func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		fmt.Println("Numbers: ", i)
	}

}

func square(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		fmt.Printf("Square of %d == %d\n", i, i*i)
	}
}

var GFG = 0

func work(wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	m.Lock()
	GFG = GFG + 1
	m.Unlock()

}

func main() {

	
	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(1)
	go printNumbers(&wg)

	wg.Add(1)
	go printLetters(&wg)

	wg.Add(1)
	go square(&wg)

	for i := 0; i < 10; i++ {
		{
			wg.Add(1)
			go work(&wg, &m)
		}
	}

	wg.Wait()
	fmt.Println("main finished")
	fmt.Println("Value of x", GFG)
}
