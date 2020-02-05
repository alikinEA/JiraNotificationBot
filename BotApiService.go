package main

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const TelegramUrl = "https://api.telegram.org/"
const JiraUrl = "https://jira.proitr.ru/browse/"
const TelegramUrlApiGetMe = "/getMe"
const TelegramUrlApiGetUpdates = "/getUpdates"

type BotApiService struct {
	botId     string
	proxyUser string
	proxyPass string
	proxyIp   string
	chatId    int
}

func (settings BotApiService) sendMessageToChat(message string) string {
	sendUrl := "/sendMessage?chat_id=" + strconv.Itoa(settings.chatId) + "&text=" + message
	return sendRequest(settings, sendUrl)
}

func (settings BotApiService) getMe() string {
	return sendRequest(settings, TelegramUrlApiGetMe)
}

func (settings BotApiService) getUpdates() string {
	return sendRequest(settings, TelegramUrlApiGetUpdates)
}

func sendRequest(settings BotApiService, method string) string {
	var auth = &proxy.Auth{User: settings.proxyUser, Password: settings.proxyPass}
	dialSocksProxy, err := proxy.SOCKS5("tcp", settings.proxyIp, auth, proxy.Direct)
	if err != nil {
		log.Fatalln("Error connecting to proxy:", err)
	}
	tr := &http.Transport{
		Dial:                dialSocksProxy.Dial,
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}

	// Create client
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	resp, err := client.Get(TelegramUrl + "bot" + settings.botId + method)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}
