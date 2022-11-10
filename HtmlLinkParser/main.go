package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func main() {
	file, _ := os.Open("test.html")
	defer file.Close()

	rootNode, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	links := make([]Link, 0, 10)
	var searchLinks func(*html.Node)
	searchLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					var text string
					if n.FirstChild != nil {
						text = getText(n.FirstChild)
					}
					links = append(links, Link{Href: attr.Val, Text: text})
					break
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			searchLinks(child)
		}
	}

	searchLinks(rootNode)

	fmt.Println(links)
}

func getText(n *html.Node) string {
	var sb strings.Builder
	var t func(n *html.Node)
	t = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.Join(strings.Fields(n.Data), " ")
			sb.WriteString(text)
			sb.WriteRune(' ')
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			t(child)
		}

		if n.NextSibling != nil {
			t(n.NextSibling)
		}
	}

	t(n)

	return strings.Trim(sb.String(), " ")
}
