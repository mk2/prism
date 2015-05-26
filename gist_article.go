package prism

import "github.com/boltdb/bolt"

const (
	ArticleGistIDBucket = "ArticleGistIDBucket"
)

type GistArticle struct {
	ArticleInterface
	GistID string
}

func (l *GistArticle) initGistArticle(values map[string]interface{}) {

}

func (l *GistArticle) saveGistArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleGistIDBucket))

	bID := i2b(l.GetID())

	b.Put(bID, s2b(l.GistID))

	return nil

}
