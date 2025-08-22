package main

import "fmt"

var Name = "Jon"

func main() {
	output := Wrap(Hello)
	fmt.Println(output())
	Name = "Bob"
	fmt.Println(output())
}

func Hello() string {
	return "hello, " + Name
}

func Wrap(stringer func() string) func() string {
	return func() string {
		return "prefix-" + stringer() + "-suffix"
	}
}
