package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	. "github.com/ankit16-19/rasoi/dao"
	. "github.com/ankit16-19/rasoi/models"
	"github.com/gorilla/mux"
)

var cdao = CouponDAO{}

// GetAllCoupon :
func GetAllCoupon(w http.ResponseWriter, r *http.Request) {
	coupon, err := cdao.FindAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, coupon)
}

// GetCouponByUserID :
func GetCouponByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coupon, err := cdao.FindByUserID(params["userid"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Coupon ID")
		return
	}
	RespondWithJSON(w, http.StatusOK, coupon)
}

// GetCouponByDateAndID :
func GetCouponByDateAndID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coupon, err := cdao.FindByDateAndID(params["id"], params["weekstartdate"])
	if err != nil {
		// if response date is empty
		if err.Error() == "not found" {
			RespondWithJSON(w, http.StatusOK, coupon)
			return
		}
		RespondWithError(w, http.StatusBadRequest, "invalid Coupon ID and Date")
		fmt.Print(err)
		return
	}
	RespondWithJSON(w, http.StatusOK, coupon)
}

// CreateCoupon : add new Coupon
func CreateCoupon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var coupon Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	coupon.ID = bson.NewObjectId()
	coupon.WeekStartDay = FirstDayofWeek(time.Now().AddDate(0, 0, 7))
	if err := cdao.Insert(coupon); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, coupon)
}

// UpdateCoupon :
func UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var coupon Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := cdao.Update(coupon); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeleteCoupon :
func DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var coupon Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := cdao.Delete(coupon); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func init() {
	cdao.Collection = "coupons"
}
