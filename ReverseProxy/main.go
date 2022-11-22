package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	// TODO: handle text/event-stream content types
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

	for k, v := range remoteResp.Trailer {
		for _, vv := range v {
			remoteResp.Header.Set(k, vv)
		}
	}
}

func main() {
	var port int
	if len(os.Args) < 2 {
		port = 8080
	} else {
		args := os.Args[1:]

		var err error
		port, err = strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
	}

	remoteAddress, err := url.Parse("https://www.google.com")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), &ReverseProxyHandler{remoteAddress: remoteAddress})
	if err != nil {
		panic(err)
	}
}
