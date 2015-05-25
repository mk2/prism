package prism

import "github.com/boltdb/bolt"

const (
	ArticleMarkdownTextBucket = "ArticleMarkdownTextBucket"
)

type MarkdownArticle struct {
	EntityInterface
	Text string
}

func (m *MarkdownArticle) initMarkdownArticle(args ...interface{}) {

}

func (l *MarkdownArticle) saveMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleMarkdownTextBucket))

	b.Put(i2b(l.GetID()), s2b(l.Text))

	return nil

}
