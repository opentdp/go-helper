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

	os.MkdirAll(filepath.Dir(dbname), 0755)

	return sqlite.Open(dbname + args.Option)

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
