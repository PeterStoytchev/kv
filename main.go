package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	dir string
)

func create(w http.ResponseWriter, req *http.Request) {
	id := uuid.New().String()

	f, err := os.Create(id)
	defer f.Close()

	if err != nil {
		fmt.Fprintf(w, "500 Internal Server Error\n%s\n", err)
		return
	}

	io.Copy(f, req.Body)
	fmt.Fprintf(w, "200 OK\n%s\n", id)
}

func get(w http.ResponseWriter, req *http.Request) {
	split := strings.Split(req.URL.Path, "/")
	if len(split) < 3 {
		fmt.Fprintf(w, "500 Internal Server Error\n")
		return
	}

	id := split[2]

	f, err := os.Open(id)
	defer f.Close()

	if err != nil {
		fmt.Fprintf(w, "500 Internal Server Error\n%s\n", err)
		return
	}

	io.Copy(w, f)
	fmt.Fprintf(w, "\n200 OK\n")
}

func patch(w http.ResponseWriter, req *http.Request) {
	split := strings.Split(req.URL.Path, "/")
	if len(split) < 3 {
		fmt.Fprintf(w, "500 Internal Server Error\n")
		return
	}

	id := split[2]
	f, err := os.Create(id)

	defer f.Close()
	if err != nil {
		fmt.Fprintf(w, "500 Internal Server Error\n%s\n", err)
		return
	}

	io.Copy(f, req.Body)
	fmt.Fprintf(w, "200 OK\n")
}

func delete(w http.ResponseWriter, req *http.Request) {
	split := strings.Split(req.URL.Path, "/")
	if len(split) < 3 {
		fmt.Fprintf(w, "500 Internal Server Error\n")
		return
	}

	id := split[2]
	os.Remove(id)

	fmt.Fprintf(w, "200 OK\n")
}

func main() {
	if len(os.Args) < 3 {
		panic("Not enough arugments")
	}

	dir = os.Args[2]

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	os.Chdir(dir)

	http.HandleFunc("/create", create)
	http.HandleFunc("/get/", get)
	http.HandleFunc("/patch/", patch)
	http.HandleFunc("/delete/", delete)

	http.ListenAndServe(os.Args[1], nil)
}
