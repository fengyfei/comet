package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	admin "github.com/fengyfei/comet/admin/controller/gin"
	category "github.com/fengyfei/comet/category/controller/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	c := admin.New(dbConn)
	c.RegisterRouter(router)

	categoryContoller := category.Register(dbConn, "students", "test", router)

	router.Run(":8000")
}
