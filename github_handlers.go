package prism

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GithubOAuthHandlers(res http.ResponseWriter, req *http.Request) {

	tokens := strings.Split(req.URL.Path, "/")
	action := tokens[2]

	switch action {

	case "login":
		githubLoginHandler(res, req)

	case "callback":
		githubCallbackHandler(res, req)

	}

	RespondErr(res, req, http.StatusBadRequest, "invalid request")
}

func githubLoginHandler(res http.ResponseWriter, req *http.Request) {

	var clientID, status string = GetVar(req, "GithubClientID").(string), "dummy"

	loginURL := fmt.Sprintf("http://github.com/login/oauth/authorize?client_id=%s&status=%s", clientID, status)

	res.Header().Set("Location", loginURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func githubCallbackHandler(res http.ResponseWriter, req *http.Request) {

	var code,
		status,
		clientID,
		clientSecret,
		atURL string = req.Form.Get("code"), req.Form.Get("status"), GetVar(req, "GithubClientID").(string), GetVar(req, "GithubClientSecret").(string), "https://github.com/login/oauth/access_token"

	dbg.Printf("Status: %s", status)

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", code)

	atres, _ := http.PostForm(atURL, q)

	bytes, _ := ioutil.ReadAll(atres.Body)

	authRes, _ := url.ParseQuery(string(bytes))

	accessToken := authRes.Get("access_token")

	dbg.Printf("AccessToken: %s", accessToken)
}
