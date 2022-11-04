package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Page struct {
	Name string
	Body []byte
}

func loadPage(name string) (*Page, error) {
	body, err := os.ReadFile("Pages/" + name + ".html")
	if err != nil {
		return nil, err
	}
	return &Page{Name: name, Body: body}, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("port argument is required")
		return
	}

	port, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("port must be integer. arg:", args[0])
		return
	}

	http.HandleFunc("/", pageHandler)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

func pageHandler(res http.ResponseWriter, req *http.Request) {
	page, err := loadPage(req.URL.Path[1:])
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
	} else {
		res.Write(page.Body)
	}
}
