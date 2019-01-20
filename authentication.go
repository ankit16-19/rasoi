package main

import "net/http"

// AuthenticationMiddleware :
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// check if we need to update menu date
	if err := UpdateMenuDateIfWeekChange(); err != nil {
		//handle error
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// call next function when exiting middleware
		next(w, r)
	})
}
