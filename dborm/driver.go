package dborm

import (
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func useSqlite(args *Param) gorm.Dialector {

	name := args.Name

	option := args.Option
	if option == "" {
		option = "?_pragma=busy_timeout=5000&_pragma=journa_mode(WAL)"
	}

	return sqlite.Open(name + option)

}

func useMysql(args *Param) gorm.Dialector {

	host := args.Host
	user := args.User
	passwd := args.Passwd
	name := args.Name

	option := args.Option
	if option == "" {
		option = "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return mysql.Open(user + ":" + passwd + "@tcp(" + host + ")/" + name + option)

}
