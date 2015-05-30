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
	ID           int
	IDBucketName string
	IDKey        string

	Accessible string

	created string
	updated string
}

func (e *Entity) GetID() int {
	return e.ID
}

func (e *Entity) saveEntity(db *bolt.DB) {

}

func (e *Entity) newID(db *bolt.DB) (int, error) {

	tx, _ := db.Begin(true)
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists(s2b(e.IDBucketName))

	lastID := b2i(b.Get(s2b(e.IDKey)))
	newID := lastID + 1

	b.Put(s2b(e.IDKey), i2b(newID))

	result := tx.Commit()

	e.ID = newID

	return newID, result
}

func (e *Entity) Created() time.Time {
	t := time.Now()
	e.created = t.Format(CreateUpdatedLayout)

	return t
}

func (e *Entity) Updated() time.Time {
	t := time.Now()
	e.updated = t.Format(CreateUpdatedLayout)

	return t
}
