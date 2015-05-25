package prism

import (
	"log"
	"net/http"
)

func ArticlesSearchHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesHandlers(res http.ResponseWriter, req *http.Request) {

	method := req.Method

	log.Printf("%v", req)

	switch method {
	case "GET":
		ArticlesGetHandler(res, req)
	case "POST":
		ArticlesPostHandler(res, req)
	case "PUT":
		ArticlesPutHandler(res, req)
	}

}

func ArticlesGetHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesPostHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesPutHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesDeleteHandler(res http.ResponseWriter, req *http.Request) {

}
