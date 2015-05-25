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
	Entity
	Terms         []*Meta
	WeightedTerms map[*Meta]float64
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
	articleID, _ := article.NewArticleID(db)

	article.ID = articleID

	if articleType == ArticleTypeLink {
		article.initLinkArticle(values)
	} else if articleType == ArticleTypeGist {
		article.initGistArticle(values)
	} else if articleType == ArticleTypeMarkdown {
		article.initMarkdownArticle(values)
	}

	return

}

func LoadArticle(db *bolt.DB, articleID int) (article *Article) {

	tx, _ := db.Begin(true)
	defer tx.Rollback()

	b := tx.Bucket(s2b(ArticleTypeBucket))
	articleType := b2i(b.Get(i2b(articleID)))

	// value set
	switch articleType {
	case ArticleTypeLink:
	case ArticleTypeGist:
	case ArticleTypeMarkdown:
	}

	return
}

func (a *Article) NewArticleID(db *bolt.DB) (int, error) {

	a.IDBucketName = ArticleIDBucket
	a.IDKey = "articleID"

	return a.newID(db)
}

func (a *Article) Save(db *bolt.DB) {

	if a.ArticleType == ArticleTypeLink {
		db.Update(a.saveLinkArticle)
	} else if a.ArticleType == ArticleTypeMarkdown {
		db.Update(a.saveMarkdownArticle)
	}

}
