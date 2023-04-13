package dborm

import (
	"github.com/open-tdp/go-helper/logman"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

type Config struct {
	Type   string
	Host   string
	User   string
	Passwd string
	Name   string
	Option string
}

func Connect(args *Config) {

	config := &gorm.Config{
		Logger: newLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	if db, err := gorm.Open(dialector(args), config); err != nil {
		logman.Fatal("Connect to databse failed", "error", err)
	} else {
		Db = db
	}

}

func dialector(args *Config) gorm.Dialector {

	switch args.Type {
	case "sqlite":
		return useSqlite(args)
	case "mysql":
		return useMysql(args)
	default:
		logman.Fatal("Database type error", "type", args.Type)
	}

	return nil

}
