package prism

import "github.com/boltdb/bolt"

const (
	ArticleLinkURLBucket = "ArticleLinkURLBucket"
)

type LinkArticle struct {
	article *Article
	LinkURL string
}

func (a *LinkArticle) initLinkArticle(values map[string]interface{}) {

	a.LinkURL, _ = values["LinkURL"].(string)

}

func (a *LinkArticle) loadLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	a.LinkURL = b2s(b.Get(s2b(a.article.GetID())))

	return nil

}

func (a *LinkArticle) saveLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	ID := a.article.GetID()
	bID := s2b(ID)

	b.Put(bID, s2b(a.LinkURL))

	return nil

}
