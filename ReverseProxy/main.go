package main

import (
	"flag"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ReverseProxyHandler struct {
	remoteAddress *url.URL
}

func (proxy *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Host = proxy.remoteAddress.Host
	r.URL.Host = proxy.remoteAddress.Host
	r.URL.Scheme = proxy.remoteAddress.Scheme
	r.RequestURI = ""

	remoteResp, err := http.DefaultClient.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	defer remoteResp.Body.Close()

	for k, v := range remoteResp.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		w.Header().Set("X-Forwarded-For", ip)
	}

	remoteRespContentType := remoteResp.Header.Get("Content-Type")

	writeBodyDone := make(chan bool)

	if remoteRespMediaType, _, _ := mime.ParseMediaType(remoteRespContentType); remoteRespMediaType == "text/event-stream" {
		go func() {
			for {
				select {
				case <-time.Tick(10 * time.Millisecond):
					w.(http.Flusher).Flush()
				case <-writeBodyDone:
					return
				}
			}
		}()
	}
	// TODO: handle http2

	trailerKeys := make([]string, len(remoteResp.Trailer))
	i := 0
	for k := range remoteResp.Trailer {
		trailerKeys[i] = k
		i++
	}

	if len(trailerKeys) > 0 {
		w.Header().Set("Trailer", strings.Join(trailerKeys, ","))
	}

	w.WriteHeader(remoteResp.StatusCode)
	io.Copy(w, remoteResp.Body)

	defer close(writeBodyDone)

	for k, v := range remoteResp.Trailer {
		for _, vv := range v {
			remoteResp.Header.Set(k, vv)
		}
	}
}

func main() {
	var targetUrl string
	var port int

	flag.IntVar(&port, "port", 8080, "port")
	flag.StringVar(&targetUrl, "url", "https://www.google.com", "url")
	flag.Parse()

	remoteAddress, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), &ReverseProxyHandler{remoteAddress: remoteAddress})
	if err != nil {
		panic(err)
	}
}
