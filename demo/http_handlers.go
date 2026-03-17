package demo

import "net/http"

func RunHttpHandlers() {
	handler := http.FileServer(http.Dir("public"))
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}
}

func test3() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public"+r.URL.Path)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
