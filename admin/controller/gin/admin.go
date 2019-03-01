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
	type admin struct {
		Name     string `json:"name"      validate:"required,alphanum,min=6,max=30"`
		Password string `json:"password"  validete:"printascii,min=6,max=30"`
	}

	first := admin{"Admin", "111111"}
	err := mysql.CreateTable(c.db, &first.Name, &first.Password)
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
		Name     string `json:"name"      validate:"required,alphanum,min=2,max=30"`
		Password string `json:"password"  validete:"printascii,min=6,max=30"`
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
			Name     string `json:"name" validate:"required,alphanum,min=2,max=30"`
			Password string `json:"pwd"  validete:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
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
			ID    uint32 `json:"id"    validate:"required"`
			Email string `json:"email" validate:"required,email"`
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
			ID     uint32 `json:"id"    validate:"required"`
			Mobile string `json:"mobile" validate:"required,numeric,len=11"`
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
			ID          uint32 `json:"id"            validate:"required"`
			Password    string `json:"password"      validete:"printascii,min=6,max=30"`
			NewPassword string `json:"newpassword"   validete:"printascii,min=6,max=30"`
			Confirm     string `json:"confirm"       validete:"printascii,min=6,max=30"`
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
