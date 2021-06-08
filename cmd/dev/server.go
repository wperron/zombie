// Copyright 2021 William Perron. All rights reserved. MIT License.

// Command dev starts up a small webserver used for local testing.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World")
		log.Println("200 OK /")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
