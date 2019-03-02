package mysql

import (
	"database/sql"
	"errors"
	"time"
)

var (
	errInvaildInsert         = errors.New("insert comment: insert affected 0 rows")
	errInvalidChangeCategory = errors.New("change status: affected 0 rows")
)

//目录
type Category struct {
	CategoryId uint
	ParentId   uint //为0则是根目录
	Name       string
	Status     int8
	CreateTime time.Time
}

func CreateDB(db *sql.DB, createDB string) error {
	_, err := db.Exec(createDB)
	return err
}

func CreateTable(db *sql.DB, createTable string) error {
	_, err := db.Exec(createTable)
	return err
}

//自动设定 id 和 status状态和 创建时间
func InsertCategory(db *sql.DB, insert string, parentId uint, name string) (uint, error) {
	result, err := db.Exec(insert, parentId, name)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvaildInsert
	}

	categoryId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(categoryId), nil
}

//改变目录状态
func ChangeCategoryStatus(db *sql.DB, change string, category uint, status int8) error {
	result, err := db.Exec(change, status, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

//改变目录名称
func ChangeCategoryName(db *sql.DB, change string, category uint, name string) error {
	result, err := db.Exec(change, name, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

//显示一个父亲的所有孩子
func LisitChirldrenByParentId(db *sql.DB, Lisit string, parentId uint) ([]*Category, error) {
	var (
		categoryId uint
		name       string
		status     int8
		creatTime  time.Time

		categorys []*Category
	)

	rows, err := db.Query(Lisit, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&categoryId, &parentId, &name, &status, &creatTime); err != nil {
			return nil, err
		}

		category := &Category{
			CategoryId: categoryId,
			ParentId:   parentId,
			Name:       name,
			Status:     status,
			CreateTime: creatTime,
		}
		categorys = append(categorys, category)
	}

	return categorys, nil
}
