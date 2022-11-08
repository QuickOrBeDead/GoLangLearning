package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

var storyData map[string]any
var storyTmpl *template.Template

func main() {
	storyTmpl = template.Must(template.ParseFiles("Templates/Story.html"))
	storyJsonContent, err := os.ReadFile("Choose Your Own Adventure.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(storyJsonContent, &storyData)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", storyHandler)
	http.ListenAndServe(":8080", nil)
}

func storyHandler(wr http.ResponseWriter, req *http.Request) {
	storyName := req.URL.Path[1:]
	if storyName == "" {
		storyName = "intro"
	}

	storyTmpl.Execute(wr, storyData[storyName])
}
