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
	defer db.Close()

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

	mux.HandleFunc("/stats",
		prism.WithVars(prism.WithBoltDB(db, prism.StatsHandler)))

	mux.HandleFunc("/articles/",
		withBaseDecorators(db, prism.ArticlesCRUDHandlers))

	mux.HandleFunc("/articles",
		withBaseDecorators(db, prism.ArticlesSearchHandler))

	mux.HandleFunc("/auth/github/",
		withBaseDecorators(db, prism.GithubOAuthHandlers))

	mux.HandleFunc("/gists",
		withBaseDecorators(db, prism.GistsHandlers))

}

func withBaseDecorators(db *bolt.DB, fn http.HandlerFunc) http.HandlerFunc {

	return prism.WithCORS(prism.WithVars(prism.WithEnvVars(prism.WithSessionStore(prism.WithBoltDB(db, prism.WithUser(fn))))))

}
