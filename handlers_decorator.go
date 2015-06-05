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

		sstore := sessions.NewCookieStore(s2b(sessionSecret))

		SetVar(req, "SessionStore", sstore)

		fn(res, req)
	}
}

func WithUser(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		db := GetVar(req, "boltDB").(*bolt.DB)

		sessionStore := GetVar(req, "SessionStore").(*sessions.CookieStore)
		session, _ := sessionStore.Get(req, "prism")

		accessToken, exist := session.Values["gh_access_token"]

		dbg.Printf("AccessToken In Request: %v", accessToken)

		var u *User

		if exist {
			// アクセストークンがあればユーザー情報をロード
			u = LoadGithubUser(db, accessToken.(string))
		} else {
			// 無ければ匿名ユーザー情報を作成する
			u = NewUser(db)
		}

		SetVar(req, "CurrentUser", u)

		fn(res, req)
	}
}

func WithCORS(fn http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		if env.Debug {
			res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		} else {
			res.Header().Set("Access-Control-Allow-Origin", "https://prism-client.github.io")
		}
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(res, req)
	}
}
