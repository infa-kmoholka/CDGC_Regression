package apmservice

import "net/http"

func Init() {

	
	http.HandleFunc("/compareRelease", compareRelease)
	http.HandleFunc("/test", test)
	http.HandleFunc("/htmlReport/", htmlReport)
	

}
