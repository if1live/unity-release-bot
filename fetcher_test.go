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
		{
			"http://netstorage.unity3d.com/unity/38b4efef76f0/UnityDownloadAssistant-5.5.0f3.exe",
			"5.5.0f3",
		},
		{"UnityDownloadAssistant-5.5.0f3.exe", "5.5.0f3"},
		{"Unity-12.34.56f78.exe", "12.34.56f78"},
		{"UnityDaydreamDownloadAssistant-5.4.2f2-GVR13.exe", "5.4.2f2-GVR13"},
		{"UnityStandardAssetsSetup-5.5.1f1.exe", "5.5.1f1"},
		{"CacheServer-5.5.1f1.zip", "5.5.1f1"},
	}

	for _, c := range cases {
		v := extractVersion(c.filepath)
		if v != c.version {
			t.Error("Expected ", c.version, ", got ", v)
		}
	}
}

func TestGetLatestVersion(t *testing.T) {
	src, _ := ioutil.ReadFile("testdata/download-stable-5.5.1f1.html")
	f := &FakeHTTPFetcher{
		idx:     0,
		sources: []string{string(src)},
	}
	version := getLatestVersion(f)
	assert.Equal(t, "5.5.1f1", version)
}

func TestMakeStableReleaseNoteURL(t *testing.T) {
	cases := []struct {
		version string
		link    string
	}{
		{"5.5.0f3", "https://unity3d.com/kr/unity/whats-new/unity-5.5.0"},
		{"5.5.0", "https://unity3d.com/kr/unity/whats-new/unity-5.5.0"},
		{"invalid", ""},
	}
	for _, c := range cases {
		v := makeStableReleaseNoteURL(c.version)
		if v != c.link {
			t.Error("Expected ", c.link, ", got ", v)
		}
	}
}
