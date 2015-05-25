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

type EntityInterface interface {
	GetID() int
}

type Entity struct {
	EntityInterface
	ID           int
	IDBucketName string
	IDKey        string
	Created      string
	Updated      string
	Accessible   string
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

func (e *Entity) created() {
	t := time.Now()
	e.Created = t.Format(CreateUpdatedLayout)
}

func (e *Entity) updated() {
	t := time.Now()
	e.Updated = t.Format(CreateUpdatedLayout)

}
