package prism

import (
	"log"
	"testing"

	"github.com/boltdb/bolt"
)

const (
	TestDatabaseName = "test.boltdb"
)

func NewTestDB() *bolt.DB {

	db, err := bolt.Open(TestDatabaseName, 0600, nil)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	DeleteBuckets(db)
	CreateBuckets(db)

	return db
}

func TestNewArticleID(t *testing.T) {

	var article, nextArticleID, currentArticleID = &Article{}, 0, 0

	db := NewTestDB()

	for i := 0; i < 10; i++ {
		nextArticleID, _ = article.newArticleID(db)

		log.Printf("current:%d new:%d", currentArticleID, nextArticleID)

		if (currentArticleID + 1) != nextArticleID {
			t.Errorf("Unexpected nextArticleID: %d", nextArticleID)
		}

		currentArticleID = nextArticleID
	}
}
