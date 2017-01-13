package main

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/05.3.html

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type VersionDatabase interface {
	close()
	insert(version, category, link string, date time.Time) int64
	fetch(version string) (VersionRow, bool)
	all() []VersionRow
	execute(sql string)
}

type VersionRow struct {
	uid      int
	version  string
	category string
	link     string
	date     time.Time
	created  time.Time
}

type FakeVersionDatabase struct {
	rows []VersionRow
}

func (d *FakeVersionDatabase) close() {

}
func (d *FakeVersionDatabase) insert(version, category, link string, date time.Time) int64 {
	for _, r := range d.rows {
		if r.version == version {
			return -1
		}
	}

	v := VersionRow{
		uid:      len(d.rows) + 1,
		version:  version,
		category: category,
		date:     date,
		link:     link,
		created:  time.Now(),
	}
	d.rows = append(d.rows, v)
	return int64(v.uid)
}
func (d *FakeVersionDatabase) fetch(version string) (VersionRow, bool) {
	for _, r := range d.rows {
		if r.version == version {
			return r, true
		}
	}

	v := VersionRow{}
	return v, false
}
func (d *FakeVersionDatabase) all() []VersionRow {
	return d.rows
}
func (d *FakeVersionDatabase) execute(sql string) {

}

type SqliteVersionDatabase struct {
	db *sql.DB
}

func NewDB(filename string) VersionDatabase {
	if len(filename) == 0 {
		return &FakeVersionDatabase{
			rows: []VersionRow{},
		}
	}

	db, err := sql.Open("sqlite3", filename)
	check(err)
	return &SqliteVersionDatabase{
		db: db,
	}
}

func (d *SqliteVersionDatabase) close() {
	d.db.Close()
	d.db = nil
}

func (d *SqliteVersionDatabase) insert(version, category, link string, date time.Time) int64 {
	stmt, err := d.db.Prepare("INSERT INTO versions(version, category, link, date) values(?,?,?,?)")
	check(err)

	res, err := stmt.Exec(version, category, link, date)
	if err != nil {
		return -1
	}

	id, err := res.LastInsertId()
	check(err)

	return id
}

func (d *SqliteVersionDatabase) fetch(version string) (VersionRow, bool) {
	stmt, err := d.db.Prepare("SELECT * FROM versions WHERE version = ?")
	check(err)

	rows, err := stmt.Query(version)
	check(err)

	var v VersionRow
	found := false
	for rows.Next() {
		err = rows.Scan(&v.uid, &v.version, &v.category, &v.link, &v.date, &v.created)
		check(err)

		found = true
		break
	}
	rows.Close()

	return v, found
}

func (d *SqliteVersionDatabase) all() []VersionRow {
	versions := []VersionRow{}

	rows, err := d.db.Query("SELECT * FROM versions")
	check(err)

	var v VersionRow
	for rows.Next() {
		err = rows.Scan(&v.uid, &v.version, &v.category, &v.link, &v.date, &v.created)
		check(err)
		versions = append(versions, v)
	}
	rows.Close()
	return versions
}

func (d *SqliteVersionDatabase) execute(sql string) {
	d.db.Exec(sql)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type DatabaseAccessor struct {
	db     VersionDatabase
	sender Sender
}

const (
	insertModeFinish       = 0
	insertModeForce        = 1
	insertModeCheckVersion = 2
)

func NewDBAccessor(db VersionDatabase, sender Sender) *DatabaseAccessor {
	return &DatabaseAccessor{
		db:     db,
		sender: sender,
	}
}

func (d *DatabaseAccessor) Run(modeCh chan int, rowCh chan VersionRow) {
	running := true
	for running == true {
		select {
		case m := <-modeCh:
			switch m {
			case insertModeFinish:
				log.Println("stop db accessor")
				running = false

			case insertModeForce:
				row := <-rowCh
				d.db.insert(row.version, row.category, row.link, row.date)

			case insertModeCheckVersion:
				row := <-rowCh
				_, found := d.db.fetch(row.version)
				if found {
					return
				}

				d.db.insert(row.version, row.category, row.link, row.date)
				msg := makeMessage(row.version, row.category, row.link)
				d.sender.send(msg)
				log.Printf("New version found : %s\n", row.version)
			}

		default:
			// CPU 100% 먹는거 방지
			interval := 1 * time.Second
			time.Sleep(interval)
		}
	}
}
