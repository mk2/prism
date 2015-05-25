package prism

import "github.com/boltdb/bolt"

const (
	ArticleLinkURLBucket = "ArticleLinkURLBucket"
)

type LinkArticle struct {
	EntityInterface
	LinkURL string
}

func (l *LinkArticle) initLinkArticle(values map[string]interface{}) {

}

func (l *Article) loadLinkArticle(db *bolt.DB) {

}

func (l *LinkArticle) saveLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	bID := i2b(l.GetID())

	b.Put(bID, s2b(l.LinkURL))

	return nil

}
