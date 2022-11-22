package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)
	fmt.Printf("serving on %s", address)

	http.HandleFunc("/", helloServer)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error executing: %s", err)
		os.Exit(1)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])

}
