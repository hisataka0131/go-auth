package controller

import (
	"github.com/gorilla/mux"
	"go-auth/service"
	"log"
	"net/http"
)

func Router() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", service.SignUp).Methods("POST")
	router.HandleFunc("/login", service.Login).Methods("POST")
	router.HandleFunc("/verify", service.TokenVerifyMiddleWare(service.VerifyEndpoint)).Methods("GET")

	// log.Fatal は、異常を検知すると処理の実行を止めてくれる
	log.Fatal(http.ListenAndServe(":8000", router))

}
