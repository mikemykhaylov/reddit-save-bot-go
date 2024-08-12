package cmd

import (
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
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

			config := &config.ServerConfig{
				Port: port,
			}

			return server.NewServer(config)
		},
	}
)

func init() {
	serveCmd.Flags().Uint16P("port", "p", 8080, "Port to listen on")
	// check ENV variable with viper
	if err := viper.BindPFlag("port", serveCmd.Flags().Lookup("port")); err != nil {
		panic(err)
	}
	viper.SetDefault("port", 8080)
	if err := viper.BindEnv("port", "PORT"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(serveCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
