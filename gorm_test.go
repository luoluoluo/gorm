package gorm

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db = newDB()

func TestCount(t *testing.T) {
	count, err := GetORM(db).Table("user").Count()
	if err != nil {
		t.Error(count, err)
		return
	}
}

func TestPaginate(t *testing.T) {
	count, rows, err := GetORM(db).Table("user").Select("User").Paginate(1, 10)
	if err != nil {
		t.Error(count, rows, err)
		return
	}
}
func TestRow(t *testing.T) {
	row, err := GetORM(db).Table("user").LeftJoin("user_info", "user.User=process_list.User").Select("user.User").NestHaving(func(orm *ORM) {
		orm.Having("user.User", "in", "luo1", "luo2").Having("user.Host", "=", 1)
	}).OrNestHaving(func(orm *ORM) {
		orm.Having("user.User", "in", "luo3", "luo4").Having("user.Host", "=", 2)
	}).Where("user.User", "=", "luo5").OrWhere("user.User", "=", "luo6").Order("user.Host DESC").Limit(0, 10).Row()
	if err != nil {
		t.Error(row, err)
		return
	}
}

func TestRows(t *testing.T) {
	rows, err := GetORM(db).Table("user").Select("User").NestWhere(func(orm *ORM) {
		orm.Where("User", "in", "luo1", "luo2").Where("Host", "=", 1)
	}).OrNestWhere(func(orm *ORM) {
		orm.Where("User", "in", "luo3", "luo4").Where("Host", "=", 2)
	}).Order("Host DESC").Limit(0, 10).Rows()
	if err != nil {
		t.Error(rows, err)
		return
	}
}

func TestUpdate(t *testing.T) {
	var data = map[string]interface{}{
		"User": "luo1",
	}
	rowsAffected, err := GetORM(db).Debug(true).Table("user").Where("User", "=", 11111).Update(data)
	if err != nil {
		t.Error(rowsAffected, err)
		return
	}
}

func TestInsert(t *testing.T) {
	var data = map[string]interface{}{
		"User": "luo1",
		"Host": 1,
	}
	id, err := GetORM(db).Table("user").Insert(data)
	if err != nil {
		t.Error(id, err)
		return
	}
}

func TestDelte(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	rowsAffected, err := GetORM(db).Debug(true).Table("user").Where("User", "=", "11111").Delete()
	if err != nil {
		tx.Rollback()
		t.Error(rowsAffected, err)
		return
	}
	tx.Commit()
}

func newDB() *sql.DB {
	db, err := sql.Open("mysql", "debian-sys-maint:uT4b7MkkcwdHN9TF@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
