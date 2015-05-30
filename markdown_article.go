package prism

import "github.com/boltdb/bolt"

const (
	ArticleMarkdownTextBucket = "ArticleMarkdownTextBucket"
)

type MarkdownArticle struct {
	ArticleInterface
	Text string
}

func (m *MarkdownArticle) initMarkdownArticle(values map[string]interface{}) {

	m.Text, _ = values["Text"].(string)

}

func (l *MarkdownArticle) loadMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	l.Text = b2s(b.Get(i2b(l.GetID())))

	return nil

}

func (l *MarkdownArticle) saveMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	b.Put(i2b(l.GetID()), s2b(l.Text))

	return nil

}
