package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gorilla/mux"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/api"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/handler"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
	"github.com/spf13/viper"
)

var (
	messageHandler *handler.MessageHandler
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

	if update.Message == nil {
		log.Warn("Update doesn't contain a message, quietly ignoring")
		w.WriteHeader(http.StatusOK)
		return
	}

	err = messageHandler.HandleMessage(ctx, update.Message)
	if err != nil {
		// we only return error if we'd like to retry
		log.Error("Failed to handle message", "cause", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return 200
	w.WriteHeader(http.StatusOK)
}

func NewServer(serverConfig *config.ServerConfig) error {
	httpClient := &http.Client{}

	telegramAPI := api.NewTelegramAPI(viper.GetString(config.TelegramBotTokenKey), httpClient)
	redditAPI := api.NewRedditAPI(viper.GetString(config.RedditClientIDKey), viper.GetString(config.RedditClientSecretKey), httpClient)

	messageHandler = handler.NewMessageHandler(telegramAPI, redditAPI)

	r := mux.NewRouter()

	helloHandleFunc := logger.WithLogging(helloHandle)
	r.HandleFunc("/", helloHandleFunc)

	botHandleFunc := logger.WithLogging(botHandle)
	r.HandleFunc("/webhook", botHandleFunc).Methods(http.MethodPost)

	address := fmt.Sprintf("0.0.0.0:%d", serverConfig.Port)

	return http.ListenAndServe(address, r)
}
