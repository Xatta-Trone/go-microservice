package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	type mailMsg struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var rqPayload mailMsg

	err := app.readJSON(w, r, &rqPayload)

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    rqPayload.From,
		To:      rqPayload.To,
		Subject: rqPayload.Subject,
		Data:    rqPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to" + rqPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
