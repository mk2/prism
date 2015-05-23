package prism

import (
	"net/http"

	"github.com/boltdb/bolt"
)

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

func WithLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fn(res, req)
	}
}
