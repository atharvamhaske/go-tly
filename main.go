package main

import "fmt"

func main() {
	p := new(int)
	fmt.Println(p)
	fmt.Println(*p)

	x := *p

	p1, p2 := &x, &x
	fmt.Println(&x)
	fmt.Println(p1)
	fmt.Println(*p1) //derefernce operator give value at p1
	fmt.Println(&p1)
	fmt.Println(p2)
}
