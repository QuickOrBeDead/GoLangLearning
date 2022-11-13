package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/QuickOrBeDead/GoLangLearning/datastructures"

	"golang.org/x/net/html"
)

func main() {
	rootUrl, err := url.Parse("https://www.calhoun.io")
	if err != nil {
		fmt.Println("invalid url")
		return
	}

	rootNode, err := parseHtml(rootUrl.String())
	if err != nil {
		fmt.Println("error parsing html: ", err.Error())
		return
	}

	links := make(datastructures.LinkSet, 128)
	links.Add("/")

	var searchLinks func(*html.Node)
	searchLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					currentUrl, err := url.Parse(href)

					if err == nil {
						if currentUrl.Host == "" || (currentUrl.Host == rootUrl.Host && currentUrl.Scheme == rootUrl.Scheme) {
							if currentUrl.Path != "" {
								if !links.Contains(currentUrl.Path) {
									links.Add(currentUrl.Path)
									if len(links) >= 50 {
										return
									}

									time.Sleep(1 * time.Second)
									newUrl := rootUrl.JoinPath(currentUrl.Path).String()
									node, err := parseHtml(newUrl)
									if err == nil {
										searchLinks(node)
									} else {
										fmt.Println("error parsing html: ", err.Error())
									}
								}
							}
						} else {
							break
						}
					}

					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchLinks(c)
		}
	}

	searchLinks(rootNode)

	fmt.Println("links count: ", len(links))
	for k := range links {
		fmt.Println(k)
	}
}

func parseHtml(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("page get error. status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rootNode, err := html.Parse(bytes.NewReader(respBody))
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}
