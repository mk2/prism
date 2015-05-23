package main

import (
	"net/http"
	"time"

	"github.com/mk2/prism"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	mux := http.NewServeMux()
	setupHandlers(mux)

	srv := &graceful.Server{

		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    ":13333",
			Handler: mux,
		},
	}

	srv.ListenAndServe()

}

func setupHandlers(mux *http.ServeMux) {
	db := prism.NewDB()

	mux.HandleFunc("/articles/",
		prism.WithVars(prism.WithBoltDB(db, prism.ArticlesHandler)))

}
