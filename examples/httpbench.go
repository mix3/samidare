package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mix3/samidare"
)

type myWorker struct {
	URL    string
	client *http.Client
}

func (w *myWorker) Setup() {
	w.client = &http.Client{}
	time.Sleep(1 * time.Second)
}

func (w *myWorker) Teardown() {
	time.Sleep(1 * time.Second)
}

func (w *myWorker) Process() (subscore int) {
	time.Sleep(100 * time.Millisecond)
	resp, err := w.client.Get(w.URL)
	if err == nil {
		defer resp.Body.Close()
		_, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			return 1
		}
	} else {
		log.Printf("err: %v, resp: %#v", err, resp)
	}
	return 0
}

func main() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		})
		http.ListenAndServe(":19301", nil)
	}()
	res := samidare.Run(func() samidare.Worker {
		return &myWorker{URL: "http://localhost:19301"}
	})
	samidare.Print(res)
}
