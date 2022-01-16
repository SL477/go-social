package main

import (
	"net/http"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{}"))
} 

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", testHandler)

	const addr = "localhost:8080"
	srv := http.Server{
		Handler: m,
		Addr: addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()
}