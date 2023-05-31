package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

type notify struct {
	gitlab *gitLabClient
	slack  *slackClient
	config *Config
}

func (r *notify) notify() {
	fmt.Println("running notify function")

	//notify groups
	for _, group := range r.config.Groups {
		var groupIDs []int
		groupIDs = append(groupIDs, group.ID)

		// Add subgroups to the groups list.
		subgroupIDs, err := r.fetchSubGroups(group.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		groupIDs = append(groupIDs, subgroupIDs...)

		projectIds, err := r.fetchProjectsFromGroups(groupIDs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		mrs, err := r.fetchOpenedMergeRequests(projectIds...)
		if err != nil {
			fmt.Println(err)
			continue
		}

		summary := r.formatMergeRequestsSummary(mrs)

		r.slackNotify(group.SlackWebhook, summary)
	}

	for _, project := range r.config.Projects {
		mrs, err := r.fetchOpenedMergeRequests(project.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		summary := r.formatMergeRequestsSummary(mrs)

		r.slackNotify(project.SlackWebhook, summary)
	}
}

func (r *notify) slackNotify(slackWebhook, summary string) error {
	return r.sendSlackMessage(slackWebhook, summary)
}

func (r *notify) sendSlackMessage(webhook, message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	return r.slack.PostWebhook(webhook, &msg)
}

func (r *notify) formatMergeRequestsSummary(mrs []*MergeRequestWithApprovals) string {
	var summary string

	//mrs to projects
	mrsProjects := make(map[int][]*MergeRequestWithApprovals, 0)

	for _, mr := range mrs {
		_, ok := mrsProjects[mr.MergeRequest.ProjectID]
		if ok {
			mrsProjects[mr.MergeRequest.ProjectID] = append(mrsProjects[mr.MergeRequest.ProjectID], mr)
		} else {
			mrsProjects[mr.MergeRequest.ProjectID] = []*MergeRequestWithApprovals{mr}
		}
	}

	for gitlabProjectID, mrs := range mrsProjects {
		projectName, err := r.gitlab.GetProject(gitlabProjectID)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(mrs) == 0 {
			return fmt.Sprintf(":tada: There are no open merge requests for %s ! :tada:", projectName)
		}

		for _, mr := range mrs {
			approvedBy := strings.Join(mr.ApprovedBy, ", ")
			if approvedBy == "" {
				approvedBy = "None"
			}

			createdAtStr := mr.MergeRequest.CreatedAt.Format("2 January 2006, 15:04 MST")

			summary += fmt.Sprintf(
				":arrow_forward: %s <%s|%s>\n*Author:* %s\n*Created at:* %s\n*Approved by:* %s\n\n",
				projectName, mr.MergeRequest.WebURL, mr.MergeRequest.Title, mr.MergeRequest.Author.Name, createdAtStr, approvedBy,
			)
		}
	}

	return summary
}

func (r *notify) fetchOpenedMergeRequests(projectIDs ...int) ([]*MergeRequestWithApprovals, error) {

	var allMRs []*MergeRequestWithApprovals

	for _, projectID := range projectIDs {
		options := &gitlab.ListProjectMergeRequestsOptions{
			State:   gitlab.String("opened"),
			OrderBy: gitlab.String("updated_at"),
			Sort:    gitlab.String("desc"),
			WIP:     gitlab.String("no"),
			ListOptions: gitlab.ListOptions{
				PerPage: 50,
				Page:    1,
			},
		}

		for {
			mrs, resp, err := r.gitlab.ListProjectMergeRequests(projectID, options)
			if err != nil {
				return nil, err
			}

			for _, mr := range mrs {
				approvals, _, err := r.gitlab.GetMergeRequestApprovalsConfiguration(projectID, mr.IID)
				if err != nil {
					return nil, err
				}

				approvedBy := make([]string, len(approvals.ApprovedBy))
				for i, approver := range approvals.ApprovedBy {
					approvedBy[i] = approver.User.Name
				}

				allMRs = append(allMRs, &MergeRequestWithApprovals{
					MergeRequest: mr,
					ApprovedBy:   approvedBy,
				})
			}

			if resp.CurrentPage >= resp.TotalPages {
				break
			}

			options.Page = resp.NextPage
		}
	}

	return allMRs, nil
}

func (r *notify) fetchProjectsFromGroups(groupIDs []int) ([]int, error) {
	var projectIDs []int
	for _, groupID := range groupIDs {
		options := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 50,
				Page:    1,
			},
		}

		for {
			projects, resp, err := r.gitlab.ListGroupProjects(groupID, options)
			if err != nil {
				return nil, err
			}

			for _, project := range projects {
				projectIDs = append(projectIDs, project.ID)
			}

			if resp.CurrentPage >= resp.TotalPages {
				break
			}

			options.Page = resp.NextPage
		}
	}

	return projectIDs, nil
}

func (r *notify) fetchSubGroups(groupID int) ([]int, error) {
	var groupIDs []int

	options := &gitlab.ListSubGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	for {
		groups, resp, err := r.gitlab.ListSubGroups(groupID, options)
		if err != nil {
			return nil, err
		}

		for _, group := range groups {
			groupIDs = append(groupIDs, group.ID)
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		options.Page = resp.NextPage
	}

	return groupIDs, nil
}
