package main

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/05.3.html

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type VersionRow struct {
	UID      int       `json:"uid"`
	Version  string    `json:"version"`
	Category string    `json:"category"`
	Link     string    `json:"link"`
	Date     time.Time `json:"date"`
	Created  time.Time `json:"created"`
}

type VersionDatabase struct {
	db *sql.DB
}

func NewDB(filename string) *VersionDatabase {
	db, err := sql.Open("sqlite3", filename)
	check(err)
	return &VersionDatabase{
		db: db,
	}
}

func (d *VersionDatabase) Close() {
	d.db.Close()
	d.db = nil
}

func (d *VersionDatabase) Insert(r *VersionRow) int64 {
	stmt, err := d.db.Prepare("INSERT INTO versions(version, category, link, date) values(?,?,?,?)")
	check(err)

	res, err := stmt.Exec(r.Version, r.Category, r.Link, r.Date)
	if err != nil {
		return -1
	}

	id, err := res.LastInsertId()
	check(err)

	return id
}

func (d *VersionDatabase) Fetch(version string) (VersionRow, bool) {
	stmt, err := d.db.Prepare("SELECT * FROM versions WHERE version = ?")
	check(err)

	rows, err := stmt.Query(version)
	check(err)

	var v VersionRow
	found := false
	for rows.Next() {
		err = rows.Scan(&v.UID, &v.Version, &v.Category, &v.Link, &v.Date, &v.Created)
		check(err)

		found = true
		break
	}
	rows.Close()

	return v, found
}

func (d *VersionDatabase) All() []VersionRow {
	versions := []VersionRow{}

	rows, err := d.db.Query("SELECT * FROM versions ORDER BY date")
	check(err)

	var v VersionRow
	for rows.Next() {
		err = rows.Scan(&v.UID, &v.Version, &v.Category, &v.Link, &v.Date, &v.Created)
		check(err)
		versions = append(versions, v)
	}
	rows.Close()
	return versions
}

func (d *VersionDatabase) Execute(sql string) {
	d.db.Exec(sql)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type DatabaseAccessor struct {
	db *VersionDatabase
}

func NewDBAccessor(db *VersionDatabase) *DatabaseAccessor {
	return &DatabaseAccessor{
		db: db,
	}
}

func (d *DatabaseAccessor) Run(initCh, insertCh chan VersionRow, quitCh chan int) {
	//
	running := true
	for running {
		select {
		case init := <-initCh:
			row := init
			d.db.Insert(&row)

		case insert := <-insertCh:
			row := insert
			_, found := d.db.Fetch(row.Version)
			if found {
				continue
			}

			d.db.Insert(&row)
			//msg := makeMessage(row.Version, row.Category, row.Link)
			// TODO
			//d.sender.send(msg)
			log.Printf("New version found : %s\n", row.Version)

		case <-quitCh:
			log.Println("stop db accessor")
			running = false
			return
		}
	}
}
