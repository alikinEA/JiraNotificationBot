package main

import "database/sql"

type Issue struct {
	id              int
	statusName      string
	key             string
	assigneeLogin   string
	entityId        int
	jiraLabels      string
	jiraTesterLogin sql.NullString
}
