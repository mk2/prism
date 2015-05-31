package prism

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"github.com/joeshaw/envdecode"
	"github.com/mk2/prism/env"
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

func WithEnvVars(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		var envvar struct {
			GithubAPIKey    string `env:"GITHUB_CLIENTID,required"`
			GithubAPISecret string `env:"GITHUB_CLIENTSECRET,required"`
			SessionSecret   string `env:"SESSION_SECRET,required"`
		}

		if err := envdecode.Decode(&envvar); err != nil {
			log.Fatalln(err)
		}

		SetVar(req, "GithubClientID", envvar.GithubAPIKey)
		SetVar(req, "GithubClientSecret", envvar.GithubAPISecret)
		SetVar(req, "SessionSecret", envvar.SessionSecret)

		fn(res, req)
	}
}

func WithSessionStore(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		sessionSecret := GetVar(req, "SessionSecret").(string)

		cstore := sessions.NewCookieStore(s2b(sessionSecret))

		SetVar(req, "SessionStore", cstore)

		fn(res, req)
	}
}

func WithSessionID(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		sessionStore := GetVar(req, "SessionStore").(*sessions.CookieStore)
		session, _ := sessionStore.Get(req, "prisim")

		session.Values["id"] = "fdafasf"
		session.Save(req, res)

		fn(res, req)
	}
}

func WithCORS(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		if env.Debug {
			res.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			res.Header().Set("Access-Control-Allow-Origin", "https://prism-client.github.io/")
		}
		res.Header().Set("Access-Control-Expose-Headers", "Location")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(res, req)
	}
}
