package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		log.Printf("New Request: action: %s, data:%v\n", requestPayload.Action, requestPayload.Auth)
		app.authenticate(w, requestPayload.Auth)
	case "log":
		log.Printf("New Request: action: %s, data:%v\n", requestPayload.Action, requestPayload.Log)
		app.logItem(w, requestPayload.Log)
	case "mail":
		log.Printf("New Request: action: %s, data:%v\n", requestPayload.Action, requestPayload.Log)
		app.sendMail(w, requestPayload.Mail)
	default:
		log.Printf("New Request: action '%s' unknown \n", requestPayload.Action)
		_ = app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll and send to the auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Printf("statusCode: %d, data:%v", response.StatusCode, response.Body)
		_ = app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from the auth-service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		_ = app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated!"
	payload.Data = jsonFromService.Data
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("App.logItem: create NewRequest err: ", err)
		_ = app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println("App.logItem: send Request err: ", err)
		_ = app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("App.logItem: error response (%v)\n", response.StatusCode)
		_ = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// cal the mail service
	mailServiceURL := "http://mail-service"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, err)
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
