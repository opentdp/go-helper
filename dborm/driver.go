package dborm

import (
	"os"
	"path"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func useSqlite(args *Config) gorm.Dialector {

	dbname := args.DbName
	if args.Host != "" && !filepath.IsAbs(args.DbName) {
		dbname = path.Join(args.Host, dbname)
	}

	option := args.Option
	if option == "" {
		option = "?_pragma=busy_timeout=5000&_pragma=journa_mode(WAL)"
	}

	os.MkdirAll(filepath.Dir(dbname), 0755)

	return sqlite.Open(dbname + option)

}

func useMysql(args *Config) gorm.Dialector {

	host := args.Host
	user := args.User
	password := args.Password
	dbname := args.DbName
	option := args.Option

	if option == "" {
		option = "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return mysql.Open(user + ":" + password + "@tcp(" + host + ")/" + dbname + option)

}
