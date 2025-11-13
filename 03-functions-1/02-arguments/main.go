package main

import "fmt"

var a string

func main() {
	Greet("Alice")
	Greet("Bob")
}

func Greet(a string) {
	fmt.Printf("Hello, %v\n", a)
}
