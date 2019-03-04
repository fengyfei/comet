package order

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/base"
	"github.com/base/constants"
	"github.com/fengyfei/comet/order/model/mysql"
	"github.com/gin-gonic/gin"
	"github.com/logging/logrus"
)

type Stocker interface {
	ModifyProductStock(tx *sql.Tx, targetID uint32, num int) error
}

type UserChecker interface {
	UserCheck(tx *sql.Tx, userid uint64, productID uint32) error
}

type Config struct {
	OrderDB        string
	OrderTable     string
	ItemTable      string
	ClosedInterval int

	Stock Stocker
	User  UserChecker
}

type Controller struct {
	db  *sql.DB
	Cnf Config
}

func Register(r gin.IRouter, db *sql.DB, cnf Config) error {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	c := New(db, cnf)

	if err := c.CreateOrderTable(); err != nil {
		return err
	}

	if err := c.CreateItemTable(); err != nil {
		return err
	}

	r.POST("/api/v1/order/create", c.Insert)
	r.POST("/api/v1/order/info", c.OrderInfoByOrderID)
	r.POST("/api/v1/order/user", c.LisitOrderByUserIDAndStatus)
	r.POST("/api/v1/order/id", c.OrderIDByOrderCode)

	return nil
}

func New(db *sql.DB, cnf Config) *Controller {
	return &Controller{
		db:  db,
		Cnf: cnf,
	}
}

func (ctl *Controller) CreateDB() error {
	return ctl.service.CreateDB()
}

func (ctl *Controller) CreateOrderTable() error {
	return ctl.service.CreateOrderTable()
}
func (ctl *Controller) CreateItemTable() error {
	return ctl.service.CreateItemTable()
}

func (ctl *Controller) Insert(c *server.Context) error {
	var (
		req struct {
			UserID     uint64 `json:"userid"`
			AddressID  string `json:"addressid"`
			TotalPrice uint32 `json:"totalprice"`
			Promotion  string `json:"promotion"`
			Freight    uint32 `json:"freight"`

			Items []mysql.Item `json:"items"`
		}
		rep struct {
			ordercode string
			orderid   uint32
		}
		err error
	)
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	promotion, err := strconv.ParseBool(req.Promotion)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	times := time.Now()
	rep.ordercode = strconv.Itoa(times.Year()) + strconv.Itoa(int(times.Month())) + strconv.Itoa(times.Day()) + strconv.Itoa(times.Hour()) + strconv.Itoa(times.Minute()) + strconv.Itoa(times.Second()) + strconv.Itoa(int(req.UserID))
	order := mysql.Order{
		OrderCode:  rep.ordercode,
		UserID:     req.UserID,
		AddressID:  req.AddressID,
		TotalPrice: req.TotalPrice,
		Promotion:  promotion,
		Freight:    req.Freight,
		Created:    times,
	}

	rep.orderid, err = ctl.service.Insert(order, req.Items)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrCreateInMysql, nil))
	}
	return c.ServeJSON(base.RespStatusAndIDCODEData(constants.ErrSucceed, rep.orderid, rep.ordercode))
}

//optional
func (ctl *Controller) OrderIDByOrderCode(c *server.Context) error {
	var req struct {
		Ordercode string `json:"ordercode"`
	}
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	id, err := ctl.service.OrderIDByOrderCode(req.Ordercode)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}
	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))
}

//full info for One Order
func (ctl *Controller) OrderInfoByOrderID(c *server.Context) error {
	var req struct {
		OrderId uint32 `json:"orderid"`
	}
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	rep, err := ctl.service.OrderInfoByOrderKey(req.OrderId)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}
	return c.ServeJSON(base.RespStatusAndTwoData(constants.ErrSucceed, rep.Order, rep.Orm))
}

/*
mode:
  Unfinished = 0
  Finished   = 1
  Paid       = 2
  Consigned  = 3
  Canceled   = 4
*/
func (ctl *Controller) LisitOrderByUserIDAndStatus(c *server.Context) error {
	var req struct {
		Userid uint64 `json:"userid"`
		Status uint8  `json:"status"`
	}

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	orders, err := ctl.service.LisitOrderByUserId(req.Userid, req.Status)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, orders))
}
