package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/TechCatsLab/gin-sor/controller"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3309)/")
	if err != nil {
		panic(err)
	}

	cnf := new(controller.Config)
	cnf.CategoryDB = "newstudents"
	cnf.CategoryTable = "user"

	controller.Register(dbConn, cnf, router)
	router.Run(":8088")
}
