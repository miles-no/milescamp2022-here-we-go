package main

import (
	"fmt"
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

func PostHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Takk for meldingen ❤️")
}

func GetHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	// Drum roll
	for n := 5; n >= 1; n-- {
		fmt.Fprintf(w, "%d...\n", n)
		flush(w)
		time.Sleep(time.Second)
	}
	fmt.Fprintln(w, "Jeg har ikke blitt implementert enda 🥹")
}

// Flush flushes any buffered data to the client, if supported by the response writer.
func flush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
