package dao

import (
	. "github.com/ad1619/rasoi/dbConnection"
	. "github.com/ad1619/rasoi/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CouponDAO :
type CouponDAO struct {
	Collection string
}

// FindAll :
func (c *CouponDAO) FindAll() ([]Coupon, error) {
	var coupons []Coupon
	err := Db.C(c.Collection).Find(bson.M{}).All(&coupons)
	return coupons, err
}

// FindByUserID :
func (c *CouponDAO) FindByUserID(id string) ([]Coupon, error) {
	var coupon []Coupon
	err := Db.C(c.Collection).Find(bson.M{"userid": id}).All(&coupon)
	return coupon, err
}

// FindByDateAndID :
func (c *CouponDAO) FindByDateAndID(id string, date string) (Coupon, error) {
	var coupon Coupon
	err := Db.C(c.Collection).Find(bson.M{"userid": id, "weekstartdate": date}).One(&coupon)
	return coupon, err
}

// Insert :
func (c *CouponDAO) Insert(coupon Coupon) error {
	cc := Db.C(c.Collection)
	index := mgo.Index{
		Key:    []string{"userid", "date"},
		Unique: true,
	}
	if err := cc.EnsureIndex(index); err != nil {
		return err
	}

	err := cc.Insert(&coupon)
	return err
}

// Delete :
func (c *CouponDAO) Delete(coupon Coupon) error {
	err := Db.C(c.Collection).Remove(&coupon)
	return err
}

// Update :
func (c *CouponDAO) Update(coupon Coupon) error {
	err := Db.C(c.Collection).UpdateId(coupon.ID, &coupon)
	return err
}
