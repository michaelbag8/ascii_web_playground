package main

import (
	"fmt"
	"strings"
)

func Reverse(str string) string {
	reverse := ""
	for _, r := range str {
		reverse = string(r) + reverse
	}
	return reverse
}
func ReverseString(str string) string {
	runes := []rune(str)
	for i, r := 0, len(runes)-1; i < r; i, r = i+1, r-1 {
		runes[i], runes[r] = runes[r], runes[i]
	}
	return string(runes)
}

func reverseEachWord(str string) string {
	words := strings.Fields(str)
	for i, word := range words {
		words[i] = ReverseString(word)
	}
	return strings.Join(words, " ")
}

// pointers
func swap(a, b *int) (*int, *int) {
	temp := a
	a = b
	b = temp
	return a, b
}

func double(n int) int {
	n = n * 2
	return n
}

func doubleByPointer(n *int) int {
	*n = *n * 2
	return *n
}

func applyDiscount(price *float64, percent float64) float64 {
	discount := *price * (percent / 100)
	*price = *price - discount

	return *price
}

//mutating a person account balance through struct

type Account struct {
	Owner   string
	Balance float64
}

func deposit(acc *Account, amount float64) {
	acc.Balance += amount
}

//with Methods
type Counter struct{
	Count int
}

func (c *Counter) Addone() {
	c.Count+=2
}
func (c Counter) Report() string{
	return fmt.Sprintf("Counter is %d", c.Count)
}

type Person struct{
	Name string
	Age int
}

func (p *Person) Birthday(){
	p.Age++
}

func (m Person) Message() string{
	return fmt.Sprintf("Happy Birthday %s, This Year You Are %d Years Old", m.Name, m.Age)
}

func main() {
	word := "hello world"
	fmt.Println(Reverse(word))

	fmt.Println("-----------")
	fmt.Println(reverseEachWord(word))

	a := 5
	b := 3
	c, d := swap(&a, &b)
	fmt.Println(*c, *d)

	val := 40
	fmt.Println(double(val))
	fmt.Println("--------------")
	fmt.Println(val)

	value := 30
	fmt.Println(doubleByPointer(&value))
	fmt.Println("---------")
	fmt.Println(value)

	price := 1000.0
	//p := &price
	percent := 20.0
	fmt.Println("-----Discounted Price----")
	final := applyDiscount(&price, percent)
	fmt.Println(final)
	fmt.Println("----original price---")
	fmt.Println(price)

	// account

	acc := Account{
		Owner:   "Michael Bag",
		Balance: 4000.67,
	}
	fmt.Println("====Before Deposit====")
	fmt.Println(acc.Balance)

	deposit(&acc, 900)
	fmt.Println("====After Deposit====")
	fmt.Println(acc.Balance)

	//methods
	z := Counter{}
	z.Addone()
	z.Addone()
	z.Addone()
	z.Addone()
	fmt.Println(z.Report())


	//Message method
	person := Person{
		Name: "Michael Bag",
		Age: 440,
	}
	fmt.Println("====2026 Birthday===")
	person.Birthday()
	fmt.Println(person.Message())
	fmt.Println("====2027 Birthday===")
	person.Birthday()
	fmt.Println(person.Message())
	fmt.Println("====2028 Birthday===")
	person.Birthday()
	fmt.Println(person.Message())
}
