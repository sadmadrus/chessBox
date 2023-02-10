package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func validateRequestMethod(r *http.Request, allowedMethod string) error {
	if r.Method != allowedMethod {
		return errors.New(fmt.Sprintf("Method %s is forbidden!", r.Method))
	}
	return nil

}

func home(w http.ResponseWriter, r *http.Request) {
	about := "This is a chessBox sessions server"
	w.Write([]byte(about))

}

func find(w http.ResponseWriter, r *http.Request) {
	err := validateRequestMethod(r, http.MethodGet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	params := r.URL.Query()
	for k, v := range params {
		log.Println(k, v)
	}

}

func add(w http.ResponseWriter, r *http.Request) {
	err := validateRequestMethod(r, http.MethodPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	r.ParseForm()
	for k, v := range r.Form {
		log.Println(k, v)
	}

}

func main() {
	log.Println("Initialize session manager MS")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/find", find)
	mux.HandleFunc("/add", add)

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatal(err)
	}

}
