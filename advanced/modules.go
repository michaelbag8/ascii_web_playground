package main

// import "fmt"

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
// func main() {
// 	word := "hello world"
// 	fmt.Println(Reverse(word))

// 	fmt.Println("-----------")
// 	fmt.Println(ReverseString(word))
// }