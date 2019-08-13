package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	reviewStatus string = "РЕВЬЮ"
	toDoStatus string = "TO DO"
	testingStatus string = "ТЕСТИРОВАНИЕ"
	codeReviewStatus string = "КОД-РЕВЬЮ"
)

func initDB(dbConnectSettings string) *sql.DB {
	connStr := dbConnectSettings + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {

	fmt.Println("start up args string: ", os.Args[1:])
	botId := os.Args[1]
	proxyUser := os.Args[2]
	proxyPass := os.Args[3]
	proxyIp := os.Args[4]
	chatIdStr := os.Args[5]
	dbConnectSettings := os.Args[6]
	nickMapStr := os.Args[7]
	chatId, _ := strconv.Atoi(chatIdStr)

	fmt.Println("start up args botId:" + botId)
	fmt.Println("start up args proxyUser:" + proxyUser)
	fmt.Println("start up args proxyPass:" + proxyPass)
	fmt.Println("start up args proxyIp:" + proxyIp)
	fmt.Println("start up args chatIdStr:" + chatIdStr)
	fmt.Println("start up args nickMapStr:" + nickMapStr)

	nickSlice := strings.Split(nickMapStr, ",")

	var telegramNickNameMap = make(map[string]string)
	for _, value := range nickSlice {
		nickSliceKeyValue := strings.Split(value, ":")
		login := nickSliceKeyValue[0]
		telegramNick := nickSliceKeyValue[1]
		telegramNickNameMap[login] = telegramNick
	}

	botApiService := BotApiService{botId: botId, proxyUser: proxyUser, proxyPass: proxyPass, proxyIp: proxyIp, chatId: chatId}
	fmt.Println(botApiService.getMe())

	var db = initDB(dbConnectSettings)
	repository := IssueRepository{Db: db}

	service1 := NotificationService{
		repository: &repository,
		botApiService: &botApiService,
		statusName: reviewStatus,
		telegramNickNameMap: telegramNickNameMap}
	service2 := NotificationService{
		repository: &repository,
		botApiService: &botApiService,
		statusName: toDoStatus,
		telegramNickNameMap: telegramNickNameMap}
	service3 := NotificationService{
		repository: &repository,
		botApiService: &botApiService,
		statusName: testingStatus,
		telegramNickNameMap: telegramNickNameMap}
	service4 := NotificationService{
		repository: &repository,
		botApiService: &botApiService,
		statusName: codeReviewStatus,
		telegramNickNameMap: telegramNickNameMap}

	for {
		service1.CheckUpdateIssues()
		service2.CheckUpdateIssues()
		service3.CheckUpdateIssues()
		service4.CheckUpdateIssues()
		time.Sleep(2000 * time.Millisecond)
	}

	runRestHealthEndpoint()
}

func runRestHealthEndpoint() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	http.ListenAndServe(":8088", nil)
}