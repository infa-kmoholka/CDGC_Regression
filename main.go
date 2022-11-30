package main

import (
	"log"
	"net/http"

	"github.com/infa-kmoholka/CDGC_Regression/apmservice"
)

func main() {

	apmservice.Init()
	log.Println("Starting server on port 4047")
	log.Println(http.ListenAndServe(":4048", nil))
}
