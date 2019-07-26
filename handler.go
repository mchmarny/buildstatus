package main

import (
	"encoding/json"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	var c pushedContent
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logger.Printf("Error decoding message: %v", err)
		writeResp(w, http.StatusBadRequest, "Invalid Content")
		return
	}
	logger.Printf("Content from PubSub subscription: %v", c.Subscription)

	m := &CloudBuildNotification{}
	if err := json.Unmarshal(c.Message.Data, &m); err != nil {
		logger.Printf("Error decoding message data: %v", err)
		writeResp(w, http.StatusBadRequest, "Invalid Content")
		return
	}
	logger.Printf("Message content: %v", m)

	if err := send(m); err != nil {
		logger.Printf("Failed to send notification %v", err)
		writeResp(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	writeResp(w, http.StatusOK, "OK")
	return
}

func writeResp(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

type pushedContent struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}
