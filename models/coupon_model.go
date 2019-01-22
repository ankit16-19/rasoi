package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Coupon structure
type Coupon struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Userid       string        `bson:"userid" json:"userid"`
	Gender       string        `bson:"gender" json:"gender"`
	Amount1      int           `bson:"amount1" json:"amount1"`
	Amount2      int           `bson:"amount2" json:"mount2"`
	Total        int           `bson:"Total" json:"Total"`
	WeekStartDay time.Time     `bson:"weekstartdate" json:"weekstartdate"`
	Coupon       CouponForWeek `bson:"coupon" json:"coupon"`
}

// CouponForWeek : Coupon structure for whole week
type CouponForWeek struct {
	Mon CouponForDay `bson:"mon" json:"mon"`
	Tue CouponForDay `bson:"tue" json:"tue"`
	Wed CouponForDay `bson:"wed" json:"wed"`
	Thr CouponForDay `bson:"thr" json:"thr"`
	Fri CouponForDay `bson:"fri" json:"fri"`
	Sat CouponForDay `bson:"sat" json:"sat"`
	Sun CouponForDay `bson:"sun" json:"sun"`
}

// CouponForDay : Coupon structure for a day
type CouponForDay struct {
	Breakfast FoodType `bson:"breakfast" json:"breakfast"`
	Lunch     FoodType `bson:"lunch" json:"lunch"`
	Dinner    FoodType `bson:"dinner" json:"dinner"`
}

// FoodType : FoodStructure for a time
type FoodType struct {
	IsSelected bool `bson:"isSelected" json:"isSelected"`
	IsVeg      bool `bson:"isVeg" json:"isVeg"`
	IsMessup   bool `bson:"ismessup" json:"isMessUp"`
}
