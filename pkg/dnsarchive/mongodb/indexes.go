// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package mongodb

func (a *Archiver) createIdx() error {
	return a.createIdxResolvs()
}

func (a *Archiver) createIdxResolvs() error {
	// c := a.getCollection(ResolvColName)
	// index := mgo.Index{
	// 	Key:        []string{"id"},
	// 	Unique:     true,
	// 	DropDups:   false,
	// 	Background: false,
	// 	Sparse:     false,
	// }
	// err := c.EnsureIndex(index)
	// if err != nil {
	// 	return err
	// }
	return nil
}
