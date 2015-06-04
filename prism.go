package prism

import (
	"log"

	"github.com/boltdb/bolt"
)

func NewDB() *bolt.DB {

	db, err := bolt.Open("prism.boltdb", 0600, nil)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	CreateArticleBuckets(db)
	CreateUserBuckets(db)

	return db
}
