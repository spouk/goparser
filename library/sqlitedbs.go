package library

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"errors"
	"fmt"
	"sync"
	"github.com/spouk/gocheats/utils"
	"time"
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
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	password TEXT);`

	testInsert = `INSERT INTO tester (password) VALUES(?);`

	testSelect = `select * from tester order by random() limit 1;`
)

type Database struct {
	nameDbs    string
	DB         *sql.DB
	sync.WaitGroup
	sync.RWMutex
	Stock      []*TestJob
	Randomizer *utils.Randomizer
}
type TestJob struct {
	Password string
}

func NewDatabase(nameDbs string) *Database {
	return &Database{
		nameDbs:    nameDbs,
		Randomizer: utils.NewRandomize(),
	}
}
func (d *Database) OpenDbs() error {
	db, err := sql.Open("sqlite3", d.nameDbs)
	if err != nil {
		return err
	}
	d.DB = db
	//set WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	_, err = db.Exec("PRAGMA cache_size = 2000;")
	_, err = db.Exec("PRAGMA default_cache_size = 2000;")
	_, err = db.Exec("PRAGMA page_size = 4096;")
	_, err = db.Exec("PRAGMA synchronous = NORMAL;")
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	_, err = db.Exec("PRAGMA busy_timeout = 1;")

	//create table
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
	sh, err := d.DB.Prepare(testInsert)
	if err != nil {
		return err
	}
	defer sh.Close()
	for i := 0; i < 10; i++ {
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
func (d *Database) generatorJob(countJob int) {
	for x := 0; x < countJob; x++ {
		d.Stock = append(d.Stock, &TestJob{Password: d.Randomizer.Letters()})
	}
}
func (d *Database) RunTestManager(countWorker int) {
	log.Println("runtestmanager....")
	go d.Manager(countWorker)
	log.Println("runtestmanager waiting....")
	d.Wait()
	log.Println("runtestmanager all ok")
}
func (d *Database) Manager(countWorker int) {
	log.Println("manager starting")
	//defer d.Done()
	var chanEnd = make(chan bool)
	var timer = time.NewTimer(time.Second * 10)

	for x := 0; x < countWorker; x ++ {
		d.Add(2)
		go d.worker(chanEnd, fmt.Sprintf("%v", x), "writerStock")
		go d.worker(chanEnd, fmt.Sprintf("%v", x), "readerStock")
	}

	<-timer.C
	log.Println("TIMER ENABLE, close all")
	close(chanEnd)
	d.Wait()
}
func (d *Database) worker(chanEnd chan bool, name string, role string) {
	var counterJobEnd = 0
	defer func() {
		d.Done()
		log.Println(fmt.Sprintf("[%s] [%3d] Worker `%v` starting", role, counterJobEnd, name))
	}()
	log.Println(fmt.Sprintf("[%s] Worker `%v` starting", role, name))
	for {
		select {
		case <-chanEnd:
			return
		default:
			//if len(d.Stock) == 0 {
			//	return
			//}
			switch role {
			case "writerStock":
				d.Lock()
				d.Stock = append(d.Stock, &TestJob{Password: d.Randomizer.Letters()})
				d.Unlock()
				counterJobEnd++
				time.Sleep(time.Second * 5)
			case "readerStock":
				if len(d.Stock) > 0 {
					d.Lock()
					element := d.Stock[0]
					d.Stock = append(d.Stock[0:], d.Stock[0+1:]...)
					d.Unlock()
					shtm, err := d.DB.Prepare(testInsert)
					if err != nil {
						log.Println(err)
					} else {
						_, err := shtm.Exec(element.Password)
						if err != nil {
							log.Println(err)
						}
						shtm.Close()
					}
					counterJobEnd++
					time.Sleep(time.Second * 2)
				} else {
					log.Printf("[%d] Stock empty\n", name)
					time.Sleep(time.Second * 1)
				}
			default:
				log.Println(fmt.Sprintf("[%s] [%s]  Worker `%v` starting", role, "ERROR WRONG ROLE", name))
				return
			}
		}
	}
}
