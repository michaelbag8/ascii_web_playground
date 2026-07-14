package main

import (
	"fmt"
	"sync"
)


func worker(id int, wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Println("worker", id, "done")
}

func task(i int, wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Println("student", i, "done")
}

func main(){

	var wg sync.WaitGroup

	for i:= 1; i <=3;i++{
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("all workers finished")


	var w sync.WaitGroup

	for i:= 1; i <=5;i++{
		w.Add(1)
		go task(i, &w)
	}
	w.Wait()
	fmt.Println("all students are done with their assignment")

}


