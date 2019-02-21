package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	. "github.com/ankit16-19/rasoi/dao"
	. "github.com/ankit16-19/rasoi/models"
	"github.com/gorilla/mux"
)

var cdao = CouponDAO{}

// GetAllCoupon :
func GetAllCoupon(w http.ResponseWriter, r *http.Request) {
	// Move to authentication when feature added
	UpdateMenuDateIfWeekChange()
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
	params["id"] = strings.ToUpper(params["id"])
	coupon, err := cdao.FindByDateAndID(params["id"], params["date"])
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
	// Create ID for mongodb
	coupon.ID = bson.NewObjectId()
	coupon.Userid = strings.ToUpper(coupon.Userid)
	// Get next week first date
	t, err := time.Parse("2006-01-02", GetDateFromTime(FirstDayofWeek(time.Now().AddDate(0, 0, 7))))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	CalculateCouponPrice(&coupon)
	coupon.WeekStartDay = t
	// Add coupon to datebase
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
	coupon.Userid = strings.ToUpper(coupon.Userid)
	CalculateCouponPrice(&coupon)
	if err := cdao.Update(coupon); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeleteCouponByID :
func DeleteCouponByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	params["id"] = strings.ToUpper(params["id"])
	if err := cdao.DeleteByID(params["id"]); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func init() {
	cdao.Collection = "coupons"
}

// CalculateCouponPrice :
func CalculateCouponPrice(c *Coupon) {
	c.Amount1, c.Amount2, c.Total = 0, 0, 0
	price := []int{10, 25, 25}
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	daytime := []string{"Breakfast", "Lunch", "Dinner"}

	for _, day := range days {
		for j, dt := range daytime {
			// if Coupon selected
			if reflect.ValueOf(&c.Coupon).Elem().FieldByName(day).FieldByName(dt).FieldByName("IsSelected").Bool() {
				// if Messup
				if reflect.ValueOf(&c.Coupon).Elem().FieldByName(day).FieldByName(dt).FieldByName("IsMessup").Bool() {
					if j == 0 {
						c.Amount1 += price[j]
					} else if j == 1 {
						c.Amount1 += price[j]
					} else if j == 2 {
						c.Amount1 += price[j]
					}

				} else {
					if j == 0 {
						c.Amount2 += price[j]
					} else if j == 1 {
						c.Amount2 += price[j]
					} else if j == 2 {
						c.Amount2 += price[j]
					}
				}
			}
		}
	}
	c.Total = c.Amount1 + c.Amount2
}
