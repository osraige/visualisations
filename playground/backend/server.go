package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/osraige/visualisations/playground/backend/goexports"
)

func envDefault(key, or string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return or
}

func main() {

	confAssets := envDefault("TA_ASSETS", "../frontend/public")
	confListenAddr := envDefault("TA_LISTEN_ADDR", ":8484")
	http.Handle("/", http.FileServer(http.Dir(confAssets)))
	http.HandleFunc("/svg", func(w http.ResponseWriter, r *http.Request) {
		param, ok := r.URL.Query()["data"]
		toEval := strings.Join(param, "")
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		i.Use(goexports.Symbols)
		_, err := i.Eval(toEval)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(toEval)
			log.Println(err)
			return
		}
		v, err := i.Eval("genSvg")
		if err != nil {
			panic(err)
		}
		genSvg := v.Interface().(func(io.Writer))
		genSvg(w)
	})
	log.Printf("starting on %q\n", confListenAddr)
	if err := http.ListenAndServe(confListenAddr, nil); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
