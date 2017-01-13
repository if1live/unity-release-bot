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

func watchRSS(ctx *Context, rssurl string, category string, delay time.Duration) {
	log.Printf("RSS [%s] : %s\n", category, rssurl)
	feed, err := rss.Fetch(rssurl)
	if err != nil {
		panic(err)
	}

	initialized := false
	for {
		for _, item := range feed.Items {
			version := strings.Replace(item.Title, "Patch", "", -1)
			version = strings.Trim(version, " ")
			row := VersionRow{
				version:  version,
				category: category,
				date:     item.Date,
				link:     item.Link,
			}

			if !initialized {
				ctx.initCh <- row
				initialized = true
			} else {
				ctx.insertCh <- row
			}
		}

		time.Sleep(delay)
		feed.Update()
		log.Printf("RSS[%s] : %s\n", category, rssurl)
	}
}

func watchLatestVersion(ctx *Context, category string, delay time.Duration) {
	initialized := false
	for {
		f := &RealHTTPFetcher{}
		version := getLatestVersion(f)
		log.Printf("Latest Version [%s] : %s\n", category, version)

		link := makeStableReleaseNoteURL(version)
		row := VersionRow{
			version:  version,
			category: category,
			date:     time.Now(),
			link:     link,
		}

		if !initialized {
			ctx.initCh <- row
			initialized = true
		} else {
			ctx.insertCh <- row
		}

		time.Sleep(delay)
	}
}

var logfilename string

func init() {
	flag.StringVar(&logfilename, "log", "", "log filename")
}

type Context struct {
	config   *Config
	accessor *DatabaseAccessor

	initCh   chan VersionRow
	insertCh chan VersionRow
	quitCh   chan int
}

var ctx Context

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
	defer db.close()

	c := NewConfig()
	c = nil

	ctx = Context{
		config:   c,
		accessor: NewDBAccessor(db, NewSender(c)),

		initCh:   make(chan VersionRow, 10),
		insertCh: make(chan VersionRow, 10),
		quitCh:   make(chan int),
	}

	interval := 15 * time.Minute
	go watchRSS(&ctx, rssPatch, categoryPatch, interval)
	go watchRSS(&ctx, rssBeta, categoryBeta, interval)
	go watchLatestVersion(&ctx, categoryStable, interval)

	ctx.accessor.Run(ctx.initCh, ctx.insertCh, ctx.quitCh)
}
