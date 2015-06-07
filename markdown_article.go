package prism

import "github.com/boltdb/bolt"

const (
	ArticleMarkdownTextBucket = "ArticleMarkdownTextBucket"
)

type MarkdownArticle struct {
	article *Article
	Text    string
}

func (a *MarkdownArticle) initMarkdownArticle(values map[string]interface{}) {

	a.Text, _ = values["MarkdownText"].(string)
}

func (a *MarkdownArticle) loadMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	a.Text = b2s(b.Get(s2b(a.article.id)))

	return nil
}

func (a *MarkdownArticle) saveMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	b.Put(s2b(a.article.id), s2b(a.Text))

	return nil
}
