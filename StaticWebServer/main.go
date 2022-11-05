package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type Page struct {
	Name string
	Body []byte
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func loadPage(name string) (*Page, error) {
	path := "Pages/" + name + ".html"
	if _, err := os.Stat(path); err == nil {
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		return &Page{Name: name, Body: body}, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else {
		return nil, err
	}
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
	pageName := validPath.FindStringSubmatch(req.URL.Path)
	if pageName == nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := loadPage(req.URL.Path[1:])
	if err != nil {
		http.Error(res, "page load error", http.StatusInternalServerError)
	} else if page == nil {
		http.NotFound(res, req)
	} else {
		res.Write(page.Body)
	}
}
