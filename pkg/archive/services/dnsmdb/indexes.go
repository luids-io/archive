// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package dnsmdb

import "github.com/globalsign/mgo"

func (a *Archiver) createIdx() error {
	return a.createIdxResolvs()
}

func (a *Archiver) createIdxResolvs() error {
	c := a.getCollection(ResolvColName)
	indexes := []mgo.Index{
		{Key: []string{"timestamp"}},
		{Key: []string{"serverIP"}},
		{Key: []string{"clientIP"}},
		{Key: []string{"name"}},
		{Key: []string{"resolvedIPs"}},
		{Key: []string{"tldPlusOne"}},
	}
	for _, idx := range indexes {
		err := c.EnsureIndex(idx)
		if err != nil {
			return err
		}
	}
	return nil
}
