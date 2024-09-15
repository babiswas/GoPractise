package ContextExamples

import (
	"io"
	"log"
	"net/http"
	"time"
)

func WebapiServer(w http.ResponseWriter, r *http.Request) {
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-r.Context().Done():
		log.Println("Error when processing request:", r.Context().Err())
		return
	case <-timer.C:
		log.Println("writing response...")
		_, err := io.WriteString(w, "Hello context")
		if err != nil {
			log.Println("Error when writing response", err)
		}
		return
	}
}

func Webserver() {
	http.HandleFunc("/", WebapiServer)
	log.Println("Starting web server.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
