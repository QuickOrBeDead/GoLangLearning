package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello")
	})

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
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

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
