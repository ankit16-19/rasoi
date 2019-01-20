package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Menu structure
type Menu struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	MessUP   MenuForWeek   `bson:"messUP" json:"messUP"`
	MessDown MenuForWeek   `bson:"messDown" json:"messDown"`
}

// MenuForWeek : Menu structure for whole week
type MenuForWeek struct {
	Mon MenuForDay `bson:"mon" json:"mon"`
	Tue MenuForDay `bson:"tue" json:"tue"`
	Wed MenuForDay `bson:"wed" json:"wed"`
	Thr MenuForDay `bson:"thr" json:"thr"`
	Fri MenuForDay `bson:"fri" json:"fri"`
	Sat MenuForDay `bson:"sat" json:"sat"`
	Sun MenuForDay `bson:"sun" json:"sun"`
}

// MenuForDay : Menu structure for a day
type MenuForDay struct {
	Breakfast string    `bson:"breakfast" json:"breakfast"`
	Lunch     string    `bson:"lunch" json:"lunch"`
	Dinner    string    `bson:"dinner" json:"dinner"`
	Date      time.Time `bson:"date" json:"date"`
}
