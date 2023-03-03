package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func getGithubAccessToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configPath := home + "/.config/gh/hosts.yml"

	yfile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	type Host struct {
		OauthToken string `yaml:"oauth_token"`
	}

	data := make(map[interface{}]Host)
	err = yaml.Unmarshal(yfile, &data)

	token := data["github.com"].OauthToken

	return token
}
