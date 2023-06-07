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
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/google/go-github/v39/github"
	github_go "github.com/labring/gh-rebot/pkg/github-go"
	"os"
	"strconv"
	"strings"
)

func GetEnvFromAction(key string) (string, error) {
	allKey := strings.ToLower("sealos_" + key)
	val, _ := os.LookupEnv(allKey)
	if val == "" {
		return "", fmt.Errorf("not found %s", allKey)
	}
	return val, nil
}

func getPrNumber() (int, error) {
	ref := os.Getenv("GITHUB_REF") // 获取环境变量GITHUB_REF
	if ref == "" {
		return 0, fmt.Errorf("not found GITHUB_REF")
	}
	split := strings.Split(ref, "/")
	if len(split) != 4 {
		return 0, fmt.Errorf("GITHUB_REF format error")
	}
	// 获取并转换PR编号
	prNumberStr := split[2]
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		return 0, fmt.Errorf("GITHUB_REF format Atoi error")
	}
	return prNumber, nil
}

func getRepo() (string, string, error) {
	repo := os.Getenv("GITHUB_REPOSITORY") // 获取环境变量GITHUB_REF
	if repo == "" {
		return "", "", fmt.Errorf("not found GITHUB_REPOSITORY")
	}
	split := strings.Split(repo, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("GITHUB_REPOSITORY format error")
	}
	return split[0], split[1], nil
}

// PRComment is a action to comment on PR
func PRComment() error {
	fileName, err := GetEnvFromAction("filename")
	if err != nil {
		return err
	}
	replaceTag, err := GetEnvFromAction("replace_tag")
	if err != nil {
		return err
	}

	prNumber, err := getPrNumber()
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	owner, repo, err := getRepo()
	if err != nil {
		return err
	}
	logger.Info("repo:%s, filename: %s, replaceTag: %s, prNumber: %d", repo, fileName, replaceTag, prNumber)
	ctx := context.Background()
	client := github_go.GithubClient(ctx)
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return fmt.Errorf("Issues.ListComments returned error: %v", err)
	}
	hiddenReplace := fmt.Sprintf("<!-- %s -->", replaceTag)
	createComment := func() {
		comment := &github.IssueComment{Body: github.String(string(fileContent))}
		client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	}
	if hiddenReplace == "" {
		//add
		createComment()
		return nil
	}
	for _, comment := range comments {
		if comment.Body != nil && comment.ID != nil {
			if *comment.Body == string(fileContent) {
				logger.Info("The comment %d has been already added to the pull request. Skipping...", *comment.ID)
				return nil
			} else if hiddenReplace != "" && strings.LastIndex(*comment.Body, hiddenReplace) != -1 {
				client.Issues.EditComment(ctx, owner, repo, *comment.ID, &github.IssueComment{Body: github.String(string(fileContent))})
				return nil
			}
		}
	}
	createComment()
	return nil
}
