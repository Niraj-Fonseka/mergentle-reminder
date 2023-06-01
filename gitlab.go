package main

import "github.com/xanzy/go-gitlab"

//go:generate mockery --name GitLabClient
type GitLabClient interface {
	ListGroupProjects(groupID int, options *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error)
	ListSubGroups(groupID int, opt *gitlab.ListSubGroupsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Group, *gitlab.Response, error)
	ListProjectMergeRequests(projectID int, options *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	GetMergeRequestApprovalsConfiguration(projectID int, mergeRequestID int) (*gitlab.MergeRequestApprovals, *gitlab.Response, error)
}

type MergeRequestWithApprovals struct {
	MergeRequest *gitlab.MergeRequest
	ApprovedBy   []string
}

type gitLabClient struct {
	client *gitlab.Client
}

func (c *gitLabClient) ListGroupProjects(groupID int, options *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	return c.client.Groups.ListGroupProjects(groupID, options)
}

func (c *gitLabClient) ListSubGroups(groupID int, opt *gitlab.ListSubGroupsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Group, *gitlab.Response, error) {
	return c.client.Groups.ListSubGroups(groupID, opt, options...)
}

func (c *gitLabClient) ListProjectMergeRequests(projectID int, options *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return c.client.MergeRequests.ListProjectMergeRequests(projectID, options)
}

func (c *gitLabClient) GetMergeRequestApprovalsConfiguration(projectID int, mergeRequestID int) (*gitlab.MergeRequestApprovals, *gitlab.Response, error) {
	return c.client.MergeRequestApprovals.GetConfiguration(projectID, mergeRequestID)
}

func (c *gitLabClient) GetProject(projectID int, opt *gitlab.GetProjectOptions) (*gitlab.Project, *gitlab.Response, error) {
	return c.client.Projects.GetProject(projectID, opt)
}
