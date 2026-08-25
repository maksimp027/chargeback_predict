package main

import "fmt"

func Hello() {
	fmt.Println("Hello, World!")
}

func Sum(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func SumAndMultiply(a, b int) (sum int, mult int) {
	sum = a + b
	mult = a * b
	return sum, mult
}

func isChildren(age int) bool {
	return age < 18
}

func Action(a, b int, action func(int, int) int) int {
	return action(a, b)
}


func main() {
	var think bool
	fmt.Printf("Type %T,\n Value %v", think, think)
	think = true
	fmt.Printf("\nType %T,\n Value %v\n", think, think)

	fmt.Printf("Sum: %d\n", Sum(1, 2))

	sum, mult := SumAndMultiply(1, 2)
	fmt.Printf("Sum and Multiply: %d, %d\n", sum, mult)

	Hello()

	fmt.Printf("Action (Sum): %d\n", Action(1, 2, Sum))
	fmt.Printf("Action (Subtract): %d\n", Action(1, 2, Subtract))

	fmt.Printf("Is Children (age 17): %v\n", isChildren(17))
	fmt.Printf("Is Children (age 19): %v\n", isChildren(19))
}
