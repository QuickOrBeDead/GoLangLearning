package main

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Routes struct {
	Map *map[string]*url.URL
}

func (r *Routes) GetRemoteAddress(path string) (*url.URL, bool) {
	v, ok := (*r.Map)[path]

	return v, ok
}

type RouteConfig struct {
	Path          string `yaml:"path"`
	RemoteAddress string `yaml:"remoteAddress"`
}

type ReverseProxyConfig struct {
	Routes []RouteConfig `yaml:"routes"`
	Host   string        `yaml:"host"`
}

type ReverseProxyHandler struct {
	routes *Routes
}

func (proxy *ReverseProxyHandler) CopyResponse(w http.ResponseWriter, r io.Reader, isEventStream bool) {
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
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
	remoteAddress, ok := proxy.routes.GetRemoteAddress(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.Host = remoteAddress.Host
	r.URL.Host = remoteAddress.Host
	r.URL.Scheme = remoteAddress.Scheme
	r.URL.Path = remoteAddress.Path
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
	conf, err := loadConfig()
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(conf.Host, &ReverseProxyHandler{routes: loadRoutes(conf)})
	if err != nil {
		panic(err)
	}
}

func loadRoutes(conf *ReverseProxyConfig) *Routes {
	routes := make(map[string]*url.URL, len(conf.Routes))
	for _, v := range conf.Routes {
		remoteAddress, err := url.Parse(v.RemoteAddress)
		if err != nil {
			panic(err)
		}

		routes[v.Path] = remoteAddress
	}
	return &Routes{Map: &routes}
}

func loadConfig() (*ReverseProxyConfig, error) {
	conf := new(ReverseProxyConfig)
	file, err := os.Open("./config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configBytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
