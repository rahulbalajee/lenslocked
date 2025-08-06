package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
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
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
