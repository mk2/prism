package prism

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
)

func ArticlesSearchHandler(w http.ResponseWriter, r *http.Request) {

	dbg.Printf("Article search request incoming")

	method := r.Method
	db := GetVar(r, "boltDB").(*bolt.DB)
	u := GetVar(r, "CurrentUser").(*User)

	switch method {

	case "GET":
		as, _ := GetAllOwnerArticles(db, u.id)
		Respond(w, r, http.StatusOK, as)
		return

	}

	RespondErr(w, r, http.StatusBadRequest, "unsupported method type"+method)
}

func ArticlesCRUDHandlers(w http.ResponseWriter, r *http.Request) {

	dbg.Printf("Article crud request incoming")

	method := r.Method

	dbg.Printf("%v", r)

	path := NewPath(r.URL.Path)

	switch method {

	case "GET":
		if !path.hasID() {
			RespondErr(w, r, http.StatusBadRequest, "no id")
			return
		}
		articlesGetHandler(w, r, path.ID)
		return

	case "POST":
		articlesPostHandler(w, r)
		return

	case "PUT":
		if !path.hasID() {
			RespondErr(w, r, http.StatusBadRequest, "no id")
			return
		}
		articlesPutHandler(w, r, path.ID)
		return

	case "DELETE":
		if !path.hasID() {
			RespondErr(w, r, http.StatusBadRequest, "no id")
			return
		}
		articlesDeleteHandler(w, r, path.ID)
		return

	}

	RespondErr(w, r, http.StatusBadRequest, "not found")
}

func articlesGetHandler(w http.ResponseWriter, r *http.Request, ID string) {

	db := GetVar(r, "boltDB").(*bolt.DB)

	dbg.Printf("ID: %d", ID)

	a, _ := LoadArticle(db, ID)

	dbg.Printf("Article: %v", a)

	Respond(w, r, http.StatusOK, *a)
}

func articlesPostHandler(w http.ResponseWriter, r *http.Request) {

	db := GetVar(r, "boltDB").(*bolt.DB)
	jd := json.NewDecoder(r.Body)

	var d ArticleDto
	jd.Decode(&d)

	dbg.Printf("Posted Article Type: %v", d.ArticleType)
	dbg.Printf("Posted Article Content: %v", d.ArticleContent)

	a := NewArticle(db, map[string]interface{}{
		"ArticleType": d.ArticleType,
	})
	d.update(a)
	a.SaveArticle(db)

	Respond(w, r, http.StatusOK, map[string]interface{}{
		"ArticleID": a.GetID(),
	})
}

func articlesPutHandler(w http.ResponseWriter, r *http.Request, ID string) {

	db := GetVar(r, "boltDB").(*bolt.DB)
	jd := json.NewDecoder(r.Body)

	var d ArticleDto
	jd.Decode(&d)

	dbg.Printf("Put Article Type: %v", d.ArticleType)
	dbg.Printf("Put Article Content: %v", d.ArticleContent)

	a, _ := LoadArticle(db, ID)
	d.update(a)

	Respond(w, r, http.StatusOK, "ok")
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request, ID string) {
	RespondErr(w, r, http.StatusInternalServerError, "not implemented")
}
