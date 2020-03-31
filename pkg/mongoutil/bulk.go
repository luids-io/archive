// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package mongoutil

import (
	"sync"

	"github.com/globalsign/mgo"
)

// Bulk is used for massive inserts
type Bulk struct {
	mutex sync.Mutex
	col   *mgo.Collection
	bulk  *mgo.Bulk
	count int
	size  int
}

// NewBulk returns a new bulk for collection with size
func NewBulk(c *mgo.Collection, size int) *Bulk {
	bk := &Bulk{
		col:  c,
		bulk: c.Bulk(),
		size: size,
	}
	return bk
}

// Insert a doc in the bulk
func (bk *Bulk) Insert(doc interface{}) error {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	bk.bulk.Insert(doc)
	bk.count++
	if bk.count >= bk.size {
		_, err := bk.bulk.Run()
		if err != nil {
			return err
		}
		bk.count = 0
		bk.bulk = bk.col.Bulk()
	}
	return nil
}

// Flush bulk
func (bk *Bulk) Flush() error {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	if bk.count > 0 {
		_, err := bk.bulk.Run()
		if err != nil {
			return err
		}
		bk.count = 0
		bk.bulk = bk.col.Bulk()
	}
	return nil
}
