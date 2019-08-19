package main

import (
	"database/sql"
	"log"
	"strconv"
)

type IssueRepository struct {
	Db *sql.DB
}

func (repository *IssueRepository) getActualIssuesByStatusName(statusName *string) []Issue {
	rows, err := repository.Db.Query("select id,jira_status_name,jira_key,jira_assignee_login,entity_id,jira_labels from palantir.jira_issue " +
		"WHERE jira_status_name = '" + *statusName +
		"'and project_key='EAISTPK' and deleted_date_time is null " +
		"and jira_assignee_login is not null " +
		"and jira_labels is not null " +
		"order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var issues []Issue
	for rows.Next() {
		issue := Issue{}
		err = rows.Scan(&issue.id, &issue.statusName, &issue.key, &issue.assigneeLogin, &issue.entityId, &issue.jiraLabels, &issue.jiraTesterLogin)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, issue)
	}

	log.Println("received len: " + strconv.Itoa(len(issues)) + ", status: " + *statusName)
	return issues
}
