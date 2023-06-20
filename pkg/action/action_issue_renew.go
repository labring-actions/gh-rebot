/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"context"
	"github.com/cuisongliu/logger"
	"github.com/google/go-github/v39/github"
	github_go "github.com/labring/gh-rebot/pkg/github-go"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

// IssueRenew is new a issue
func IssueRenew() error {
	issueTitle, err := GetEnvFromAction("issue_title")
	if err != nil {
		return err
	}
	label, _ := GetEnvFromAction("issue_label")
	if err != nil {
		return err
	}
	body, _ := GetEnvFromAction("issue_body")
	bodyfile, _ := GetEnvFromAction("issue_bodyfile")
	if bodyfile != "" {
		bodyBytes, _ := os.ReadFile(bodyfile)
		body = string(bodyBytes)
	}
	issueType, err := GetEnvFromAction("issue_type")
	if err != nil {
		return err
	}
	switch issueType {
	case "day":
		issueTitle = issueTitle + " " + utils.FormatDay(time.Now())
	case "week":
		start, end := utils.FormatWeek(time.Now())
		issueTitle = issueTitle + " " + start + " to " + end
	default:
		issueTitle = issueTitle + " " + utils.FormatDay(time.Now())
	}

	owner, repo, err := getRepo()
	if err != nil {
		return err
	}
	ctx := context.Background()
	client := github_go.GithubClient(ctx)

	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		Creator: types.GlobalsBotConfig.Bot.Username,
	})
	logger.Info("repo:%s, issueTitle: %s, owner: %s, create: %s", repo, issueTitle, owner, types.GlobalsBotConfig.Bot.Username)
	if err != nil {
		return err
	}
	hasIssue := false
	issueId := ""
	defer func() {
		writeGithubEnv("SEALOS_ISSUE_ID", issueId)
		logger.Info("add env SEALOS_ISSUE_ID: %s", issueId)
	}()
	issueOldTitle, _ := GetEnvFromAction("issue_title")
	for _, issue := range issues {
		if strings.HasPrefix(issue.GetTitle(), issueOldTitle) {
			logger.Debug("issue: %s, state: %s, id: %d", issue.GetTitle(), issue.GetState(), issue.GetID())
			if issue.GetTitle() == issueTitle && issue.GetState() != "closed" {
				logger.Info("issue already exist, issue: %s", issue.GetTitle())
				hasIssue = true
				issueId = strconv.Itoa(int(issue.GetID()))
				return nil
			} else {
				state := "closed"
				issueRequest := &github.IssueRequest{
					State: &state,
				}
				_, _, _ = client.Issues.Edit(ctx, owner, repo, issue.GetNumber(), issueRequest)
				logger.Info("close issue: %s", issue.GetTitle())
			}
		}
	}
	logger.Warn("issue not exist, issue: %s", issueTitle)
	if !hasIssue {
		issueRequest := &github.IssueRequest{
			Title: &issueTitle,
			Body:  &body,
			Labels: &[]string{
				label,
			},
		}
		issue, _, _ := client.Issues.Create(ctx, owner, repo, issueRequest)
		issueId = strconv.Itoa(int(issue.GetID()))
		logger.Info("create issue: %s, id: %d", issueTitle, issue.GetID())
	}

	return nil
}
