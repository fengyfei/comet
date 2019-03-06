package admin

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fengyfei/comet/admin/model/mysql"
	"github.com/gin-gonic/gin"
)

// Controller -
type Controller struct {
	db             *sql.DB
	OnLoginSucceed func(userID uint32) error
}

// New -
func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

// RegisterRouter -
func (c *Controller) RegisterRouter(r gin.IRouter) {
	name := "admin"
	password := "111111"
	err := mysql.CreateTable(c.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/admin/create", c.create)
	r.POST("/api/v1/admin/login", c.login)
	r.POST("/api/v1/admin/email", c.modifyEmail)
	r.POST("/api/v1/admin/mobile", c.modifyMobile)
	r.POST("/api/v1/admin/newpassword", c.modifyPassword)
}

func (c *Controller) create(ctx *gin.Context) {
	var admin struct {
		Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
		Password string `json:"password"  binding:"printascii,min=6,max=30"`
	}

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	// Default password
	if admin.Password == "" {
		admin.Password = "111111"
	}

	err = mysql.Create(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) login(ctx *gin.Context) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	if c.OnLoginSucceed != nil {
		err = c.OnLoginSucceed(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			ID    uint32 `json:"id"    binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyEmail(c.db, admin.ID, &admin.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyMobile(ctx *gin.Context) {
	var (
		admin struct {
			ID     uint32 `json:"id"     binding:"required"`
			Mobile string `json:"mobile" binding:"required,numeric,len=11"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMobile(c.db, admin.ID, &admin.Mobile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyPassword(ctx *gin.Context) {
	var (
		admin struct {
			ID          uint32 `json:"id"            binding:"required"`
			Password    string `json:"password"      binding:"printascii,min=6,max=30"`
			NewPassword string `json:"newpassword"   binding:"printascii,min=6,max=30"`
			Confirm     string `json:"confirm"       binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.NewPassword == admin.Password {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	if admin.NewPassword != admin.Confirm {
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.ModifyPassword(c.db, admin.ID, &admin.Password, &admin.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
