package prism

import "net/http"

func GithubOAuthHandlers(res http.ResponseWriter, req *http.Request) {

	path := NewPath(req.URL.Path)

	if !path.hasID() {
		githubOAuthInfoHandler(res, req)
		return
	}

	switch path.ID {
	case "redirect":
	}

	RespondErr(res, req, http.StatusBadRequest, "invalid request")

}

func githubOAuthInfoHandler(res http.ResponseWriter, req *http.Request) {

	clientID := GetVar(req, "GithubClientID").(string)

	Respond(res, req, http.StatusOK, map[string]string{
		"oauth_entry_url":     "https://github.com/login/oauth/authorize",
		"client_id":           clientID,
		"github_redirect_uri": "http://localhost:13333/ghoauth/redirect",
		"state":               "dummy",
	})

}
