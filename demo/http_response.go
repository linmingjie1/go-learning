package demo

import (
	"fmt"
	"net/http"
)

func RunHttpResponse() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		str := `<!doctype html>
<html>
<head> <title>Hello World</title> </head>
<body> <h1>Hello World</h1> </body>
</html>`
		_, err := w.Write([]byte(str))
		if err != nil {
			return
		}
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, err := fmt.Fprintln(w, "Internal Server Error")
		if err != nil {
			return
		}
	})

	http.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "Hello World"}`))
		if err != nil {
			return
		}
	})

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
