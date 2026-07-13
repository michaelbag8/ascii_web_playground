package main

import (
	"fmt"
	"strings"
)

func Reverse(str string) string{
	reverse := ""
	for _, r := range str{
		reverse = string(r) + reverse
	}
	return reverse
}
func ReverseString(str string) string{
	runes := []rune(str)
	for i,r:=0, len(runes)-1; i < r; i, r = i+1, r-1{
		runes[i], runes[r] = runes[r], runes[i]
	}
	return string(runes)
}

func reverseEachWord(str string) string{
	words := strings.Fields(str)
	for i, word := range words{
		words[i] = ReverseString(word)
	}
	return strings.Join(words, " ")
}

//pointers
func swap(a,b *int) (*int, *int){
	temp := a
	a = b
	b = temp
	return a, b
}

func main() {
	word := "hello world"
	fmt.Println(Reverse(word))

	fmt.Println("-----------")
	fmt.Println(reverseEachWord(word))

	a := 5
	b := 3
	c , d := swap(&a, &b)

	fmt.Println(*c,*d)
}