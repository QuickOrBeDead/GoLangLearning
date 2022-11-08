package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

var storyData map[string]any
var storyTmpl *template.Template
var storyPageCache map[string][]byte

func main() {
	storyPageCache = make(map[string][]byte)
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

	wr.Write(getStoryPageContents(storyName))
}

func getStoryPageContents(storyName string) []byte {
	content, hasContent := storyPageCache[storyName]
	if !hasContent {
		buf := new(bytes.Buffer)
		storyTmpl.Execute(buf, storyData[storyName])
		content = buf.Bytes()
		buf.Reset()
		storyPageCache[storyName] = content
	}

	return content
}
