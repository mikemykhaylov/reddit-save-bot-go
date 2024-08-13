package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gorilla/mux"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/api"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
	"github.com/spf13/viper"
)

var (
	telegramAPI *api.TelegramAPI
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go!")
}

func botHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var update gotgbot.Update

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&update)

	if err != nil {
		log.Error("Failed to decode update", "cause", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// print the sender id
	log.Info("Received message", "sender", update.Message.From.Id)

	// send a message back
	personalID := viper.GetInt64("personalID")

	err = telegramAPI.SendMessage(ctx, personalID, "Hello, Go!")
	if err != nil {
		log.Error("Failed to send message", "cause", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return 200
	w.WriteHeader(http.StatusOK)
}

func NewServer(config *config.ServerConfig) error {
	telegramAPI = api.NewTelegramAPI(
		viper.GetString("token"),
		&http.Client{},
	)

	r := mux.NewRouter()

	helloHandleFunc := logger.WithLogging(helloHandle)
	r.HandleFunc("/", helloHandleFunc)

	botHandleFunc := logger.WithLogging(botHandle)
	r.HandleFunc("/webhook", botHandleFunc).Methods(http.MethodPost)

	address := fmt.Sprintf("0.0.0.0:%d", config.Port)

	return http.ListenAndServe(address, r)
}
