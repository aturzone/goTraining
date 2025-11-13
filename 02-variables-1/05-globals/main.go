package main

import "fmt"

const taxRate = 0.1

var price = 40.0

func main() {
	tax := taxRate * price
	fmt.Println("tax:", tax)
}
