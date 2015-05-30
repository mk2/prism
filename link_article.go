package prism

import "github.com/boltdb/bolt"

const (
	ArticleLinkURLBucket = "ArticleLinkURLBucket"
)

type LinkArticle struct {
	ArticleInterface
	LinkURL string
}

func (l *LinkArticle) initLinkArticle(values map[string]interface{}) {

	l.LinkURL, _ = values["LinkURL"].(string)

}

func (l *LinkArticle) loadLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	l.LinkURL = b2s(b.Get(i2b(l.GetID())))

	return nil

}

func (l *LinkArticle) saveLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	bID := i2b(l.GetID())

	b.Put(bID, s2b(l.LinkURL))

	return nil

}
