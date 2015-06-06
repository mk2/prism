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
type EntityIface interface {
	GetID() string
	Created() time.Time
	Updated() time.Time
	IsVisible() bool
}

/*
Entity 基本的なID振り機能を持たせるための構造
*/
type Entity struct {
	id         string
	accessible string
	created    time.Time
	updated    time.Time
}

func (e *Entity) GetID() string {

	return e.id
}

func (e *Entity) IsVisible() bool {

	return e.accessible == YesAccessible
}

func (e *Entity) newID(tx *bolt.Tx, bucketName string, key string) error {

	b, _ := tx.CreateBucketIfNotExists(s2b(bucketName))

	lastID := b2i(b.Get(s2b(key)))
	newID := lastID + 1

	dbg.Printf("LastID: %d", lastID)
	dbg.Printf("NewID: %d", newID)

	b.Put(s2b(key), i2b(newID))

	e.id = i2s(newID)

	return nil
}

func (e *Entity) Created(tx *bolt.Tx, bucketName string) time.Time {

	e.created = time.Now()
	created := e.created.Format(CreateUpdatedLayout)

	b := tx.Bucket(s2b(bucketName))
	b.Put(s2b(e.id), s2b(created))

	return e.created
}

func (e *Entity) Updated(tx *bolt.Tx, bucketName string) time.Time {

	e.updated = time.Now()
	updated := e.updated.Format(CreateUpdatedLayout)

	b := tx.Bucket(s2b(bucketName))
	b.Put(s2b(e.id), s2b(updated))

	return e.updated
}
