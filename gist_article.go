package prism

import "github.com/boltdb/bolt"

const (
	ArticleGistIDBucket = "ArticleGistIDBucket"
)

type GistArticleIface interface {
	GetGistID() string
}

type GistArticle struct {
	ArticleIface
	article *Article
	gistID string
}

func (a *GistArticle) GetGistID() string {

	return a.gistID
}

func (a *GistArticle) initGistArticle(values map[string]interface{}) {

	a.gistID = values["GistID"].(string)
}

func (a *GistArticle) loadGistArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleGistIDBucket))

	bID := s2b(a.GetID())

	a.gistID = b2s(b.Get(bID))

	return nil
}

func (a *GistArticle) saveGistArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleGistIDBucket))

	bID := s2b(a.GetID())

	b.Put(bID, s2b(a.gistID))

	return nil
}
