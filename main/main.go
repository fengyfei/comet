package main

import (
	"database/sql"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	admin "github.com/fengyfei/comet/admin/controller/gin"
	category "github.com/fengyfei/comet/category/controller/gin"
	order "github.com/fengyfei/comet/order/controller/gin"
	permission "github.com/fengyfei/comet/permission/controller/gin"
	upload "github.com/fengyfei/comet/upload/controller/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// JWTMiddleware should be exported for user authentication.
	JWTMiddleware *jwt.GinJWTMiddleware
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	a := admin.New(dbConn)

	JWTMiddleware = &jwt.GinJWTMiddleware{
		Realm:   "Template",
		Key:     []byte("hydra"),
		Timeout: 24 * time.Hour,
	}
	getUID := a.ExtendJWTMiddleWare(JWTMiddleware)

	router.POST("/api/v1/admin/login", JWTMiddleware.LoginHandler)

	router.Use(func(c *gin.Context) {
		JWTMiddleware.MiddlewareFunc()(c)
	})
	router.Use(admin.CheckActive(a, getUID))

	p := permission.New(dbConn)
	router.Use(permission.CheckPermission(p, getUID))

	u := upload.New(dbConn, "http://0.0.0.1:9573", getUID)

	a.RegisterRouter(router.Group("/api/v1/admin"))
	p.RegisterRouter(router.Group("/api/v1/permission"))
	u.RegisterRouter(router.Group("/api/v1/user"))

	category.Register(dbConn, "students", "test", router)

	order.Register(router, dbConn)

	router.Run(":8000")
}
