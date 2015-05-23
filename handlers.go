package prism

import (
	"log"
	"net/http"
	"sync"

	"github.com/boltdb/bolt"
)

var vars map[*http.Request]map[string]interface{}
var varsLock sync.RWMutex

// GetVar gets the value of the key for the specified http.Request.
func GetVar(req *http.Request, key string) interface{} {
	varsLock.RLock()
	value := vars[req][key]
	varsLock.RUnlock()
	return value
}

// SetVar sets the key to the value for the specified http.Request.
func SetVar(req *http.Request, key string, value interface{}) {
	varsLock.Lock()
	vars[req][key] = value
	varsLock.Unlock()
}

// OpenVars opens the vars for the specified http.Request.
// Must be called before GetVar or SetVar is called for each
// request.
func OpenVars(req *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[req] = map[string]interface{}{}
	varsLock.Unlock()
}

// CloseVars closes the vars for the specified
// http.Request.
// Must be called when all var activity is completed to
// clean up any used memory.
func CloseVars(res *http.Request) {
	varsLock.Lock()
	delete(vars, res)
	varsLock.Unlock()
}

func WithVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		OpenVars(req)
		defer CloseVars(req)
		fn(res, req)
	}
}

func WithBoltDB(db *bolt.DB, fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		SetVar(req, "boltDB", db)
		fn(res, req)
	}
}

func RootHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesSearchHandler(res http.ResponseWriter, req *http.Request) {

}

func ArticlesHandler(res http.ResponseWriter, req *http.Request) {

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
