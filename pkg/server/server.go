package server 

import "net/http"

func Run() error {
	
	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(":7540", nil)
	
}