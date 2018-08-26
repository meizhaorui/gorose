package main

import (
	"github.com/meizhaorui/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/meizhaorui/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.GetInstance()
	fmt.Println(db)
	res, err := db.Table("users").Where([][]interface{}{{"id", ">", 2}}).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)
}
