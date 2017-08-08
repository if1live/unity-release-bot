package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func renderErrorJSON(w http.ResponseWriter, err error, errcode int) {
	type Response struct {
		Error string `json:"error"`
	}
	resp := Response{
		Error: err.Error(),
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errcode)
	w.Write(data)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Port int       `json:"port"`
		Now  time.Time `json:"now"`
	}
	resp := &Response{
		Port: port,
		Now:  time.Now(),
	}
	renderJSON(w, resp)
}

func handleLatestVersion(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Now     time.Time `json:"now"`
		Version string    `json:"version"`
	}
	h := VersionHelper{}
	v := h.FromURI(UnityDownloadURL)
	resp := &Response{
		Now:     time.Now(),
		Version: v,
	}
	renderJSON(w, resp)
}

func handleRSSPatch(w http.ResponseWriter, r *http.Request) {
	feed := NewPatchRSS()
	renderJSON(w, feed.Rows())
}

func handleRSSBeta(w http.ResponseWriter, r *http.Request) {
	feed := NewBetaRSS()
	renderJSON(w, feed.Rows())
}

func handleView(w http.ResponseWriter, r *http.Request) {
	rows := g_svr.db.All()
	renderJSON(w, rows)
}

func handleSyncPatchRSS(w http.ResponseWriter, r *http.Request) {
	handleSyncCommonRSS(w, r, NewPatchRSS())
}

func handleSyncBetaRSS(w http.ResponseWriter, r *http.Request) {
	handleSyncCommonRSS(w, r, NewBetaRSS())
}

func handleSyncCommonRSS(w http.ResponseWriter, r *http.Request, feed *UnityFeed) {
	rows := feed.Rows()
	uids := make([]int64, len(rows))
	for i, r := range rows {
		uid := g_svr.db.Insert(&r)
		uids[i] = uid
	}
	renderJSON(w, uids)
}

func handleSyncLatest(w http.ResponseWriter, r *http.Request) {
	h := VersionHelper{}
	v := h.FromURI(UnityDownloadURL)
	link := h.makeStableReleaseNoteURL(v)
	row := VersionRow{
		Version:  v,
		Category: categoryStable,
		Date:     time.Now(),
		Link:     link,
	}
	uid := g_svr.db.Insert(&row)
	renderJSON(w, uid)
}

func handleDevIndex(w http.ResponseWriter, r *http.Request) {
	src := `
<ul>
<li><a href="/dev/latest-version">latest version</a></li>
<li><a href="/dev/rss-patch">rss patch</a></li>
<li><a href="/dev/rss-beta">rss beta</a></li>
<li><a href="/dev/view">view all</a></li>
<li><a href="/sync/rss-patch">sync rss patch</a></li>
<li><a href="/sync/rss-beta">sync rss beta</a></li>
<li><a href="/sync/latest">sync latest</a></li>
</ul>
`
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.Write([]byte(src))
}

type Server struct {
	port int
	db   *VersionDatabase
}

func NewServer(port int, db *VersionDatabase) *Server {
	return &Server{
		port: port,
		db:   db,
	}
}

var g_svr *Server

func (s *Server) Main() {
	g_svr = s

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/dev", handleDevIndex)

	http.HandleFunc("/dev/latest-version", handleLatestVersion)
	http.HandleFunc("/dev/rss-patch", handleRSSPatch)
	http.HandleFunc("/dev/rss-beta", handleRSSBeta)
	http.HandleFunc("/dev/view", handleView)

	http.HandleFunc("/sync/rss-patch", handleSyncPatchRSS)
	http.HandleFunc("/sync/rss-beta", handleSyncBetaRSS)
	http.HandleFunc("/sync/latest", handleSyncLatest)

	fmt.Printf("Run server : port=%d\n", s.port)
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)

}
