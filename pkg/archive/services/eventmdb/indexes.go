// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package eventmdb

import "github.com/globalsign/mgo"

func (a *Archiver) createIdx() error {
	return a.createIdxEvents()
}

func (a *Archiver) createIdxEvents() error {
	c := a.getCollection(EventColName)
	indexes := []mgo.Index{
		{Key: []string{"created"}},
		{Key: []string{"code"}},
		{Key: []string{"level"}},
		{Key: []string{"$text:description"}},
	}
	for _, idx := range indexes {
		err := c.EnsureIndex(idx)
		if err != nil {
			return err
		}
	}
	return nil
}
