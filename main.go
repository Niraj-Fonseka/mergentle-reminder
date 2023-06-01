package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

func main() {
	// Load configuration
	config, err := loadConfig(&OsEnv{})
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	glClient, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(config.GitLab.URL))
	if err != nil {
		fmt.Printf("Error creating GitLab client: %v\n", err)
		os.Exit(1)
	}

	gitlabClient := &gitLabClient{client: glClient}
	slackClient := &slackClient{}

	notify := &notify{gitlab: gitlabClient, slack: slackClient, config: config}

	// job := func() {
	// 	remind.notify()
	// }
	//scheduler.Every().Day().Run(job)
	notify.notify()

	select {}
}

func readConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
