package main

import (
	"flag"
	"log"
	"os"
	"path"

	"time"

	"strings"

	"github.com/SlyMarbo/rss"
)

func insertVersion(row VersionRow, s Sender, db VersionDatabase) {
	db.insert(row.version, row.category, row.link, row.date)
}
func checkVersion(row VersionRow, s Sender, db VersionDatabase) {
	_, found := db.fetch(row.version)
	if found {
		return
	}

	db.insert(row.version, row.category, row.link, row.date)
	msg := makeMessage(row.version, row.category, row.link)
	s.send(msg)
	log.Printf("New version found : %s\n", row.version)
}

func watchRSS(c *Config, db VersionDatabase, rssurl string, category string, delay time.Duration) {
	feed, err := rss.Fetch(rssurl)
	if err != nil {
		panic(err)
	}

	sender := NewSender(c)

	for i := 0; ; {
		for _, item := range feed.Items {
			version := strings.Replace(item.Title, "Patch", "", -1)
			version = strings.Trim(version, " ")
			row := VersionRow{
				version:  version,
				category: category,
				date:     item.Date,
				link:     item.Link,
			}

			if i == 0 {
				insertVersion(row, sender, db)
			} else {
				checkVersion(row, sender, db)
			}
		}

		time.Sleep(delay)
		feed.Update()
	}
}

func watchLatestVersion(c *Config, db VersionDatabase, category string, delay time.Duration) {
	sender := NewSender(c)

	for i := 0; ; {
		f := &RealHTTPFetcher{}
		version := getLatestVersion(f)
		link := makeStableReleaseNoteURL(version)
		row := VersionRow{
			version:  version,
			category: category,
			date:     time.Now(),
			link:     link,
		}

		if i == 0 {
			insertVersion(row, sender, db)
		} else {
			checkVersion(row, sender, db)
		}

		time.Sleep(delay)
	}
}

var logfilename string

func init() {
	flag.StringVar(&logfilename, "log", "bot.log", "log filename")
}

func main() {
	flag.Parse()

	// initialize logger
	// http: //stackoverflow.com/questions/19965795/go-golang-write-log-to-file
	// logger 초기화를 별도 함수에서 할 경우 defer 로 파일이 닫혀서 로그작성이 안된다
	// 그래서 그냥 메인함수에서 처리
	if logfilename != "" {
		filepath := path.Join(getExecutablePath(), logfilename)
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	db := NewDB("db.sqlite3")

	c := NewConfig()
	c = nil

	interval := 15 * time.Minute
	go watchRSS(c, db, rssPatch, categoryPatch, interval)
	go watchRSS(c, db, rssBeta, categoryBeta, interval)
	go watchLatestVersion(c, db, categoryStable, interval)

	for {
		delay := 1 * time.Minute
		time.Sleep(delay)
	}
}
