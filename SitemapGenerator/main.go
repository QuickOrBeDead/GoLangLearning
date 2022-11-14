package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/QuickOrBeDead/GoLangLearning/datastructures"

	"golang.org/x/net/html"
)

type urlxml struct {
	Location string   `xml:"loc"`
	XMLName  struct{} `xml:"url"`
}

func main() {
	var (
		targetUrl string
		printType string
	)
	flag.StringVar(&targetUrl, "url", "https://www.calhoun.io", "the target url to analyze")
	flag.StringVar(&printType, "type", "Sitemap", "the print type: Sitemap | Tree")
	flag.Parse()

	rootUrl, err := url.Parse(targetUrl)
	if err != nil {
		fmt.Println("invalid url")
		return
	}

	rootNode, err := parseHtml(rootUrl.String())
	if err != nil {
		fmt.Println("error parsing html: ", err.Error())
		return
	}

	switch printType {
	case "Tree":
		sitemap := &datastructures.SitemapNode{Url: "/", Children: make([]*datastructures.SitemapNode, 0), Parent: nil, Level: 0}
		searchPageLinks(rootNode, sitemap, func(u *url.URL, data any) (bool, any, *url.URL) {
			if u.Path != "" && (u.Host == "" || (u.Host == rootUrl.Host && u.Scheme == rootUrl.Scheme)) {
				c := data.(*datastructures.SitemapNode).AddChild(u.Path)
				if c == nil || c.Level > 1 {
					return false, nil, nil
				}

				return true, c, rootUrl.JoinPath(c.Url)
			}

			return false, nil, nil
		})
		printSitemap(sitemap)
	case "Sitemap":
		set := datastructures.LinkSet{}
		searchPageLinks(rootNode, 0, func(u *url.URL, data any) (bool, any, *url.URL) {
			if u.Path != "" && (u.Host == "" || (u.Host == rootUrl.Host && u.Scheme == rootUrl.Scheme)) {
				level := data.(int)
				if level > 1 || set.Contains(u.Path) {
					return false, nil, nil
				}

				set.Add(u.Path)

				return true, level + 1, rootUrl.JoinPath(u.Path)
			}

			return false, nil, nil
		})

		header := xml.ProcInst{
			Target: "xml", Inst: []byte(`version="1.0" encoding="UTF-8"`),
		}
		startElement := xml.StartElement{
			Name: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
			Attr: []xml.Attr{
				{Name: xml.Name{Space: "", Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
				{Name: xml.Name{Space: "", Local: "xsi:schemaLocation"}, Value: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd/sitemap.xsd"},
			},
		}
		buf := new(bytes.Buffer)
		xmlEnc := xml.NewEncoder(buf)
		xmlEnc.EncodeToken(header)
		xmlEnc.EncodeToken(startElement)

		for k := range set {
			xmlEnc.Encode(&urlxml{Location: k})
		}

		xmlEnc.EncodeToken(startElement.End())
		xmlEnc.Flush()

		fmt.Println(buf.String())
	}
}

func searchPageLinks(n *html.Node, data any, f func(*url.URL, any) (bool, any, *url.URL)) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				href := attr.Val
				currentUrl, err := url.Parse(href)

				if err == nil {
					if currentUrl.Path != "" {
						c, d, u := f(currentUrl, data)
						if !c {
							break
						}

						time.Sleep(1 * time.Second)
						node, err := parseHtml(u.String())
						if err == nil {
							searchPageLinks(node, d, f)
						} else {
							fmt.Println("error parsing html: ", err.Error())
						}
					}

					break
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		searchPageLinks(c, data, f)
	}
}

func printSitemap(s *datastructures.SitemapNode) {
	fmt.Printf("%s%s\n", strings.Repeat("   ", s.Level), s.Url)

	for _, c := range s.Children {
		printSitemap(c)
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
