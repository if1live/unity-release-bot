package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
)

const UnityDownloadURL = "https://unity3d.com/kr/get-unity/download?ref=professional"

type VersionHelper struct {
}

func (h *VersionHelper) extractVersion(fp string) string {
	// http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe
	filename := path.Base(fp)
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	idx := strings.Index(filename, "-")
	if idx < 0 {
		return ""
	}
	found := filename[idx+1:]
	return found
}

func (h *VersionHelper) makeStableReleaseNoteURL(v string) string {
	// 5.5.0f3 -> https://unity3d.com/kr/unity/whats-new/unity-5.5.0
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	m := re.FindString(v)
	if len(m) == 0 {
		return ""
	}
	return "https://unity3d.com/kr/unity/whats-new/unity-" + m
}

func (h *VersionHelper) fromHTML(src string) string {
	// http://netstorage.unity3d.com/unity/88d00a7498cd/WindowsStandardAssetsInstaller/UnityStandardAssetsSetup-5.5.1f1.exe
	// 윈도우 인스톨러를 기준으로 잡아내면 될듯?
	re := regexp.MustCompile(`"http://netstorage.unity3d.com/unity/\w+/WindowsStandardAssetsInstaller/(.+\.exe)"`)
	links := re.FindAllStringSubmatch(src, -1)
	link := links[0][1]
	version := h.extractVersion(link)
	return version
}

func (h *VersionHelper) FromURI(uri string) string {
	isHTTP := strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://")
	if isHTTP {
		resp, err := http.Get(uri)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()

		src := buf.String()
		return h.fromHTML(src)
	}

	src, _ := ioutil.ReadFile(uri)
	return h.fromHTML(string(src))
}
