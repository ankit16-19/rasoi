package main

import (
	"fmt"
	"net/http"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes array to store all routes
type Routes []Route

// Index function to test API status
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status: live ; API version: v0.0.1")
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},

	// All Coupan Routes
	Route{"GetAllCoupon", "GET", "/coupon", GetAllCoupon},
	Route{"GetCouponByUserID", "GET", "/coupon/{userid}", GetCouponByUserID},
	Route{"GetCouponByDateAndID", "GET", "/coupon/{id}/{date}", GetCouponByDateAndID},
	Route{"CreateCoupon", "POST", "/coupon", CreateCoupon},
	Route{"UpdateCoupon", "PUT", "/coupon", UpdateCoupon},
	Route{"DeleteCouponByID", "DELETE", "/coupon/{id}", DeleteCouponByID},

	// All menu Routes
	Route{"GetAllMenu", "GET", "/menu", GetAllMenu},
	Route{"CreateMenu", "POST", "/menu", CreateMenu},
	Route{"UpdateMenu", "PUT", "/menu", UpdateMenu},
}
