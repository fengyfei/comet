package service

import (
	"database/sql"

	"github.com/TechCatsLab/gin-sor/config"

	"github.com/TechCatsLab/gin-sor/model/mysql"
)

type TransactService struct {
	db   *sql.DB
	SQLS []string
}

const (
	mysqlCreateDatabase = iota
	mysqlCreateTable
	mysqlInsert
	mysqlUpdateStatus
	mysqlUpdateName
	mysqlSelectByParentID
)

func NewCategoryService(c *config.Config, db *sql.DB) *TransactService {
	ts := &TransactService{
		db: db,
		SQLS: []string{
			`CREATE DATABASE IF NOT EXISTS ` + c.CategoryDB,
			`CREATE TABLE IF NOT EXISTS ` + c.CategoryDB + `.` + c.CategoryTable + `(
				categoryId INT(11) NOT NULL AUTO_INCREMENT COMMENT '类别id',
				parentId INT(11) DEFAULT NULL  COMMENT '父类别id',
				name VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
				status TINYINT(1) DEFAULT '1' COMMENT '状态1-在售，2-废弃',
				createTime DATETIME DEFAULT current_timestamp COMMENT '创建时间',
				PRIMARY KEY (categoryId),
				INDEX(parentId)
			)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4`,
			`INSERT INTO ` + c.CategoryDB + `.` + c.CategoryTable + `(parentId,name) VALUES (?,?)`,
			`UPDATE ` + c.CategoryDB + `.` + c.CategoryTable + `SET status = ? WHERE categoryId = ? LIMIT 1`,
			`UPDATE ` + c.CategoryDB + `.` + c.CategoryTable + `SET name = ? WHERE categoryId = ? LIMIT 1`,
			`SELECT * FROM ` + c.CategoryDB + `.` + c.CategoryTable + ` WHERE parentId = ?`,
		},
	}
	return ts
}
func (ts *TransactService) CreateDB() error {
	return mysql.CreateDB(ts.db, ts.SQLS[mysqlCreateDatabase])
}

func (ts *TransactService) CreateTable() error {
	return mysql.CreateTable(ts.db, ts.SQLS[mysqlCreateTable])
}

//返回插入的编号
func (ts *TransactService) Insert(parentId uint, name string) (uint, error) {
	return mysql.InsertCategory(ts.db, ts.SQLS[mysqlInsert], parentId, name)
}

func (ts *TransactService) ChangeCategoryStatus(categoryId uint, status int8) error {
	return mysql.ChangeCategoryStatus(ts.db, ts.SQLS[mysqlUpdateStatus], categoryId, status)
}

func (ts *TransactService) ChangeCategoryName(categoryId uint, name string) error {
	return mysql.ChangeCategoryName(ts.db, ts.SQLS[mysqlUpdateName], categoryId, name)
}

//返回父级目录为parentId的目录
func (ts *TransactService) LisitChirldrenByParentId(parentId uint) ([]*mysql.Category, error) {
	categorys, err := mysql.LisitChirldrenByParentId(ts.db, ts.SQLS[mysqlSelectByParentID], parentId)
	if err != nil {
		return nil, err
	}

	return categorys, nil
}
