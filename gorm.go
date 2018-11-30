package gorm

import (
	"database/sql"
	"fmt"
	"strings"
)

// ORM orm struct
type ORM struct {
	db     *sql.DB
	debug  bool
	table  string
	cols   string
	where  string
	having string
	order  string
	join   string
	group  string
	limit  string
	args   []interface{}
}

// GetORM sqlx.GetORM(orm)
func GetORM(db *sql.DB) *ORM {
	return &ORM{db: db, cols: "*"}
}

// Debug sqlx.Debug(true)
func (orm *ORM) Debug(debug bool) *ORM {
	orm.debug = debug
	return orm
}

// Table orm.Table("user")
func (orm *ORM) Table(table string) *ORM {
	orm.table = table
	return orm
}

// Select orm.Select("name, age"])
func (orm *ORM) Select(cols string) *ORM {
	orm.cols = cols
	return orm
}

// Order  orm.Order("id desc")
func (orm *ORM) Order(query string) *ORM {
	orm.order = fmt.Sprintf("ORDER BY %s", query)
	return orm
}

// Group orm.Group("name")
func (orm *ORM) Group(query string) *ORM {
	orm.group = fmt.Sprintf("GROUP BY %s", query)
	return orm
}

// Limit orm.Limit(0, 10)
func (orm *ORM) Limit(offset, count int) *ORM {
	orm.limit = fmt.Sprintf("LIMIT %d, %d", offset, count)
	return orm
}

// LeftJoin orm.LeftJoin("user_info", "userinfo.user_id=user.id")
func (orm *ORM) LeftJoin(table string, on string) *ORM {
	orm.join += fmt.Sprintf("LEFT JOIN %s ON %s", table, on)
	return orm
}

// RightJoin orm.RightJoin("user_info", "userinfo.user_id=user.id")
func (orm *ORM) RightJoin(table string, on string) *ORM {
	orm.join += fmt.Sprintf("RIGHT JOIN %s ON %s", table, on)
	return orm
}

// InnerJoin orm.InnerJoin("user_info", "userinfo.user_id=user.id")
func (orm *ORM) InnerJoin(table string, on string) *ORM {
	orm.join += fmt.Sprintf("INNER JOIN %s ON %s", table, on)
	return orm
}

// NestWhere e: (name="luo1" and sex=1) or (name="luo2" and sex=2)
// orm.NestWhere(func(orm){
//		orm.Where("name", "=", "luo1").Where("sex", 1)
// }).OrNestWhere(func(orm)){
//		orm.Where("name" "=", "luo2").Where("sex", 2)
// })
func (orm *ORM) NestWhere(f func(*ORM)) *ORM {
	if orm.where != "" {
		orm.where += fmt.Sprintf(" %s ", "AND")
	}
	orm.where += "("
	f(orm)
	orm.where += ")"
	return orm
}

// OrNestWhere e: (name="luo1" and sex=1) or (name="luo2" and sex=2)
// orm.NestWhere(func(orm){
//		orm.Where("name", "=", "luo1").Where("sex", 1)
// }).OrNestWhere(func(orm)){
//		orm.Where("name" "=", "luo2").Where("sex", 2)
// })
func (orm *ORM) OrNestWhere(f func(*ORM)) *ORM {
	if orm.where != "" {
		orm.where += fmt.Sprintf(" %s ", "OR")
	}
	orm.where += "("
	f(orm)
	orm.where += ")"
	return orm
}

// NestHaving e: (name="luo1" and sex=1) or (name="luo2" and sex=2)
// orm.NestWhere(func(orm){
//		orm.Where("name", "=", "luo1").Where("sex", 1)
// }).OrNestWhere(func(orm)){
//		orm.Where("name" "=", "luo2").Where("sex", 2)
// })
func (orm *ORM) NestHaving(f func(*ORM)) *ORM {
	if orm.having != "" {
		orm.having += fmt.Sprintf(" %s ", "AND")
	}
	orm.having += "("
	f(orm)
	orm.having += ")"
	return orm
}

// OrNestHaving e: (name="luo1" and sex=1) or (name="luo2" and sex=2)
// orm.NestWhere(func(orm){
//		orm.Where("name", "=", "luo1").Where("sex", 1)
// }).OrNestWhere(func(orm)){
//		orm.Where("name" "=", "luo2").Where("sex", 2)
// })
func (orm *ORM) OrNestHaving(f func(*ORM)) *ORM {
	if orm.having != "" {
		orm.having += fmt.Sprintf(" %s ", "OR")
	}
	orm.having += "("
	f(orm)
	orm.having += ")"
	return orm
}

// Where orm.Where("name", "=", "luo")
func (orm *ORM) Where(col string, op string, args ...interface{}) *ORM {
	return orm.buildCondition(false, "and", col, op, args...)
}

// OrWhere orm.OrWhere("name", "luo1", "luo2")
func (orm *ORM) OrWhere(col string, op string, args ...interface{}) *ORM {
	return orm.buildCondition(false, "or", col, op, args...)
}

// Having orm.Having("name", "=", "luo")
func (orm *ORM) Having(col string, op string, args ...interface{}) *ORM {
	return orm.buildCondition(true, "and", col, op, args...)
}

