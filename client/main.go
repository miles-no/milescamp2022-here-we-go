package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var DefaultClient = &http.Client{Timeout: time.Minute}

func init() { log.SetFlags(log.Ltime) }

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	cmd, url := os.Args[1], os.Args[2]
	switch strings.ToUpper(cmd) {
	case "CONSUME":
		if err := consume(url, os.Stdout); err != nil {
			log.Fatal(err)
		}
	case "PRODUCE":
		if err := produce(url, os.Stdout, os.Stdin); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func usage() {
	log.Fatal("Usage: %s (CONSUME or PRODUCE) (url)", os.Args[0])
}

func consume(url string, w io.Writer) error              { return doRequest("GET", url, w, nil) }
func produce(url string, w io.Writer, r io.Reader) error { return doRequest("POST", url, w, r) }

func doRequest(method, url string, w io.Writer, r io.Reader) error {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return err
	}
	log.Printf("%s %s", req.Method, url)
	resp, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Printf("%d - %s", resp.StatusCode, resp.Status)
	_, err = io.Copy(w, resp.Body)
	return err
}
