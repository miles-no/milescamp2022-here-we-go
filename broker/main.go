package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/messages", PostHandler).Methods("POST")
	router.HandleFunc("/messages/next", GetHandler).Methods("GET")

	srv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		Handler:      router,
	}
	log.Printf("Listening on %s...", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type Message string

var messages = make(chan Message)

func PostHandler(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	select {
	case messages <- Message(data):
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Meldingen er levert ❤️")
		log.Println("Message received and forwarded.")
	case <-time.After(20 * time.Second):
		w.WriteHeader(http.StatusRequestTimeout)
		fmt.Fprintln(w, "Tidsavbrudd - Ingen mottakere å sende til 😞")
		log.Println("Message received, but timed out while waiting for a recipient.")
	case <-req.Context().Done():
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Message received, but the producer cancelled while waiting for a recipient.")
	}
}

func GetHandler(w http.ResponseWriter, req *http.Request) {
	for {
		log.Println("Client is waiting for the next message...")
		select {
		case msg := <-messages:
			fmt.Fprintln(w, msg)
			flush(w)
		case <-req.Context().Done():
			log.Println("Client exited")
			return
		}
	}
}

// Flush flushes any buffered data to the client, if supported by the response writer.
func flush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
