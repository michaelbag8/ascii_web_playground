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

func square(wg *sync.WaitGroup){
	defer wg.Done()
	for i:=1; i<=10; i++{
		fmt.Printf("Square of %d == %d\n", i, i*i)
	}
}
func main() {
	var wg sync.WaitGroup
	
	wg.Add(3)
	go printNumbers(&wg)
	go printLetters(&wg)
	go square(&wg)

	wg.Wait()
	fmt.Println("main finished")
}