// OrHaving orm.OrHaving("name", "in", "luo1", "luo2")
func (orm *ORM) OrHaving(col string, op string, args ...interface{}) *ORM {
	return orm.buildCondition(true, "or", col, op, args...)
}

func (orm *ORM) buildCondition(isHaving bool, andor string, col string, op string, args ...interface{}) *ORM {
	var condition string
	op = strings.ToUpper(op)
	var places []string
	for _, arg := range args {
		places = append(places, "?")
		orm.args = append(orm.args, arg)
	}
	switch op {
	case "IN", "NOT IN":
		condition = fmt.Sprintf("%s %s(%s)", col, op, strings.Join(places, ","))
	case "BETWEEN", "NOT BETWEEN":
		condition = fmt.Sprintf("%s %s ? AND ?", col, op)
	default:
		condition = fmt.Sprintf("%s%s?", col, op)
	}
	if isHaving {
		if orm.having != "" && orm.having[len(orm.having)-1:] != "(" {
			condition = fmt.Sprintf("%s %s", strings.ToUpper(andor), condition)
		}
		if orm.having == "" {
			orm.having = fmt.Sprintf("HAVING %s", condition)
		} else if orm.having == "(" {
			orm.having = fmt.Sprintf("HAVING (%s", condition)
		} else if orm.having[len(orm.having)-1:] == "(" {
			orm.having += fmt.Sprintf("%s", condition)
		} else {
			orm.having += fmt.Sprintf(" %s", condition)
		}
	} else {
		if orm.where != "" && orm.where[len(orm.where)-1:] != "(" {
			condition = fmt.Sprintf("%s %s", strings.ToUpper(andor), condition)
		}
		if orm.where == "" {
			orm.where = fmt.Sprintf("WHERE %s", condition)
		} else if orm.where == "(" {
			orm.where = fmt.Sprintf("WHERE (%s", condition)
		} else if orm.where[len(orm.where)-1:] == "(" {
			orm.where += fmt.Sprintf("%s", condition)
		} else {
			orm.where += fmt.Sprintf(" %s", condition)
		}
	}
	return orm
}

// Insert orm.Insert(["name": "luo", "age": 1])
func (orm *ORM) Insert(data map[string]interface{}) (int64, error) {
	var cols []string
	var args []interface{}
	var places []string
	for k, v := range data {
		cols = append(cols, k)
		places = append(places, "?")
		args = append(args, v)
	}
	var query = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", orm.table, strings.Join(cols, ","), strings.Join(places, ","))
	res, err := orm.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update orm.Update(["name": "luo", "age": 1])
func (orm *ORM) Update(data map[string]interface{}) (int64, error) {
	var cols []string
	var args []interface{}
	for k, v := range data {
		cols = append(cols, k+"=?")
		args = append(args, v)
	}
	orm.args = append(args, orm.args...)
	var query = fmt.Sprintf("UPDATE %s SET %s %s", orm.table, strings.Join(cols, ","), orm.where)
	res, err := orm.Exec(query, orm.args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete delete record
func (orm *ORM) Delete() (int64, error) {
	var query = fmt.Sprintf("DELETE FROM %s %s", orm.table, orm.where)
	res, err := orm.Exec(query, orm.args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Row select one
func (orm *ORM) Row() (map[string]interface{}, error) {
	row := make(map[string]interface{})
	rows, err := orm.Rows()
	if err != nil {
		return row, err
	}
	if len(rows) != 0 {
		row = rows[0]
	}
	return row, nil
}

// Query orm.Query("SELECT * FROM user WHERE name=?", "luo")
func (orm *ORM) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	if orm.debug {
		fmt.Printf(query, args...)
	}
	stmt, err := orm.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return res, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return res, err
	}

	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return res, err
		}
		vmap := make(map[string]interface{}, len(scanArgs))
		for i, col := range values {
			vmap[cols[i]] = col
		}
		res = append(res, vmap)
	}
	return res, nil
}

// Exec orm.Exec("DELETE FROM user WHERE name=?", "luo")
func (orm *ORM) Exec(query string, args ...interface{}) (sql.Result, error) {
	if orm.debug {
		fmt.Printf(query, args...)
	}
	stmt, err := orm.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

// Rows select list
func (orm *ORM) Rows() ([]map[string]interface{}, error) {
	var query = fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s %s", orm.cols, orm.table, orm.join, orm.where, orm.having, orm.order, orm.group, orm.limit)
	return orm.Query(query, orm.args...)
}

// Count count(*)
func (orm *ORM) Count() (int64, error) {
	orm.cols = "COUNT(*) as count"
	row, err := orm.Row()
	if err != nil {
		return 0, err
	}
	return row["count"].(int64), nil
}

// Paginate orm.Paginate(1, 10)
func (orm *ORM) Paginate(page int, size int) (int64, []map[string]interface{}, error) {
	var offset = (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	rows, err := orm.Limit(offset, size).Rows()
	if err != nil {
		return 0, nil, err
	}
	orm.limit = ""
	count, err := orm.Count()
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}
