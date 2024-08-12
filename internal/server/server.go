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

	// send a message back
	personalID := viper.Get("personalID").(int64)

	telegramAPI.SendMessage(ctx, personalID, "Hello, Go!")

	// return 200
	w.WriteHeader(http.StatusOK)
}

func NewServer(config *config.ServerConfig) error {
	r := mux.NewRouter()

	helloHandleFunc := logger.WithLogging(helloHandle)
	r.HandleFunc("/", helloHandleFunc)

	botHandleFunc := logger.WithLogging(botHandle)
	r.HandleFunc("/webhook", botHandleFunc).Methods(http.MethodPost)

	address := fmt.Sprintf("0.0.0.0:%d", config.Port)

	return http.ListenAndServe(address, r)
}

func init() {
	token := viper.Get("token").(string)
	if token == "" {
		panic("token is required")
	}
	telegramAPI = api.NewTelegramAPI(token)
}
