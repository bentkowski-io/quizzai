package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/bentkowski-io/cfglib"
	"github.com/bentkowski-io/quizzai/gemini"
	"github.com/bentkowski-io/quizzai/openai"
)

var counter atomic.Int64

func main() {
	p, err := cfglib.NewFileEnvProvider(".env")
	if err != nil {
		log.Fatal(err)
	}
	cfg := cfglib.NewWithEnvProvider(p)
	h := func(w http.ResponseWriter, r *http.Request) {
		counter.Add(1)
		w.Write([]byte(openai.Example(cfg)))
	}
	http.HandleFunc("/openai", h)

	h = func(w http.ResponseWriter, r *http.Request) {
		counter.Add(1)
		if err := gemini.GenerateContentFromText(w, cfg); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	http.HandleFunc("/gemini", h)

	addr := cfg.ReadString("server.address", ":8080")
	log.Default().Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}
