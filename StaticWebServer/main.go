package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
)

type Page struct {
	Name string
	Body []byte
}

var validPagePath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func loadFile(name string) ([]byte, error) {
	if _, err := os.Stat(name); err == nil {
		body, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}

		return body, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else {
		return nil, err
	}
}

func loadPage(name string) (*Page, error) {
	body, err := loadFile(path.Join("Pages", fmt.Sprintf("%s.html", name)))
	if err != nil {
		return nil, err
	} else if body == nil {
		return nil, nil
	}

	return &Page{Name: name, Body: body}, nil
}

func loadData(name string) (any, error) {
	content, err := loadFile(path.Join("Pages", fmt.Sprintf("%s.json", name)))
	if err != nil {
		return nil, err
	} else if content == nil {
		return nil, nil
	}

	var data any
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
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
	http.HandleFunc("/css/", cssHandler)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

func cssHandler(res http.ResponseWriter, req *http.Request) {
	var cssName = req.URL.Path[len("/css/"):]
	if cssName == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	content, err := loadFile(path.Join("Pages", "css", cssName))
	if err != nil {
		http.Error(res, "css load error", http.StatusInternalServerError)
	} else if content == nil {
		http.NotFound(res, req)
	}

	res.Header().Set("Content-Type", "text/css; charset=utf-8")
	res.Write(content)
}

func pageHandler(res http.ResponseWriter, req *http.Request) {
	pageNameMatch := validPagePath.FindStringSubmatch(req.URL.Path)
	if pageNameMatch == nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	pageName := pageNameMatch[1]
	page, err := loadPage(pageName)
	if err != nil {
		http.Error(res, "page load error", http.StatusInternalServerError)
	} else if page == nil {
		http.NotFound(res, req)
	}

	data, err := loadData(pageName)
	if err != nil {
		http.Error(res, "page data load error", http.StatusInternalServerError)
		return
	}

	if data != nil {
		t, err := template.New("t").
			Funcs(template.FuncMap{"mod": func(i, j int) int { return i % j }}).
			Funcs(template.FuncMap{"add": func(i, j int) int { return i + j }}).
			Parse(string(page.Body))

		if err != nil {
			http.Error(res, "page template parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(res, data)
		if err != nil {
			http.Error(res, "page template execute error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		res.Write(page.Body)
	}
}
