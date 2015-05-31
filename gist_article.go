package prism

import "github.com/boltdb/bolt"

const (
	ArticleGistIDBucket = "ArticleGistIDBucket"
)

type GistArticle struct {
	ArticleInterface
	article *Article
	GistID  string
}

func (a *GistArticle) initGistArticle(values map[string]interface{}) {

	a.GistID = values["GistID"].(string)
}

func (a *GistArticle) loadGistArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleGistIDBucket))

	bID := s2b(a.GetID())

	a.GistID = b2s(b.Get(bID))

	return nil
}

func (a *GistArticle) saveGistArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleGistIDBucket))

	bID := s2b(a.GetID())

	b.Put(bID, s2b(a.GistID))

	return nil
}
