package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/QuickOrBeDead/GoLangLearning/lexer"
)

// https://www.w3.org/TR/css-syntax-3/#tokenizing-and-parsing
func main() {
	colors := make(map[lexer.TokenType]string)
	colors[lexer.Ident] = "blue"
	colors[lexer.Function] = "blue"
	colors[lexer.AtKeyword] = "blue"
	colors[lexer.Hash] = "red"
	colors[lexer.String] = "green"
	colors[lexer.BadString] = "green"
	colors[lexer.Url] = "blue"
	colors[lexer.Number] = "gray"
	colors[lexer.Dimension] = "blue"
	colors[lexer.Percentage] = "blue"
	colors[lexer.Whitespace] = "blue"
	colors[lexer.LeftParenthesis] = "blue"
	colors[lexer.RightParenthesis] = "blue"
	colors[lexer.LeftBrace] = "orange"
	colors[lexer.RightBrace] = "orange"
	colors[lexer.Colon] = "yellow"
	colors[lexer.Semicolon] = "yellow"
	colors[lexer.Comma] = "white"
	colors[lexer.Comment] = "blue"
	colors[lexer.At] = "blue"
	colors[lexer.CDO] = "blue"
	colors[lexer.CDC] = "blue"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, err := loadFile("./Pages/Index.html")
		if page == nil || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := template.New("t").
			Funcs(template.FuncMap{"html": func(s string) template.HTML { return template.HTML(s) }}).
			Parse(string(page))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		data := make(map[string]interface{})

		if r.Method == http.MethodPost {
			var sb strings.Builder
			css := r.FormValue("cssText")
			l := lexer.Lexer{Text: []rune(css)}
			for v := l.NextToken(); v.Type != lexer.EOF; v = l.NextToken() {
				if v.Type == lexer.Whitespace {
					sb.WriteString(string(v.Val))
				} else {
					color, ok := colors[v.Type]
					if !ok {
						color = "white"
					}

					sb.WriteString("<span style=\"color:")
					sb.WriteString(color)
					sb.WriteString("\">")
					sb.WriteString(string(v.Val))
					sb.WriteString("</span>")
				}
			}

			data["Result"] = sb.String()
			data["CssText"] = css
		}

		t.Execute(w, data)
	})
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./Pages/css"))))
	http.ListenAndServe(":8080", nil)
}

func loadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err == nil {
		return os.ReadFile(path)
	} else if err == os.ErrNotExist {
		return nil, nil
	} else {
		return nil, err
	}
}
