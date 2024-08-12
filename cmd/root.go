package cmd

import (
	"context"
	"errors"

	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "reddit-save-bot",
		Short:   "Reddit Save Bot is a bot that receives a Telegram message with a Reddit link to a video and sends it back",
		Version: server.Version(),
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := viper.GetUint16("port")

			logger.FromContext(context.Background()).Info("Starting server", "port", port)

			token := viper.GetString("token")
			if token == "" {
				logger.FromContext(context.Background()).Error("Telegram bot token is required")
				return errors.New("telegram bot token is required")
			}

			personalID := viper.GetString("personalID")
			if personalID == "" {
				logger.FromContext(context.Background()).Warn("Telegram personal ID is not set, bot in public mode")
			}

			config := &config.ServerConfig{
				Port: port,
			}

			return server.NewServer(config)
		},
	}
)

func init() {
	serveCmd.Flags().Uint16P("port", "p", 8080, "Port to listen on")
	if err := viper.BindPFlag("port", serveCmd.Flags().Lookup("port")); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("port", "PORT"); err != nil {
		panic(err)
	}
	viper.SetDefault("port", 8080)

	serveCmd.Flags().StringP("token", "t", "", "Telegram bot token")
	if err := viper.BindPFlag("token", serveCmd.Flags().Lookup("token")); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("token", "TELEGRAM_BOT_TOKEN"); err != nil {
		panic(err)
	}

	serveCmd.Flags().StringP("personalID", "", "", "Telegram personal ID")
	if err := viper.BindPFlag("personalID", serveCmd.Flags().Lookup("personalID")); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("personalID", "PERSONAL_ID"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(serveCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
