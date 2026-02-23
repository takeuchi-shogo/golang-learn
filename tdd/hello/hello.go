package main

import "fmt"

const helloPrefiix = "hello"

func Hello(name string) string {
	return helloPrefiix + name
}

func main() {
	fmt.Println(Hello("world"))
}
