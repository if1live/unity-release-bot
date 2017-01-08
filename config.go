package main

import (
	"io/ioutil"
	"path"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kardianos/osext"
	yaml "gopkg.in/yaml.v2"
)

const (
	categoryStable = "stable"
	categoryPatch  = "patch"
	categoryBeta   = "beta"

	rssPatch = "https://unity3d.com/unity/qa/patch-releases/latest.xml"
	rssBeta  = "https://unity3d.com/unity/beta/latest.xml"
)

type Config struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
}

func getExecutablePath() string {
	path, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	return path
}

func NewConfig() *Config {
	filename := "config.yaml"
	filepath := path.Join(getExecutablePath(), filename)

	c := Config{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &c
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return &c
	}
	return &c
}

func (c *Config) createAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api := anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)
	return api
}
