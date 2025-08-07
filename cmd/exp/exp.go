package main

import (
	"html/template"
	"os"
)

type User struct {
	Name      string
	Bio       string
	NestedMap NestedMap
}

type NestedMap struct {
	MapEntry map[int]string
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Rahul Balajee",
		Bio:  `<script>alert("PWNED");</script>`,
		NestedMap: NestedMap{
			MapEntry: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
