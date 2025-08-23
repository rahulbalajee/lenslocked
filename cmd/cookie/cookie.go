package main

import (
	"fmt"
	"net/http"
)

const port = ":3005"

func main() {
	cookie := http.Cookie{
		Name:  "Set-C",
		Value: "This is a test cookie",
		//HttpOnly: false,
		Path: "/test-cookie",
	}

	mux := http.ServeMux{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &cookie)
	})

	mux.HandleFunc("/set-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(cookie.Value))
	})

	mux.HandleFunc("/read-cookie", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Set-C")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Cookie read: %v", cookie)
	})

	err := http.ListenAndServe(port, &mux)
	if err != nil {
		panic(err)
	}

}
