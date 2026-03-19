package controllers

import "net/http"

func registAboutRoutes() {
	http.HandleFunc("/about", aboutHandler)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About Page"))
}
