package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"

	. "github.com/ad1619/rasoi/dao"
	. "github.com/ad1619/rasoi/models"
)

var mdao = MenuDAO{}

// GetAllMenu :
func GetAllMenu(w http.ResponseWriter, r *http.Request) {
	menu, err := mdao.FindAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, menu)
}

// CreateMenu : add new Menu
func CreateMenu(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var menu Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	menu.ID = bson.NewObjectId()

	// set Date for every entry
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	dates := WholeWeekDates(time.Now())

	for i := range days {
		reflect.ValueOf(&menu.MessUP).Elem().FieldByName(days[i]).FieldByName("Date").SetString(GetDateFromTime(dates[i]))
		reflect.ValueOf(&menu.MessDown).Elem().FieldByName(days[i]).FieldByName("Date").SetString(GetDateFromTime(dates[i]))
	}

	if err := mdao.Insert(menu); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, menu)
}

// UpdateMenu :
func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var menu Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// set Date for every entry
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	dates := WholeWeekDates(time.Now())

	for i := range days {
		reflect.ValueOf(&menu.MessUP).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
		reflect.ValueOf(&menu.MessDown).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
	}
	if err := mdao.Update(menu); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// UpdateMenuDateIfWeekChange :
func UpdateMenuDateIfWeekChange() error {
	menus, err := mdao.FindAll()
	if err != nil {
		return err
	}
	// compare last sunday date with curretn date
	if time.Now().After(menus[0].MessUP.Sun.Date) {
		if err := mdao.Update(menus[0]); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	mdao.Collection = "menu"
}
