package prism

import "github.com/boltdb/bolt"

const (
	_ = iota
	ArticleTypeLink
	ArticleTypeGist
	ArticleTypeMarkdown
)

const (
	ArticleIDBucket   = "ArticleIdBucket"
	ArticleTypeBucket = "ArticleTypeBucket"
)

type Article struct {
	ID            int
	Terms         []*Term
	WeightedTerms map[*Term]float64
	ArticleType   int
	LinkArticle
	GistArticle
	MarkdownArticle
}

func CreateBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		ArticleIDBucket,
		ArticleTypeBucket,
		ArticleLinkURLBucket,
		ArticleMarkdownTextBucket,
	}

	for _, requiredBucket := range requiredBuckets {
		_, err := tx.CreateBucketIfNotExists(s2b(requiredBucket))

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func DeleteBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		ArticleIDBucket,
		ArticleTypeBucket,
		ArticleLinkURLBucket,
		ArticleMarkdownTextBucket,
	}

	for _, requiredBucket := range requiredBuckets {
		err := tx.DeleteBucket(s2b(requiredBucket))

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func NewArticle(db *bolt.DB, values map[string]interface{}) (article *Article) {

	articleType, _ := values["articleType"].(int)
	articleID, _ := newArticleID(db)

	article.ID = articleID

	if articleType == ArticleTypeLink {
		article.initLinkArticle(values)
	} else if articleType == ArticleTypeMarkdown {
		article.initMarkdownArticle(values)
	}

	return

}

func (a *Article) Save(db *bolt.DB) {

	if a.ArticleType == ArticleTypeLink {
		db.Update(a.saveLinkArticle)
	} else if a.ArticleType == ArticleTypeMarkdown {
		db.Update(a.saveMarkdownArticle)
	}

}

func newArticleID(db *bolt.DB) (int, error) {

	tx, _ := db.Begin(true)
	defer tx.Rollback()

	b := tx.Bucket(s2b(ArticleIDBucket))

	lastArticleID := b2i(b.Get(s2b("lastArticleID")))
	nextArticleID := lastArticleID + 1

	b.Put(s2b("lastArticleID"), i2b(nextArticleID))

	result := tx.Commit()

	return nextArticleID, result
}
