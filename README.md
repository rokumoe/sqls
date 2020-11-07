# sqls

go sql utility

## nonull

wrap type to a shadow type for skip NULL value from query result.

### usage

```go
func QueryNullable() {
	type User struct {
		Id   int `nonull:"-"`
		Name string
		Age  int
	}
	db := sqlx.MustOpen("sqlite3", ":memory:")
	defer db.Close()
	db.MustExec(`CREATE TABLE user (id int not null primary key, name varchar(10), age int)`)
	db.MustExec(`insert into user (id, name, age) values (?, ?, ?)`, 1, "user1", 10)
	db.MustExec(`insert into user (id, name) values (?, ?)`, 2, "user2")
	var user User
	err := db.Get(&user, `select * from user where id = ?`, 2)
	fmt.Println("try get user 2", err)
	err = db.Get(nonull.Wrap(&user), `select * from user where id = ?`, 2)
	fmt.Println("get user 2", err, user)
	var users []User
	err = db.Select(nonull.Wrap(&users), `select * from user`, 2)
	fmt.Println("get all users", err, users)
}
```
