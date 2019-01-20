package dao

import (
	. "github.com/ad1619/rasoi/dbConnection"
	. "github.com/ad1619/rasoi/models"
	"gopkg.in/mgo.v2/bson"
)

// MenuDAO :
type MenuDAO struct {
	Collection string
}

// FindAll :
func (c *MenuDAO) FindAll() ([]Menu, error) {
	var menus []Menu
	err := Db.C(c.Collection).Find(bson.M{}).All(&menus)
	return menus, err
}

// Update :
func (c *MenuDAO) Update(menu Menu) error {
	err := Db.C(c.Collection).UpdateId(menu.ID, &menu)
	return err
}

// Insert :
func (c *MenuDAO) Insert(menu Menu) error {
	err := Db.C(c.Collection).Insert(&menu)
	return err
}
