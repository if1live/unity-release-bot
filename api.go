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

func mainServer(port int) {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/dev/latest-version", handleLatestVersion)
	http.HandleFunc("/dev/rss-patch", handleRSSPatch)
	http.HandleFunc("/dev/rss-beta", handleRSSBeta)

	fmt.Printf("Run server : port=%d\n", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
