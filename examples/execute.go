package main

import (
	"github.com/meizhaorui/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/meizhaorui/gorose"
)

func main() {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	user, err := db.Execute("update users set name=? where id=?", "fizz8", 4)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(user)
}
