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
)

type ReverseProxyHandler struct {
	remoteAddress *url.URL
}

func (proxy *ReverseProxyHandler) CopyResponse(w http.ResponseWriter, r io.Reader, isEventStream bool) {
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err == io.EOF {
			break
		}

		if n == 0 {
			break
		}

		w.Write(buf[:n])
		if isEventStream {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
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

	remoteRespContentType := remoteResp.Header.Get("Content-Type")
	isEventStream := false
	if remoteRespMediaType, _, _ := mime.ParseMediaType(remoteRespContentType); remoteRespMediaType == "text/event-stream" {
		isEventStream = true
	}

	proxy.CopyResponse(w, remoteResp.Body, isEventStream)

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
	flag.StringVar(&targetUrl, "url", "http://localhost:8888", "url")
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
