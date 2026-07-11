package main

import (
	"fmt"
	"time"
)

func sendEmail(message string) {
	go func() {
		time.Sleep(time.Millisecond * 250)
		fmt.Printf("Email received: '%s'\n", message)
	}()
	fmt.Printf("Email sent: '%s'\n", message)
}


func test(message string) {
	sendEmail(message)
	time.Sleep(time.Millisecond * 500)
	fmt.Println("========================")
}



func send(ch chan int) {
    ch <- 99
}


// func splitIntoTwo(n []int) ([]int, []int){
// 	mid := len(n)/2
// 	return  n[:mid], n[mid:]
// }

func splitAnySlice[T any](s []T) ([]T, []T) {
    mid := len(s)/2
    return s[:mid], s[mid:]
}

func main() {
	test("Hello there Kaladin!")
	test("Hi there Shallan!")
	test("Hey there Dalinar!")


	ch := make(chan int)
    go send(ch)
    fmt.Println(<-ch)

	s := []int{1,2,3,4,5,6,7,8,9,10,11,12}
	fmt.Println(splitAnySlice(s))

	words :=[]string{"a","b","c","d","e","f"}
	fmt.Println(splitAnySlice(words))
}

