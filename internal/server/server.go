package server

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go!")
}

func botHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var update gotgbot.Update

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&update)

	if err != nil {
		logger.FromContext(ctx).Error("Failed to decode update", "cause", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// print the sender id
	logger.FromContext(ctx).Info("Received message", "sender", update.Message.From.Id)

	// return 200
	w.WriteHeader(http.StatusOK)
}

func NewServer(config *config.ServerConfig) error {
	r := mux.NewRouter()

	helloHandleFunc := logger.WithLogging(helloHandle)
	r.HandleFunc("/hello", helloHandleFunc)

	botHandleFunc := logger.WithLogging(botHandle)
	r.HandleFunc("/", botHandleFunc).Methods(http.MethodPost)

	address := fmt.Sprintf("localhost:%d", config.Port)

	return http.ListenAndServe(address, r)
}
