package main

import (
	"errors"
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

// with Methods
type Counter struct {
	Count int
}

func (c *Counter) Addone() {
	c.Count += 2
}
func (c Counter) Report() string {
	return fmt.Sprintf("Counter is %d", c.Count)
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) Birthday() {
	p.Age++
}

func (m Person) Message() string {
	return fmt.Sprintf("Happy Birthday %s, This Year You Are %d Years Old", m.Name, m.Age)
}

type Task struct {
	Title string
	Done  bool
}

func (t *Task) MarkDone() {
	t.Done = true
}

func (t Task) Verify() string {
	if t.Done {
		return fmt.Sprintf("The task: %q has been marked done", t.Title)
	}
	return fmt.Sprintf("The task %q is not done yet.", t.Title)
}

//errors in Go
func parseAge(input int) (int, error){
	if input < 0{
		return 0, errors.New("age can not be negative")
	}
	return input, nil
}

type ValidationError struct{
	Field string
	Reason string
}
//Because *ValidationError implements Error() string, it satisfies the 
// error interface and can be returned wherever error 
func (e *ValidationError) Error() string{
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Reason)
}

func ValidateUsername(name string) error{
	if len(name) < 3{
		return &ValidationError{Field: "username",Reason: "too short"}
	}
	return nil
}

//error task
func divide(a, b float64) (float64, error){
	if b == 0.0{
		return 0.0, errors.New("not divisible by 0")
	}
	return a / b, nil
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
		Age:  440,
	}
	fmt.Println("\033[31m====2026 Birthday===\033[0m")
	person.Birthday()
	fmt.Println(person.Message())
	fmt.Println("\033[32m====2027 Birthday===\033[0m")
	person.Birthday()
	fmt.Println(person.Message())
	fmt.Println("\033[33m====2028 Birthday===\033[0m")
	person.Birthday()
	fmt.Println(person.Message())

	//Task MarkDone Methods
	task := Task{
		Title: "Learning Advanced Go",
		Done:  false,
	}
	fmt.Println("\033[35m====Task Marked Done===\033[0m")
	task.MarkDone()
	fmt.Println(task.Verify())

	//custom error
	fmt.Println("\033[36m====Custom Error Type===\033[0m")
	err := ValidateUsername("ol")
	if err!=nil{
		fmt.Println(err)
	}

	fmt.Println("====Division====")
	m, err := divide(2,0)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(m)
}
