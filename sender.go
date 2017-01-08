package main

import (
	"fmt"
	"net/url"
	"strings"
)

type Sender interface {
	send(text string)
}

func makeMessage(version, category, link string) string {
	lines := []string{
		fmt.Sprintf("New Unity3D Released!"),
		fmt.Sprintf("Version = %s", version),
		fmt.Sprintf("Type = %s", category),
		fmt.Sprintf("Detail %s", link),
	}
	return strings.Join(lines, "\n")
}

func NewSender(c *Config) Sender {
	if c == nil {
		return &FakeSender{
			texts: []string{},
		}
	}
	return &TwitterSender{
		config: c,
	}
}

type TwitterSender struct {
	config *Config
}

func (s *TwitterSender) send(text string) {
	api := s.config.createAPI()
	v := url.Values{}
	api.PostTweet(text, v)
}

type FakeSender struct {
	texts []string
	last  string
}

func (s *FakeSender) send(text string) {
	s.texts = append(s.texts, text)
	s.last = text

	fmt.Printf("Send : %s\n", text)
}
