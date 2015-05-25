package prism

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"github.com/joeshaw/envdecode"
)

func WithVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		OpenVars(req)
		defer CloseVars(req)
		fn(res, req)
	}
}

func WithBoltDB(db *bolt.DB, fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		SetVar(req, "boltDB", db)
		fn(res, req)
	}
}

func WithLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fn(res, req)
	}
}

func WithEnvValues(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var envvar struct {
			GithubAPIKey    string `env:"GITHUB_APIKEY,required"`
			GithubAPISecret string `env:"GITHUB_APISECRET,required"`
			SessionSecret   string `env:"SESSION_SECRET,required"`
		}

		if err := envdecode.Decode(&envvar); err != nil {
			log.Fatalln(err)
		}

		SetVar(req, "GithubAPIKey", envvar.GithubAPIKey)
		SetVar(req, "GithubAPISecret", envvar.GithubAPISecret)
		SetVar(req, "SessionSecret", envvar.SessionSecret)

		fn(res, req)
	}
}

func WithCookieStore(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		sessionSecret := GetVar(req, "SessionSecret").(string)

		cstore := sessions.NewCookieStore(s2b(sessionSecret))

		SetVar(req, "CookieStore", cstore)

		fn(res, req)
	}
}
