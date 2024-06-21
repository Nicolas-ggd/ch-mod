package main

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db"
	"github.com/Nicolas-ggd/ch-mod/pkg/api"
)

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		panic(err)
	}

	serve := api.ServeAPI(dbConn)
	serve.Run(":8090")
}
