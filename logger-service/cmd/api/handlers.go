package main

import (
	"log"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) writeLog(w http.ResponseWriter, r *http.Request) {
	// read json int var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Printf("handlers.writeLog: err:%v\n", err)
		_ = app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	log.Printf("handlers.writeLog: event%v\n", event)

	_ = app.writeJSON(w, http.StatusAccepted, resp)
}
