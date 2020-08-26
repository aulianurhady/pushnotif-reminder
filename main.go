package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron"
)

type Request struct {
	Notification PayloadFCM `json:"notification"`
	To           string     `json:"to"`
}

type PayloadFCM struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	Icon        string `json:"icon,omitempty"`
	ClickAction string `json:"click_action,omitempty"`
}

type Response struct {
	MulticastID  int64           `json:"multicast_id"`
	Success      int             `json:"success"`
	Failure      int             `json:"failure"`
	CanonicalIds int             `json:"canonical_ids"`
	Results      []PayloadResult `json:"results"`
}

type PayloadResult struct {
	MessageID string `json:"message_id"`
}

func main() {
	c := cron.New()
	c.AddFunc("7 9 * * *", func() {
		sendPushNotif()
	})
	log.Println("Cron is starting...")
	c.Start()
	time.Sleep(2 * time.Minute)
}

func sendPushNotif() {
	reqStruct := &Request{
		Notification: PayloadFCM{
			Title: "Reminder Meeting",
			Body:  "Sebentar lagi meeting, brow!",
		},
		To: "",
	}
	requestBody, err := json.Marshal(reqStruct)
	if err != nil {
		log.Fatalln(err)
	}

	httpRequest, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "")
	if err != nil {
		log.Fatalln(err)
	}

	client := http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var respBody Response
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Body: ", respBody)
}
