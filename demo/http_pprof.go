package demo

import (
	"net/http"
	_ "net/http/pprof"
)

func RunHttpPprof() {
	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public"+r.URL.Path)
	})

	// http://localhost:9000/debug/pprof/
	// go tool pprof http://localhost:9000/debug/pprof/heap
	go http.ListenAndServe(":9000", nil)

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
