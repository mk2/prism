package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mk2/prism"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	log.Println("Start setup Prism")

	db := prism.NewDB()

	mux := http.NewServeMux()

	setupHandlers(db, mux)

	srv := &graceful.Server{

		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    ":13333",
			Handler: mux,
		},
	}

	log.Println("End setup Prism and start serving")

	srv.ListenAndServe()

}

func setupHandlers(db *bolt.DB, mux *http.ServeMux) {

	mux.HandleFunc("/articles/",
		withBaseDecorators(db, prism.ArticlesHandlers))

	mux.HandleFunc("/ghoauth/",
		withBaseDecorators(db, prism.GithubOAuthHandlers))

}

func withBaseDecorators(db *bolt.DB, fn http.HandlerFunc) http.HandlerFunc {

	return prism.WithCORS(prism.WithVars(prism.WithEnvVars(prism.WithSessionStore(prism.WithBoltDB(db, fn)))))

}
