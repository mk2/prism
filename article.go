package prism

import (
	"bytes"
	"github.com/boltdb/bolt"
)

const (
	ArticleTypeLink     = "link"
	ArticleTypeGist     = "gist"
	ArticleTypeMarkdown = "markdown"
)

const (
	ArticleIDBucket                 = "ArticleIdBucket"
	ArticleOwnerIDToArticleIDBucket = "ArticleIDToArticleOwnerIDBucket"
	ArticleIDToArticleOwnerIDBucket = "ArticleIDToArticleOwnerIDBucket"
	ArticleCreatedBucket            = "ArticleCreatedBucket"
	ArticleUpdatedBucket            = "ArticleUpdatedBucket"
	ArticleAccessibleBucket         = "ArticleAccessibleBucket"
	ArticleTypeBucket               = "ArticleTypeBucket"
)

type ArticleIface interface {
	GistArticleIface
}

type Article struct {
	ArticleIface

	Entity
	LinkArticle
	GistArticle
	MarkdownArticle

	articleType    string
	articleOwnerID string
}

func CreateArticleBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		ArticleIDBucket,
		ArticleOwnerIDToArticleIDBucket,
		ArticleCreatedBucket,
		ArticleUpdatedBucket,
		ArticleAccessibleBucket,
		ArticleGistIDBucket,
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

func DeleteArticleBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		ArticleIDBucket,
		ArticleOwnerIDToArticleIDBucket,
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
func NewArticle(db *bolt.DB, options map[string]interface{}) *Article {

	var a Article
	a.articleType, _ = options["ArticleType"].(string)
	a.articleOwnerID, _ = options["ArticleOwner"].(string)

	db.Update(func(tx *bolt.Tx) error {

		a.newArticleID(tx)
		a.saveArticleOwner(tx)
		a.saveArticleType(tx)

		created := a.Created(tx, ArticleCreatedBucket)
		updated := a.Updated(tx, ArticleUpdatedBucket)

		dbg.Printf("Created: %v", created)
		dbg.Printf("Updated: %v", updated)

		switch a.articleType {

		case ArticleTypeLink:
			a.LinkArticle.article = &a
			a.initLinkArticle(options)
			a.saveLinkArticle(tx)

		case ArticleTypeGist:
			a.GistArticle.article = &a
			a.initGistArticle(options)
			a.saveGistArticle(tx)

		case ArticleTypeMarkdown:
			a.MarkdownArticle.article = &a
			a.initMarkdownArticle(options)
			a.saveGistArticle(tx)

		}

		return nil

	})

	return &a
}

func LoadArticlesByOwnerID(db *bolt.DB, ownerID string) ([]*Article, error) {

	var as []*Article = make([]*Article, 10)
	var aids []string = make([]string, 10)

	db.View(func(tx *bolt.Tx) error {

		prefix := s2b("owner-" + ownerID)
		c := tx.Bucket(s2b(ArticleOwnerIDToArticleIDBucket)).Cursor()

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {

			aids = append(aids, b2s(v))

		}

		return nil
	})

	for _, aid := range aids {

		a, _ := LoadArticle(db, aid)
		as = append(as, a)

	}

	return as, nil
}

/*
LoadArticle アーティクルを読み込む
*/
func LoadArticle(db *bolt.DB, id string) (*Article, error) {

	var a Article
	a.id = id

	err := db.View(func(tx *bolt.Tx) error {

		a.loadArticleType(tx)
		a.loadArticleOwnerID(tx)

		dbg.Printf("ArticleType: %v", a.articleType)

		// value set
		switch a.articleType {

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

		switch a.articleType {

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

func (a *Article) saveArticleOwner(tx *bolt.Tx) error {

	var b *bolt.Bucket

	b = tx.Bucket(s2b(ArticleOwnerIDToArticleIDBucket))
	b.Put(s2b("owner-"+a.articleOwnerID+a.id), s2b(a.id))

	b = tx.Bucket(s2b(ArticleIDToArticleOwnerIDBucket))
	b.Put(s2b("article-"+a.id+a.articleOwnerID), s2b(a.articleOwnerID))

	return nil
}

func (a *Article) loadArticleOwnerID(tx *bolt.Tx) error {

	c := tx.Bucket(s2b(ArticleOwnerIDToArticleIDBucket)).Cursor()

	prefix := s2b("article-" + a.articleOwnerID)

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {

		a.articleOwnerID = b2s(v)
	}

	return nil
}

func (a *Article) saveArticleType(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleTypeBucket))
	b.Put(s2b(a.id), s2b(a.articleType))

	return nil
}

func (a *Article) loadArticleType(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleTypeBucket))
	a.articleType = b2s(b.Get(s2b(a.id)))

	return nil
}

func (a *Article) newArticleID(tx *bolt.Tx) error {

	return a.newID(tx, ArticleIDBucket, "articleID")
}
