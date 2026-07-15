package main

import (
	font "fmt"
	"time"
	"math/rand"
	//"testing"
)


func f(n int){
	for i:= range 2{
		font.Println(n, ":", i)
		amt := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * amt)
	}
}

func pinger(msg chan string){
	for i := 0; ; i++{
		msg <- "ping"
	}
}

func printer(c <-chan string){
	for {
		font.Println(<- c)
		time.Sleep(time.Second * 1)
	}
}

func ponger(c chan<- string){
	for i:=0; ; i++{
		c <- "pong"
	}
}
func main() {
	var c chan string = make(chan string)
	go pinger(c)
	go ponger(c)
	go printer(c)
	// for i:= range 2{
	// 	go f(i)
	// }
	var input string
	font.Scanln(&input)
	
}

// message <- "hello" // means send "hello"
// msg := <- message // msg should receive "hello"

