package prism

import "github.com/boltdb/bolt"

const (
	ArticleMarkdownTextBucket = "ArticleMarkdownTextBucket"
)

type MarkdownArticle struct {
	Text string
}

func (m *MarkdownArticle) initMarkdownArticle(args ...interface{}) {

}

func (l *Article) saveMarkdownArticle(tx *bolt.Tx) error {

	b := tx.Bucket(s2b(ArticleLinkURLBucket))

	b.Put(i2b(l.ID), s2b(l.Text))

	return nil

}
