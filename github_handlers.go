package prism

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func GithubOAuthHandlers(w http.ResponseWriter, r *http.Request) {

	tokens := strings.Split(r.URL.Path, "/")
	action := tokens[3]

	dbg.Printf("action: %s", action)

	switch action {

	case "login":
		githubLoginHandler(w, r)
		return

	case "callback":
		githubCallbackHandler(w, r)
		return

	}

	RespondErr(w, r, http.StatusBadRequest, "invalid request")
}

func GistsHandlers(w http.ResponseWriter, r *http.Request) {

	sessionStore := GetVar(r, "SessionStore").(*sessions.CookieStore)
	session, _ := sessionStore.Get(r, "prism")

	accessToken := session.Values["gh_access_token"]

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	cli := github.NewClient(tc)

	gists, _, _ := cli.Gists.ListAll(nil)

	Respond(w, r, http.StatusOK, gists)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {

	var clientID,
		state,
		reqURL string = GetVar(r, "GithubClientID").(string), "dummy", "https://github.com/login/oauth/authorize"

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("state", state)
	q.Set("scope", "gist")

	loginURL := fmt.Sprintf("%s?%s", reqURL, q.Encode())

	dbg.Printf("login url: %s", loginURL)

	w.Header().Set("Location", loginURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var code,
		state,
		clientID,
		clientSecret,
		atURL string = r.FormValue("code"), r.FormValue("state"), GetVar(r, "GithubClientID").(string), GetVar(r, "GithubClientSecret").(string), "https://github.com/login/oauth/access_token"

	var u *User = GetVar(r, "CurrentUser").(*User)

	var db *bolt.DB = GetVar(r, "boltDB").(*bolt.DB)

	dbg.Printf("code: %s", code)
	dbg.Printf("state: %s", state)

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", code)

	atres, _ := http.PostForm(atURL, q)

	bytes, _ := ioutil.ReadAll(atres.Body)

	authRes, _ := url.ParseQuery(string(bytes))

	accessToken := authRes.Get("access_token")

	dbg.Printf("AccessToken: %s", accessToken)

	sessionStore := GetVar(r, "SessionStore").(*sessions.CookieStore)
	session, _ := sessionStore.Get(r, "prism")

	session.Values["gh_access_token"] = accessToken
	session.Save(r, w)

	u.AccessToken = accessToken

	u.SaveUser(db)

	Respond(w, r, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"user": map[string]string{
			"name": u.name,
		},
	})
}
