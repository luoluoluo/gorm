### simple golang orm
```
db, err := sql.Open("mysql", "debian-sys-maint:uT4b7MkkcwdHN9TF@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
```
##### 一. SELECT
1. 多条记录
```
rows, err := GetORM(db).Table("user").Select("User").NestWhere(func(orm *ORM) {
		orm.Where("User", "in", "luo1", "luo2").Where("Host", "=", 1)
	}).OrNestWhere(func(orm *ORM) {
		orm.Where("User", "in", "luo3", "luo4").Where("Host", "=", 2)
	}).Order("Host DESC").Limit(0, 10).Rows()
```
2. 单条记录
```
row, err := GetORM(db).Table("user").LeftJoin("user_info", "user.User=process_list.User").Select("user.User").NestWhere(func(orm *ORM) {
		orm.Where("user.User", "in", "luo1", "luo2").Where("user.Host", "=", 1)
	}).OrNestWhere(func(orm *ORM) {
		orm.Where("user.User", "in", "luo3", "luo4").Where("user.Host", "=", 2)
	}).Where("user.User", "=", "luo5").OrWhere("user.User", "=", "luo6").Order("user.Host DESC").Row()
```
3. count 

```
SELECT COUNT(*) as count FROM user
```

```
count, err := GetORM(db).Table("user").Count()
```

4. 分页查询
```
count, rows, err := GetORM(db).Table("user").Select("User").Paginate(1, 10)
```

##### 二. UPDATE
```
var data = map[string]interface{}{
		"User": "luo1",
		"host": 1,
	}
rowsAffected, err := GetORM(db).Debug(true).Table("user").Where("User", "=", 11111).Update(data)
```

##### 三. INSERT
```
var data = map[string]interface{}{
		"User": "luo1",
		"Host": 1,
	}
id, err := GetORM(db).Table("user").Insert(data)
```

##### 四. DELETE
```
rowsAffected, err := GetORM(db).Debug(true).Table("user").Where("User", "=", "11111").Delete()
```