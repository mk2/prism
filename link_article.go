package prism

import "github.com/boltdb/bolt"

const (
	ArticleLinkURLBucket = "ArticleLinkURLBucket"
)

type LinkArticle struct {
	LinkURL string
}

func (l *LinkArticle) initLinkArticle(values map[string]interface{}) {

}

func (l *Article) loadLinkArticle(db *bolt.DB) {

}

func (l *Article) saveLinkArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	bID := i2b(l.ID)

	b.Put(bID, s2b(l.LinkURL))

	return nil

}
