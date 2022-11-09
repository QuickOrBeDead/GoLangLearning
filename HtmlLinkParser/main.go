package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	file, _ := os.Open("test.html")
	defer file.Close()

	rootNode, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	var searchLinks func(*html.Node)
	searchLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Println(attr.Val)
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			searchLinks(child)
		}
	}

	searchLinks(rootNode)
}
