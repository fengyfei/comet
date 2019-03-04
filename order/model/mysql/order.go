package mysql

import (
	"database/sql"
	"time"
)

type Order struct {
	ID         uint32
	OrderCode  string    `json:"ordercode"`
	UserID     uint64    `json:"userid"`
	ShipCode   string    `json:"shipcode"`
	AddressID  string    `json:"addressid"`
	TotalPrice uint32    `json:"totalprice"`
	PayWay     uint8     `json:"payway"`
	Promotion  bool      `json:"promotion"`
	Freight    uint32    `json:"freight"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Closed     time.Time `json:"closed"`
	Updated    time.Time `json:"updated"`
}

type Item struct {
	ProductId uint32 `json:"productid"`
	OrderID   uint32 `json:"orderid"`
	Count     uint32 `json:"count"`
	Price     uint32 `json:"price"`
	Discount  uint32 `json:"discount"`
}

type OrmOrder struct {
	*Order
	Orm []*Item
}

func CreateDB(db *sql.DB, createDB string) error {
	_, err := db.Exec(createDB)
	return err
}

func CreateTable(db *sql.DB, createTable string) error {
	_, err := db.Exec(createTable)
	return err
}

func OrderIDByOrderCode(db *sql.DB, query string, ordercode string) (uint32, error) {
	var (
		orderid uint32
		err     error
	)

	rows, err := db.Query(query, ordercode)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&orderid); err != nil {
			return 0, err
		}
	}

	return orderid, nil
}

func SelectByOrderKey(db *sql.DB, query, queryitem string, orderid uint32) (*OrmOrder, error) {
	var (
		oo OrmOrder
		o  Order
	)
	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.UserID, &o.ShipCode, &o.AddressID, &o.TotalPrice, &o.PayWay, &o.Promotion, &o.Freight, &o.Status, &o.Created, &o.Closed, &o.Updated); err != nil {
			return nil, err
		}
	}

	oo.Order = &o
	oo.Orm, err = LisitItemByOrderId(db, queryitem, orderid)
	if err != nil {
		return nil, err
	}

	return &oo, nil
}

func LisitItemByOrderId(db *sql.DB, query string, orderid uint32) ([]*Item, error) {
	var (
		ProductId uint32
		OrderID   uint32
		Count     uint32
		Price     uint32
		Discount  uint32

		items []*Item
	)

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ProductId, &OrderID, &Count, &Price, &Discount); err != nil {
			return nil, err
		}

		item := &Item{
			ProductId: ProductId,
			OrderID:   OrderID,
			Count:     Count,
			Price:     Price,
			Discount:  Discount,
		}
		items = append(items, item)
	}

	return items, nil
}

func LisitOrderByUserId(db *sql.DB, query, queryitem string, userid uint64, mode uint8) ([]*OrmOrder, error) {
	var OOs []*OrmOrder

	rows, err := db.Query(query, userid, mode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var oo OrmOrder
		var o Order
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.UserID, &o.ShipCode, &o.AddressID, &o.TotalPrice, &o.PayWay, &o.Promotion, &o.Freight, &o.Status, &o.Created, &o.Closed, &o.Updated); err != nil {
			return nil, err
		}

		oo.Order = &o
		oo.Orm, err = LisitItemByOrderId(db, queryitem, oo.ID)
		if err != nil {
			return nil, err
		}
		OOs = append(OOs, &oo)
	}

	return OOs, nil
}
