jira telegram notification bot
build app:
    go build
run app:
    ./JiraNotificationBot botId lognProxy passProxy ip:port idChat postgres://postgres:postgres@ip/schemaname jiraLogin1:@NickTelegram1,jiraLogin2:@NickTelegram2