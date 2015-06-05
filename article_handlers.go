package prism

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
)

func ArticlesSearchHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesCRUDHandlers(res http.ResponseWriter, req *http.Request) {

	dbg.Printf("Article request incoming")

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

	dbg.Printf("ID: %d", ID)

	a, _ := LoadArticle(db, ID)

	dbg.Printf("Article: %v", a)

	Respond(res, req, http.StatusOK, *a)
}

func articlesPostHandler(res http.ResponseWriter, req *http.Request) {

	db := GetVar(req, "boltDB").(*bolt.DB)
	jd := json.NewDecoder(req.Body)

	var d ArticleDto
	jd.Decode(&d)

	dbg.Printf("Posted Article Type: %v", d.ArticleType)
	dbg.Printf("Posted Article Content: %v", d.ArticleContent)

	a := NewArticle(db, map[string]interface{}{
		"ArticleType": d.ArticleType,
	})
	d.update(a)
	a.SaveArticle(db)

	Respond(res, req, http.StatusOK, map[string]interface{}{
		"ArticleID": a.GetID(),
	})
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
