package main

import (
	"fmt"
	"strconv"
)

type NotificationService struct {
	repository          *IssueRepository
	botApiService       *BotApiService
	statusName          string
	telegramNickNameMap map[string]string
	currentIssues       *[]Issue
}

func (service *NotificationService) checkUpdateIssues() {
	var issues = service.repository.getActualIssuesByStatusName(&service.statusName)
	if len(*service.currentIssues) != 0 {
		var newIssues []Issue
		for _, value1 := range issues {
			if !containsIssue(service.currentIssues, &value1) {
				newIssues = append(newIssues, value1)
			}
		}
		if len(newIssues) > 0 {
			for _, value := range newIssues {
				assignee := value.assigneeLogin
				if val, ok := service.telegramNickNameMap[value.assigneeLogin]; ok {
					assignee = val
				}

				message := "таск: " + JiraUrl + value.key +
					", переведен в статус: " + value.statusName +
					", исполнителем назначен: " + assignee +
					", labels: " + value.jiraLabels

				fmt.Println("Message to chat: " + strconv.Itoa(service.botApiService.chatId) + ", " + message)
				service.botApiService.sendMessageToChat(message)
			}
		}
		newIssues = nil
	}
	service.currentIssues = &issues
}

func containsIssue(issues *[]Issue, issue *Issue) bool {
	for _, value := range *issues {
		if issue.id == value.id {
			return true
		}
		if issue.entityId == value.entityId {
			if value.assigneeLogin == issue.assigneeLogin {
				return true
			}
		}
	}
	return false
}
