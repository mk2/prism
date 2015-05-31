package prism

import "github.com/boltdb/bolt"

const (
	ArticleTypeLink     = "link"
	ArticleTypeGist     = "gist"
	ArticleTypeMarkdown = "markdown"
)

const (
	ArticleIDBucket         = "ArticleIdBucket"
	ArticleCreatedBucket    = "ArticleCreatedBucket"
	ArticleUpdatedBucket    = "ArticleUpdatedBucket"
	ArticleAccessibleBucket = "ArticleAccessibleBucket"
	ArticleTypeBucket       = "ArticleTypeBucket"
)

type ArticleInterface interface {
	EntityInterface
}

type Article struct {
	Entity
	ArticleType string
	LinkArticle
	GistArticle
	MarkdownArticle

	Terms []*Meta
}

func CreateBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		ArticleIDBucket,
		ArticleCreatedBucket,
		ArticleUpdatedBucket,
		ArticleAccessibleBucket,
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
		ArticleCreatedBucket,
		ArticleUpdatedBucket,
		ArticleAccessibleBucket,
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
func NewArticle(db *bolt.DB, values map[string]interface{}) *Article {

	var a Article
	a.ArticleType, _ = values["ArticleType"].(string)

	db.Update(func(tx *bolt.Tx) error {

		a.newArticleID(tx)
		a.saveArticleType(tx)

		created := a.Created(tx, ArticleCreatedBucket)
		updated := a.Updated(tx, ArticleUpdatedBucket)

		dbg.Printf("Created: %v", created)
		dbg.Printf("Updated: %v", updated)

		switch a.ArticleType {

		case ArticleTypeLink:
			a.LinkArticle.article = &a
			a.initLinkArticle(values)
			a.saveLinkArticle(tx)

		case ArticleTypeGist:
			a.GistArticle.article = &a
			a.initGistArticle(values)
			a.saveGistArticle(tx)

		case ArticleTypeMarkdown:
			a.MarkdownArticle.article = &a
			a.initMarkdownArticle(values)
			a.saveGistArticle(tx)

		}

		return nil

	})

	return &a
}

/*
LoadArticle アーティクルを読み込む
*/
func LoadArticle(db *bolt.DB, ID string) (*Article, error) {

	var a Article
	a.ID = ID

	err := db.View(func(tx *bolt.Tx) error {

		a.loadArticleType(tx)

		dbg.Printf("ArticleType: %v", a.ArticleType)

		// value set
		switch a.ArticleType {

		case ArticleTypeLink:
			a.LinkArticle.article = &a
			a.loadLinkArticle(tx)

		case ArticleTypeGist:
			a.GistArticle.article = &a
			a.loadGistArticle(tx)

		case ArticleTypeMarkdown:
			a.MarkdownArticle.article = &a
			a.loadMarkdownArticle(tx)

		}

		dbg.Printf("Finish Reading")

		return nil
	})

	dbg.Printf("LoadArticle: %v", a)

	return &a, err
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

func (a *Article) saveArticleType(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleTypeBucket))
	b.Put(s2b(a.ID), s2b(a.ArticleType))

	return nil
}

func (a *Article) loadArticleType(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleTypeBucket))
	a.ArticleType = b2s(b.Get(s2b(a.ID)))

	return nil
}

func (a *Article) newArticleID(tx *bolt.Tx) error {

	return a.newID(tx, ArticleIDBucket, "articleID")
}
