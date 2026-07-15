package main

import(
	"fmt"
	"os"
)


func menu(){
	
	for {
	var ans string

	fmt.Println("1. Admin")
	fmt.Println("2. Student")
	fmt.Println("3. Examiner")
	fmt.Println("4. Quit")
	fmt.Scanf("%s", &ans)

	switch ans {
	case "1":
		fmt.Println("==Welcome Admin==")
	case "2":
		fmt.Println("==Welcome Student==")
	case "3":
		fmt.Println("==Welcome Examiner==")
	case "4":
		fmt.Println("==Goodbye==")
		os.Exit(1)
		
	}
	}
	
}

func main() {
	fmt.Println("======Welcome To The Exam======")
	menu()
}
