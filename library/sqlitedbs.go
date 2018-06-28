package library

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"errors"
	"fmt"
)

const (
	createdbsOrigin = `create table if not exists stocker (
id integer primary key autoincrement,
link text,
result integer,
search text,
context text,
created integer);`
	testingDbs = `CREATE TABLE IF NOT EXISTS tester (
	id INTEGER PRIMARY KEY autoincrement,
	password TEXT);`
	testInsert = `insert into tester (password) values(?);`
)

type Database struct {
	nameDbs string
	DB *sql.DB
}

func NewDatabase(nameDbs string) *Database {
	return &Database{
		nameDbs: nameDbs,
	}
}
func (d *Database) OpenDbs() error {
	db, err := sql.Open("sqlite3", d.nameDbs)
	if err != nil {
		return err
	}
	d.DB = db
	shmt, err := db.Prepare(testingDbs)
	if err != nil {
		return err
	}
	defer shmt.Close()
	res, err := shmt.Exec()
	if err != nil {
		return err
	}
	log.Print(res)
	return nil
}
func (d *Database) WriteRecordTester() error {
	if d.DB == nil {
		return errors.New("WARNING: DB instance is EMPTY\n")
	}
	sh, err  := d.DB.Prepare(testInsert)
	if err != nil {
		return err
	}
	defer sh.Close()
	for i:=0; i < 10; i++{
		res, err := sh.Exec(fmt.Sprint("Username#%d", i))
		if err != nil {
			log.Print(err)
		} else {
			log.Println(res)
		}
		ids, err := res.LastInsertId()
		if err != nil {
			log.Printf("Error %v\n", err)
		} else {
			log.Printf("IDS: %v\n", ids)
		}


	}
	return nil
}