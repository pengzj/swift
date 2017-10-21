package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var dbMap map[string]*sql.DB = make(map[string]*sql.DB)

func Register(name, dbType, dsn string)  {
	if dbMap[name] != nil {
		log.Fatal(name + "has registered twice")
	}
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

	dbMap[name] = db
}

func Get(name string) *sql.DB  {
	if dbMap[name] != nil {
		log.Fatal("name not exists")
	}
	return dbMap[name]
}

func CloseDB()  {
	for _, db := range dbMap {
		db.Close()
	}
}