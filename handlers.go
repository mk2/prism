package prism

import (
	"net/http"

	"github.com/boltdb/bolt"
)

func StatsHandler(res http.ResponseWriter, req *http.Request) {

	db := GetVar(req, "boltDB").(*bolt.DB)
	stats := db.Stats()

	Respond(res, req, http.StatusOK, stats)
}
