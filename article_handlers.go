package prism

import (
	"log"
	"net/http"
)

func ArticlesSearchHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesHandlers(res http.ResponseWriter, req *http.Request) {

	log.Println("Article request incoming")

	method := req.Method

	log.Printf("%v", req)

	path := NewPath(req.URL.Path)

	if !path.hasID() {
		RespondErr(res, req, http.StatusBadRequest, "no id")
		return
	}

	switch method {
	case "GET":
		articlesGetHandler(res, req, path.ID)
		return
	case "POST":
		articlesPostHandler(res, req, path.ID)
		return
	case "PUT":
		articlesPutHandler(res, req, path.ID)
		return
	case "DELETE":
		articlesDeleteHandler(res, req, path.ID)
		return
	}

	RespondErr(res, req, http.StatusBadRequest, "not found")

}

func articlesGetHandler(res http.ResponseWriter, req *http.Request, ID string) {
	RespondErr(res, req, http.StatusInternalServerError, "not implemented")
}

func articlesPostHandler(res http.ResponseWriter, req *http.Request, ID string) {
	RespondErr(res, req, http.StatusInternalServerError, "not implemented")
}

func articlesPutHandler(res http.ResponseWriter, req *http.Request, ID string) {
	RespondErr(res, req, http.StatusInternalServerError, "not implemented")
}

func articlesDeleteHandler(res http.ResponseWriter, req *http.Request, ID string) {
	RespondErr(res, req, http.StatusInternalServerError, "not implemented")
}
