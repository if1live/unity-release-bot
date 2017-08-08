package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVersion(t *testing.T) {
	cases := []struct {
		filepath string
		version  string
	}{
		// url
		{
			"http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe",
			"5.5.0f3",
		},
		// 5.x
		{
			"UnityDownloadAssistant-5.5.0f3.exe",
			"5.5.0f3",
		},
		// postfix
		{
			"UnityDaydreamDownloadAssistant-5.4.2f2-GVR13.exe",
			"5.4.2f2-GVR13",
		},
		// different prefix
		{
			"UnityStandardAssetsSetup-5.5.1f1.exe",
			"5.5.1f1",
		},
		// unity 2017
		{
			"UnityDownloadAssistant-2017.1.0p2.exe",
			"2017.1.0p2",
		},
		// format extension
		{
			"Unity-12.34.56f78.exe",
			"12.34.56f78",
		},
	}

	for _, c := range cases {
		v := extractVersion(c.filepath)
		if v != c.version {
			t.Error("Expected ", c.version, ", got ", v)
		}
	}
}

func TestGetLatestVersion(t *testing.T) {
	cases := []struct {
		file    string
		version string
	}{
		{"testdata/download-stable-5.5.1f1.html", "5.5.1f1"},
		{"testdata/download-stable-2017.1.0f3.html", "2017.1.0f3"},
	}

	for _, c := range cases {
		src, _ := ioutil.ReadFile(c.file)
		f := &FakeHTTPFetcher{
			idx:     0,
			sources: []string{string(src)},
		}
		version := getLatestVersion(f)
		assert.Equal(t, c.version, version)
	}
}

func TestMakeStableReleaseNoteURL(t *testing.T) {
	cases := []struct {
		version string
		link    string
	}{
		{"5.5.0f3", "https://unity3d.com/kr/unity/whats-new/unity-5.5.0"},
		{"5.5.0", "https://unity3d.com/kr/unity/whats-new/unity-5.5.0"},
		{"2017.1.0f3", "https://unity3d.com/kr/unity/whats-new/unity-2017.1.0"},
		{"2017.1.0", "https://unity3d.com/kr/unity/whats-new/unity-2017.1.0"},
		{"invalid", ""},
	}
	for _, c := range cases {
		v := makeStableReleaseNoteURL(c.version)
		if v != c.link {
			t.Error("Expected ", c.link, ", got ", v)
		}
	}
}
