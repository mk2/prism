package prism

import "github.com/boltdb/bolt"

const (
	ArticleMarkdownTextBucket = "ArticleMarkdownTextBucket"
)

type MarkdownArticle struct {
	ArticleInterface
	Text string
}

func (a *MarkdownArticle) initMarkdownArticle(values map[string]interface{}) {

	a.Text, _ = values["Text"].(string)

}

func (a *MarkdownArticle) loadMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	a.Text = b2s(b.Get(s2b(a.GetID())))

	return nil

}

func (a *MarkdownArticle) saveMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	b.Put(s2b(a.GetID()), s2b(a.Text))

	return nil

}
