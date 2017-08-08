package main

import (
	"strings"

	"github.com/SlyMarbo/rss"
)

const (
	categoryStable = "stable"
	categoryPatch  = "patch"
	categoryBeta   = "beta"

	rssPatch = "https://unity3d.com/unity/qa/patch-releases/latest.xml"
	rssBeta  = "https://unity3d.com/unity/beta/latest.xml"
)

type UnityFeed struct {
	category string
	rssurl   string
	feed     *rss.Feed
}

func NewFeed(category, rssurl string) *UnityFeed {
	feed, err := rss.Fetch(rssurl)
	if err != nil {
		panic(err)
	}

	return &UnityFeed{
		category: category,
		rssurl:   rssurl,
		feed:     feed,
	}
}

func NewPatchRSS() *UnityFeed {
	return NewFeed(categoryPatch, rssPatch)
}

func NewBetaRSS() *UnityFeed {
	return NewFeed(categoryBeta, rssBeta)
}

func (f *UnityFeed) Update() {
	f.feed.Update()
}

func (f *UnityFeed) Rows() []VersionRow {
	rows := []VersionRow{}
	for _, item := range f.feed.Items {
		version := strings.Replace(item.Title, "Patch", "", -1)
		version = strings.Trim(version, " ")
		row := VersionRow{
			Version:  version,
			Category: f.category,
			Date:     item.Date,
			Link:     item.Link,
		}
		rows = append(rows, row)
	}
	return rows
}
