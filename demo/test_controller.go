package demo

import (
	"go-learning/controllers"
	"net/http"
)

func RunTestController() {
	controllers.RegistRoutes()

	http.ListenAndServe(":8080", nil)
}
