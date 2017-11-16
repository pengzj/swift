package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pengzj/swift/logger"
)

var dbMap map[string]*sql.DB = make(map[string]*sql.DB)

func Register(name, dbType, dsn string)  {
	if dbMap[name] != nil {
		logger.Fatal(name, " has registered twice")
	}
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		logger.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		logger.Fatal("Error on opening database connection: ", err.Error())
	}

	dbMap[name] = db
}

func Get(name string) *sql.DB  {
	if dbMap[name] != nil {
		logger.Fatal("name not exists")
	}
	return dbMap[name]
}

func CloseDB()  {
	for _, db := range dbMap {
		db.Close()
	}
}