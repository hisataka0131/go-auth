package main

import (
	"go-auth/controller"
	"go-auth/db"
)

func main() {

	db.Init()
	defer db.Close()

	controller.Router()

}
