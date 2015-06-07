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

	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func WithBoltDB(db *bolt.DB, fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		SetVar(r, "boltDB", db)
		fn(w, r)
	}
}

func WithLogin(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func WithEnvVars(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var envvar struct {
			GithubAPIKey    string `env:"GITHUB_CLIENTID,required"`
			GithubAPISecret string `env:"GITHUB_CLIENTSECRET,required"`
			SessionSecret   string `env:"SESSION_SECRET,required"`
		}

		if err := envdecode.Decode(&envvar); err != nil {
			log.Fatalln(err)
		}

		SetVar(r, "GithubClientID", envvar.GithubAPIKey)
		SetVar(r, "GithubClientSecret", envvar.GithubAPISecret)
		SetVar(r, "SessionSecret", envvar.SessionSecret)

		fn(w, r)
	}
}

func WithSessionStore(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		sessionSecret := GetVar(r, "SessionSecret").(string)

		sstore := sessions.NewCookieStore(s2b(sessionSecret))

		SetVar(r, "SessionStore", sstore)

		fn(w, r)
	}
}

func WithUser(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		db := GetVar(r, "boltDB").(*bolt.DB)

		sessionStore := GetVar(r, "SessionStore").(*sessions.CookieStore)
		session, _ := sessionStore.Get(r, "prism")

		accessToken, exist := session.Values["gh_access_token"]

		dbg.Printf("AccessToken In Request: %v", accessToken)

		var u *User

		if exist {
			// アクセストークンがあればユーザー情報をロード
			dbg.Printf("Already logined")
			u = LoadGithubUser(db, accessToken.(string))
		} else {
			// 無ければ匿名ユーザー情報を作成する
			dbg.Printf("Anonymous User")
			u = NewUser(db)
		}

		dbg.Printf("CurrentUser: %v", u)

		SetVar(r, "CurrentUser", u)

		fn(w, r)
	}
}

func WithCORS(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if env.Debug {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://prism-client.github.io")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}
