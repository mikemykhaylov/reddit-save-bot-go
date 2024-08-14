package config

var (
	TelegramBotTokenKey   string = "token"
	PersonalIDKey         string = "personalID"
	RedditClientIDKey     string = "redditClientID"
	RedditClientSecretKey string = "redditClientSecret"
)

type ServerConfig struct {
	Port uint16
}
