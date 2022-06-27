package main

import "fmt"

var log = fmt.Println

func Hello() string {
	return "Hello, world!"
}

func main() {
	log(Hello())
}