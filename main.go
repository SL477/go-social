package main

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/SL477/go-social/internal/database"
	"errors"
)

type errorBody struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Add headers
	w.Header().Set("Content-Type", "application/json")

	// Write JSON body
	response, _ := json.Marshal(payload)
	// deal with err ...
	w.Write(response)

	// Write status code
	w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, err error) {
	e := errorBody{
		Error: err.Error(),
	}
	respondWithJSON(w, 200, e)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{}"))*/

	// you can use any compatible type, but let's use our database package's User type for practice
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, errors.New("test error"))
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", testHandler)
	m.HandleFunc("/err", testErrHandler)

	const addr = "localhost:8080"
	srv := http.Server{
		Handler: m,
		Addr: addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()
}