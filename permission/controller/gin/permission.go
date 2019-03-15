/*
 * Revision History:
 *     Initial: 2019/03/14        Yang ChengKai
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"

	mysql "github.com/fengyfei/comet/permission/model/mysql"
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

//RegisterRouter -
func (c *Controller) RegisterRouter(r gin.IRouter) {
	err := mysql.CreateTable(c.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/addrole", c.createRole)
	r.POST("/modifyrole", c.modifyRole)
	r.POST("/activerole", c.modifyRoleActive)
	r.POST("/getallrole", c.roleList)
	r.POST("/idgetrole", c.getRoleByID)

	r.POST("/addurl", c.addURLPermission)
	r.POST("/removeurl", c.removeURLPermission)
	r.POST("/urlgetrole", c.urlPermissions)
	r.POST("/geturl", c.permissions)

	r.POST("/addrelation", c.addRelation)
	r.POST("/removerelation", c.removeRelation)
	r.POST("/admingetrole", c.adminGetRoleMap)
	r.POST("/getalladmin", c.getAdminIDMap)
	r.POST("/getallroleid", c.getRoleIDMap)

}

func (c *Controller) createRole(ctx *gin.Context) {
	var (
		role struct {
			Name  string `json:"name"        binding:"required,alphanum,min=5,max=64"`
			Intro string `json:"intro"       binding:"required,alphanum,min=2,max=256"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.CreateRole(c.db, &role.Name, &role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyRole(ctx *gin.Context) {
	var (
		role struct {
			RoleID uint32 `json:"role_id"     binding:"required"`
			Name   string `json:"name"        binding:"required,alphanum,min=5,max=64"`
			Intro  string `json:"intro"       binding:"required,alphanum,min=2,max=256"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyRole(c.db, role.RoleID, &role.Name, &role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyRoleActive(ctx *gin.Context) {
	var (
		role struct {
			RoleID uint32 `json:"role_id"     binding:"required"`
			Active bool   `json:"active"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyRoleActive(c.db, role.RoleID, role.Active)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) roleList(ctx *gin.Context) {
	result, err := mysql.RoleList(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleList": result})
}

func (c *Controller) getRoleByID(ctx *gin.Context) {
	var (
		role struct {
			RoleID uint32 `json:"role_id"     binding:"required"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := mysql.GetRoleByID(c.db, role.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleByID": result})
}

func (c *Controller) addURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"         binding:"required"`
			RoleID uint32 `json:"role_id"     binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.AddURLPermission(c.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) removeURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"     binding:"required"`
			RoleID uint32 `json:"role_id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.RemoveURLPermission(c.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) urlPermissions(ctx *gin.Context) {
	var (
		url struct {
			URL string `json:"url"         binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := mysql.URLPermissions(c.db, &url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "URLPermissions": result})
}

func (c *Controller) permissions(ctx *gin.Context) {
	result, err := mysql.Permissions(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Permissions": result})
}

func (c *Controller) addRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID uint32 `json:"admin_id" binding:"required"`
			RoleID  uint32 `json:"role_id"  binding:"required"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.AddRelation(c.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) removeRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID uint32 `json:"admin_id" binding:"required"`
			RoleID  uint32 `json:"role_id"  binding:"required"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.RemoveRelation(c.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) adminGetRoleMap(ctx *gin.Context) {
	var (
		relation struct {
			AdminID uint32 `json:"admin_id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := mysql.AdminGetRoleMap(c.db, relation.AdminID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleMap": result})
}

func (c *Controller) getAdminIDMap(ctx *gin.Context) {
	result, err := mysql.GetAdminIDMap(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "AdminIDMap": result})
}

func (c *Controller) getRoleIDMap(ctx *gin.Context) {
	result, err := mysql.GetRoleIDMap(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleIDMap": result})
}
