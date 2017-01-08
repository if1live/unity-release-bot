package main

import (
	"testing"
	"time"

	"io/ioutil"

	"os"
	"path"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "unity-release-bot")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filepath := path.Join(tempDir, "test.sqlite3")
	db := NewDB(filepath)
	b, _ := ioutil.ReadFile("schema.sql")
	sql := string(b)
	db.execute(sql)

	version := "5.5.0f3"
	link := "http://google.com"
	category := "alpha"
	db.insert(version, category, link, time.Now())

	r, found := db.fetch(version)
	assert.Equal(t, found, true)
	assert.Equal(t, r.version, version)
	assert.Equal(t, r.link, link)
	assert.Equal(t, r.category, category)

	_, found = db.fetch("5.5.0f1")
	assert.Equal(t, found, false)

	db.close()
}
