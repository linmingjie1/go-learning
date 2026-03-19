package controllers

import "net/http"

func registHomeRoutes() {
	http.HandleFunc("/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home Page"))
}
