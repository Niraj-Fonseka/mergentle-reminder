package main

import "github.com/slack-go/slack"

//go:generate mockery --name SlackClient
type SlackClient interface {
	PostWebhook(payload *slack.WebhookMessage) error
}

type slackClient struct {
}

func (c *slackClient) PostWebhook(webhook string, payload *slack.WebhookMessage) error {
	return slack.PostWebhook(webhook, payload)
}
