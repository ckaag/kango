package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := welcome("John")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
