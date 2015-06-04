package prism

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/sessions"
)

func GithubOAuthHandlers(res http.ResponseWriter, req *http.Request) {

	tokens := strings.Split(req.URL.Path, "/")
	action := tokens[3]

	dbg.Printf("action: %s", action)

	switch action {

	case "login":
		githubLoginHandler(res, req)
		return

	case "callback":
		githubCallbackHandler(res, req)
		return

	}

	RespondErr(res, req, http.StatusBadRequest, "invalid request")
}

func githubLoginHandler(res http.ResponseWriter, req *http.Request) {

	var clientID,
		state,
		reqURL string = GetVar(req, "GithubClientID").(string), "dummy", "https://github.com/login/oauth/authorize"

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("state", state)
	q.Set("scope", "gist")

	loginURL := fmt.Sprintf("%s?%s", reqURL, q.Encode())

	dbg.Printf("login url: %s", loginURL)

	res.Header().Set("Location", loginURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func githubCallbackHandler(res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	var code,
		state,
		clientID,
		clientSecret,
		atURL string = req.FormValue("code"), req.FormValue("state"), GetVar(req, "GithubClientID").(string), GetVar(req, "GithubClientSecret").(string), "https://github.com/login/oauth/access_token"

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

	sessionStore := GetVar(req, "SessionStore").(*sessions.FilesystemStore)
	session, _ := sessionStore.Get(req, "prism")

	session.Values["gh_access_token"] = accessToken
	session.Save(req, res)

	Respond(res, req, http.StatusOK, "ok")
}
