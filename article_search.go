package prism

import "github.com/boltdb/bolt"

func GetAllOwnerArticles(db *bolt.DB, ownerID string) ([]*Article, error) {

	return LoadArticlesByOwnerID(db, ownerID)
}
