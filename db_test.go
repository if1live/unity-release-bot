package main

import (
	"testing"
	"time"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	db := NewDB(":memory:")

	b, _ := ioutil.ReadFile("schema.sql")
	sql := string(b)
	db.Execute(sql)

	ins := VersionRow{
		Version:  "5.5.0f3",
		Link:     "http://google.com",
		Category: "alpha",
		Date:     time.Now(),
	}
	db.Insert(&ins)

	found, ok := db.Fetch(ins.Version)
	assert.Equal(t, ok, true)
	assert.Equal(t, found.Version, ins.Version)
	assert.Equal(t, found.Link, ins.Link)
	assert.Equal(t, found.Category, ins.Category)

	_, ok = db.Fetch("5.5.0f1")
	assert.Equal(t, ok, false)

	db.Close()
}
