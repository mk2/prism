package prism

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

func ArticlesSearchHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesHandlers(res http.ResponseWriter, req *http.Request) {

	log.Println("Article request incoming")

	method := req.Method

	dbg.Printf("%v", req)

	path := NewPath(req.URL.Path)

	switch method {
	case "GET":
		if !path.hasID() {
			RespondErr(res, req, http.StatusBadRequest, "no id")
			return
		}
		articlesGetHandler(res, req, path.ID)
		return
	case "POST":
		articlesPostHandler(res, req)
		return
	case "PUT":
		if !path.hasID() {
			RespondErr(res, req, http.StatusBadRequest, "no id")
			return
		}
		articlesPutHandler(res, req, path.ID)
		return
	case "DELETE":
		if !path.hasID() {
			RespondErr(res, req, http.StatusBadRequest, "no id")
			return
		}
		articlesDeleteHandler(res, req, path.ID)
		return
	}

	RespondErr(res, req, http.StatusBadRequest, "not found")

}

func articlesGetHandler(res http.ResponseWriter, req *http.Request, ID string) {
	db := GetVar(req, "boltDB").(*bolt.DB)

	a, _ := LoadArticle(db, ID)

	Respond(res, req, http.StatusOK, a)
}

func articlesPostHandler(res http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)

	var a ArticleDto
	d.Decode(&a)

	dbg.Printf("Posted Article Type: %v", a.ArticleType)
	dbg.Printf("Posted Article Content: %v", a.ArticleContent)

	Respond(res, req, http.StatusOK, "ok")
}

func articlesPutHandler(res http.ResponseWriter, req *http.Request, ID string) {
	db := GetVar(req, "boltDB").(*bolt.DB)
	jd := json.NewDecoder(req.Body)

	var d ArticleDto
	jd.Decode(&d)

	dbg.Printf("Put Article Type: %v", d.ArticleType)
	dbg.Printf("Put Article Content: %v", d.ArticleContent)

	a, _ := LoadArticle(db, ID)
	d.update(a)

	Respond(res, req, http.StatusOK, "ok")
}

func articlesDeleteHandler(res http.ResponseWriter, req *http.Request, ID string) {
	RespondErr(res, req, http.StatusInternalServerError, "not implemented")
}
