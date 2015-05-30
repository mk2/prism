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
func NewArticle(db *bolt.DB, values map[string]interface{}) (a *Article) {

	db.Batch(func(tx *bolt.Tx) error {

		articleType, _ := values["articleType"].(int)
		articleID, _ := a.newArticleID(db)

		a.ID = articleID
		created := a.Created()
		updated := a.Updated()

		if articleType == ArticleTypeLink {
			a.initLinkArticle(values)
			a.saveLinkArticle(tx)
		} else if articleType == ArticleTypeGist {
			a.initGistArticle(values)
			a.saveGistArticle(tx)
		} else if articleType == ArticleTypeMarkdown {
			a.initMarkdownArticle(values)
			a.saveMarkdownArticle(tx)
		}

		dbg.Printf("Created: %v", created)
		dbg.Printf("Updated: %v", updated)

		return nil

	})

	return

}

/*
LoadArticle アーティクルを読み込む
*/
func LoadArticle(db *bolt.DB, articleID int) (*Article, error) {

	a := &Article{}

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(s2b(ArticleTypeBucket))
		articleType := b2i(b.Get(i2b(articleID)))

		// value set
		switch articleType {
		case ArticleTypeLink:
			a.loadLinkArticle(tx)
		case ArticleTypeGist:
			a.loadGistArticle(tx)
		case ArticleTypeMarkdown:
			a.loadMarkdownArticle(tx)
		}

		return nil
	})

	return a, err
}

/*
SaveArticle アーティクルを保存する
*/
func (a *Article) SaveArticle(db *bolt.DB) error {

	return db.Batch(func(tx *bolt.Tx) error {

		switch a.ArticleType {
		case ArticleTypeLink:
			a.saveLinkArticle(tx)
		case ArticleTypeGist:
			a.saveMarkdownArticle(tx)
		case ArticleTypeMarkdown:
			a.saveMarkdownArticle(tx)
		}

		return nil
	})

}

func (a *Article) newArticleID(db *bolt.DB) (int, error) {

	a.IDBucketName = ArticleIDBucket
	a.IDKey = "articleID"

	return a.newID(db)
}
