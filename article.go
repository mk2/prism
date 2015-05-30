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

type ArticleInterface interface {
	EntityInterface
}

type Article struct {
	ArticleInterface
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

/*
NewArticle 新規アーティクルを作成する
*/
func NewArticle(db *bolt.DB, values map[string]interface{}) (article *Article) {

	articleType, _ := values["articleType"].(int)
	articleID, _ := article.newArticleID(db)

	article.ID = articleID

	if articleType == ArticleTypeLink {
		article.initLinkArticle(values)
	} else if articleType == ArticleTypeGist {
		article.initGistArticle(values)
	} else if articleType == ArticleTypeMarkdown {
		article.initMarkdownArticle(values)
	}

	created := article.Created()
	updated := article.Updated()

	dbg.Printf("Created: %v", created)
	dbg.Printf("Updated: %v", updated)

	return

}

/*
LoadArticle アーティクルを読み込む
*/
func LoadArticle(db *bolt.DB, articleID int) (*Article, error) {

	article := &Article{}

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(s2b(ArticleTypeBucket))
		articleType := b2i(b.Get(i2b(articleID)))

		// value set
		switch articleType {
		case ArticleTypeLink:
		case ArticleTypeGist:
		case ArticleTypeMarkdown:
		}

		return nil
	})

	return article, err
}

/*
SaveArticle アーティクルを保存する
*/
func (a *Article) SaveArticle(db *bolt.DB) error {

	return db.Batch(func(tx *bolt.Tx) error {

		if a.ArticleType == ArticleTypeLink {
			db.Update(a.saveLinkArticle)
		} else if a.ArticleType == ArticleTypeMarkdown {
			db.Update(a.saveMarkdownArticle)
		}

		return nil
	})

}

func (a *Article) newArticleID(db *bolt.DB) (int, error) {

	a.IDBucketName = ArticleIDBucket
	a.IDKey = "articleID"

	return a.newID(db)
}
