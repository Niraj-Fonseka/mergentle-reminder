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

// func formatMergeRequestsSummary(mrs []*MergeRequestWithApprovals) string {
// 	var summary string
// 	for _, mr := range mrs {
// 		approvedBy := strings.Join(mr.ApprovedBy, ", ")
// 		if approvedBy == "" {
// 			approvedBy = "None"
// 		}

// 		createdAtStr := mr.MergeRequest.CreatedAt.Format("2 January 2006, 15:04 MST")

// 		summary += fmt.Sprintf(
// 			":arrow_forward: <%s|%s>\n*Author:* %s\n*Created at:* %s\n*Approved by:* %s\n\n",
// 			mr.MergeRequest.WebURL, mr.MergeRequest.Title, mr.MergeRequest.Author.Name, createdAtStr, approvedBy,
// 		)
// 	}

// 	return summary
// }
