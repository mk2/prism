package prism

import (
	"net/http"

	"github.com/boltdb/bolt"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {

	db := GetVar(r, "boltDB").(*bolt.DB)
	stats := db.Stats()

	Respond(w, r, http.StatusOK, stats)
}
