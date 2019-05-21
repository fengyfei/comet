/*
 * Revision History:
 *     Initial: 2019/03/14        Yang ChengKai
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"

	mysql "github.com/fengyfei/comet/admin/model/mysql"
	"github.com/gin-gonic/gin"
)

// Controller -
type Controller struct {
	db *sql.DB
}

// New -
func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

// RegisterRouter -
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	name := "Admin"
	password := "111111"
	err := mysql.CreateTable(c.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", c.create)
	r.POST("/modify/email", c.modifyEmail)
	r.POST("/modify/mobile", c.modifyMobile)
	r.POST("/modify/password", c.modifyPassword)
	r.POST("/modify/active", c.ModifyAdminActive)
}

func (c *Controller) create(ctx *gin.Context) {
	var admin struct {
		Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
		Password string `json:"password"  binding:"printascii,max=30"`
	}

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	//Default password
	if admin.Password == "" {
		admin.Password = "111111"
	}

	err = mysql.Create(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"    binding:"required"`
			Email   string `json:"email"       binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyEmail(c.db, admin.AdminID, &admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyMobile(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"     binding:"required"`
			Mobile  string `json:"mobile"       binding:"required,numeric,len=11"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMobile(c.db, admin.AdminID, &admin.Mobile)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyPassword(ctx *gin.Context) {
	var (
		admin struct {
			AdminID     uint32 `json:"admin_id"      binding:"required"`
			Password    string `json:"password"      binding:"printascii,min=6,max=30"`
			NewPassword string `json:"newpassword"   binding:"printascii,min=6,max=30"`
			Confirm     string `json:"confirm"       binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.NewPassword == admin.Password {
		ctx.Error(err)
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	if admin.NewPassword != admin.Confirm {
		ctx.Error(err)
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.ModifyPassword(c.db, admin.AdminID, &admin.Password, &admin.NewPassword)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

//ModifyAdminActive -
func (c *Controller) ModifyAdminActive(ctx *gin.Context) {
	var (
		admin struct {
			CheckID     uint32 `json:"check_id"    binding:"required"`
			CheckActive bool   `json:"check_active"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyAdminActive(c.db, admin.CheckID, admin.CheckActive)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

//Login -
func (c *Controller) Login(ctx *gin.Context) (uint32, error) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		return 0, err
	}

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Password)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
