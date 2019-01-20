package main

import (
	"log"
	"net/http"

	. "github.com/ad1619/rasoi/dbConnection"
)

func init() {
	var d = DAO{}
	d.Connect()
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":9000", router))

}
