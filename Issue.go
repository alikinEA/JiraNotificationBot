package main

type Issue struct {
	id              int
	statusName      string
	key             string
	assigneeLogin   string
	entityId        int
	jiraLabels      string
	jiraTesterLogin string
}
