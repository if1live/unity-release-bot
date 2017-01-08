package main

import (
	"bytes"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
)

func extractVersion(filepath string) string {
	// http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe
	filename := path.Base(filepath)
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	idx := strings.Index(filename, "-")
	if idx < 0 {
		return ""
	}
	found := filename[idx+1:]
	return found
}

func makeStableReleaseNoteURL(version string) string {
	// 5.5.0f3 -> https://unity3d.com/kr/unity/whats-new/unity-5.5.0
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	m := re.FindString(version)
	if len(m) == 0 {
		return ""
	}
	return "https://unity3d.com/kr/unity/whats-new/unity-" + m
}

func getLatestVersion(fetcher HTTPFetcher) string {
	initialURL := "https://store.unity.com/download?ref=personal"
	initialHTML := fetcher.Fetch(initialURL)
	// <a href="https://store.unity.com/download/thank-you?thank-you=personal&amp;os=win&amp;nid=178" class="download-btn bg-gr os-windows hide">
	initialRe := regexp.MustCompile(`"https://store.unity.com/download/thank-you.+" `)
	results := initialRe.FindAllString(initialHTML, -1)
	nexturl := results[0]
	nexturl = strings.Replace(nexturl, `"`, "", -1)
	nexturl = strings.Trim(nexturl, " ")
	nexturl = strings.Replace(nexturl, "&amp;", "&", -1)

	html := fetcher.Fetch(nexturl)
	// downloadUrl = 'http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe';
	re := regexp.MustCompile(`downloadUrl = '(.+)';`)
	links := re.FindAllStringSubmatch(html, -1)
	link := links[0][1]
	// http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe

	version := extractVersion(link)
	return version
}

type HTTPFetcher interface {
	Fetch(rawurl string) string
}

type RealHTTPFetcher struct {
}

func (f *RealHTTPFetcher) Fetch(rawurl string) string {
	resp, err := http.Get(rawurl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	resp.Body.Close()

	return buf.String()
}

type FakeHTTPFetcher struct {
	sources []string
	idx     int
}

func (f *FakeHTTPFetcher) Fetch(rawurl string) string {
	source := f.sources[f.idx]
	f.idx = (f.idx + 1) % len(f.sources)
	return source
}
