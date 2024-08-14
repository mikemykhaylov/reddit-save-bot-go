package config

var (
	TelegramBotTokenKey string = "token"
	PersonalIDKey       string = "personalID"
)

type ServerConfig struct {
	Port uint16
}
