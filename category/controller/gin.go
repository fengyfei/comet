package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logging/logrus"

	"github.com/TechCatsLab/gin-sor/config"
	"github.com/TechCatsLab/gin-sor/service"
)

type Controller struct {
	service *service.TransactService
}

func Register(db *sql.DB, cnf *config.Config, r gin.IRouter) error {
	c := New(db, cnf)

	if err := c.CreateDB(); err != nil {
		return err
	}

	if err := c.CreateTable(); err != nil {
		return err
	}

	r.POST("/api/v1/category/create", c.Insert)
	r.POST("/api/v1/category/modify/status", c.ChangeCategoryStatus)
	r.POST("/api/v1/category/modify/name", c.ChangeCategoryName)
	r.POST("/api/v1/category/children", c.LisitChirldrenByParentId)
	return nil
}

func New(db *sql.DB, c *config.Config) *Controller {
	return &Controller{
		service: service.NewCategoryService(c, db),
	}
}

func (con *Controller) CreateDB() error {
	return con.service.CreateDB()
}

func (con *Controller) CreateTable() error {
	return con.service.CreateTable()
}

func (con *Controller) Insert(c *gin.Context) {
	var (
		req struct {
			ParentId uint   `json:"parentId"`
			Name     string `json:"name"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := con.service.Insert(req.ParentId, req.Name)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Insert", "id": id})
	return
}

func (con *Controller) ChangeCategoryStatus(c *gin.Context) {
	var (
		req struct {
			CategoryId uint `json:"categoryId"`
			Status     int8 `json:"status"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := con.service.ChangeCategoryStatus(req.CategoryId, req.Status)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ChangeCategoryStatus"})
	return
}

func (con *Controller) ChangeCategoryName(c *gin.Context) {
	var (
		req struct {
			CategoryId uint   `json:"categoryId"`
			Name       string `json:"name"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := con.service.ChangeCategoryName(req.CategoryId, req.Name)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ChangeCategoryName"})
	return
}

func (con *Controller) LisitChirldrenByParentId(c *gin.Context) {
	var (
		req struct {
			ParentId uint `json:"parentId"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categorys, err := con.service.LisitChirldrenByParentId(req.ParentId)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "LisitChirldrenByParentId", "category": categorys})
	return
}
