package prism

import (
	"time"

	"github.com/boltdb/bolt"
)

const (
	CreateUpdatedLayout = "Jan 2, 2006 at 3:04pm (MST)"
	YesAccessible       = "yes"
	NoAccessible        = "no"
)

/*
EntityInterface Entityの基本的な情報にアクセスするためのインターフェース
*/
type EntityInterface interface {
	GetID() int
	Created() time.Time
	Updated() time.Time
}

/*
Entity 基本的なID振り機能を持たせるための構造
*/
type Entity struct {
	ID int

	Accessible string
	created    time.Time
	updated    time.Time
}

func (e *Entity) GetID() int {
	return e.ID
}

func (e *Entity) newID(tx *bolt.Tx, bucketName string, key string) error {

	b, _ := tx.CreateBucketIfNotExists(s2b(bucketName))

	lastID := b2i(b.Get(s2b(key)))
	newID := lastID + 1

	b.Put(s2b(key), i2b(newID))

	e.ID = newID

	return nil
}

func (e *Entity) Created(tx *bolt.Tx, bucketName string) time.Time {
	e.created = time.Now()
	created := e.created.Format(CreateUpdatedLayout)

	b := tx.Bucket(s2b(bucketName))
	b.Put(i2b(e.ID), s2b(created))

	return e.created
}

func (e *Entity) Updated(tx *bolt.Tx, bucketName string) time.Time {
	e.updated = time.Now()
	updated := e.updated.Format(CreateUpdatedLayout)

	b := tx.Bucket(s2b(bucketName))
	b.Put(i2b(e.ID), s2b(updated))

	return e.updated
}
