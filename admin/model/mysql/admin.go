/*
 * Revision History:
 *     Initial: 2019/03/14        Yang ChengKai
 */

package mysql

import (
	"database/sql"
	"errors"

	"github.com/fengyfei/comet/pkgs/salt"
)

const (
	mysqlUserCreateTable = iota
	mysqlUserInsert
	mysqlUserLogin
	mysqlUserModifyEmail
	mysqlUserModifyMobile
	mysqlUserGetPassword
	mysqlUserModifyPassword
	mysqlUserModifyActive
	mysqlUserGetIsActive
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed  = errors.New("invalid username or password")

	adminSQLString = []string{
		`CREATE TABLE IF NOT EXISTS admin (
			admin_id    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name     	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			password 	VARCHAR(512) NOT NULL DEFAULT ' ',
			mobile   	VARCHAR(32) UNIQUE DEFAULT NULL,
			email    	VARCHAR(128) UNIQUE DEFAULT NULL,
			active   	BOOLEAN DEFAULT TRUE,
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (admin_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO admin (name,password,active)  VALUES (?,?,?)`,
		`SELECT admin_id,password FROM admin WHERE name = ? LOCK IN SHARE MODE`,
		`UPDATE admin SET email=? WHERE admin_id = ? LIMIT 1`,
		`UPDATE admin SET mobile=? WHERE admin_id = ? LIMIT 1`,
		`SELECT password FROM admin WHERE admin_id = ?  LOCK IN SHARE MODE`,
		`UPDATE admin SET password = ? WHERE admin_id = ? LIMIT 1`,
		`UPDATE admin SET active = ? WHERE admin_id = ? LIMIT 1`,
		`SELECT active FROM admin WHERE admin_id = ? LOCK IN SHARE MODE`,
	}
)

// CreateTable create admin table.
func CreateTable(db *sql.DB, name, password *string) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}

	Create(db, name, password)
	return nil
}

//Create create an administrative user
func Create(db *sql.DB, name, password *string) error {
	hash, err := salt.Generate(password)
	if err != nil {
		return err
	}

	result, err := db.Exec(adminSQLString[mysqlUserInsert], name, hash, true)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//Login the administrative user logins
func Login(db *sql.DB, name, password *string) (uint32, error) {
	var (
		id  uint32
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserLogin], name).Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}

	if !salt.Compare([]byte(pwd), password) {
		return 0, errLoginFailed
	}

	return id, nil
}

// ModifyEmail the administrative user updates email
func ModifyEmail(db *sql.DB, id uint32, email *string) error {

	result, err := db.Exec(adminSQLString[mysqlUserModifyEmail], email, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyMobile the administrative user updates mobile
func ModifyMobile(db *sql.DB, id uint32, mobile *string) error {

	result, err := db.Exec(adminSQLString[mysqlUserModifyMobile], mobile, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyPassword the administrative user updates password
func ModifyPassword(db *sql.DB, id uint32, password, newPassword *string) error {
	var (
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserGetPassword], id).Scan(&pwd)
	if err != nil {
		return err
	}

	if !salt.Compare([]byte(pwd), password) {
		return errLoginFailed
	}

	hash, err := salt.Generate(newPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec(adminSQLString[mysqlUserModifyPassword], hash, id)

	return err
}

//ModifyAdminActive the administrative user updates active
func ModifyAdminActive(db *sql.DB, id uint32, active bool) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyActive], active, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil

}

//IsActive return user.Active and nil if query success.
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var (
		isActive bool
	)

	db.QueryRow(adminSQLString[mysqlUserGetIsActive], id).Scan(&isActive)
	return isActive, nil
}
