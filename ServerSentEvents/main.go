package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("port argument is required")
		return
	}

	port, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for i := 0; i < 10; i++ {
			fmt.Fprintf(w, "data: %v\n\n", i)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			time.Sleep(2 * time.Second)
		}
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
