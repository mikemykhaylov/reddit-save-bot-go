package handler

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/google/uuid"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/api"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/config"
	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
	"github.com/spf13/viper"
)

var (
	StartCommand               = "/start"
	StartCommandResponse       = "Hello, send a Reddit link to a video to begin"
	VideoSizeLimit       int64 = 50 * 1024 * 1024
)

type MessageHandler struct {
	TelegramAPI *api.TelegramAPI
	RedditAPI   *api.RedditAPI
}

func NewMessageHandler(telegramAPI *api.TelegramAPI, redditAPI *api.RedditAPI) *MessageHandler {
	return &MessageHandler{
		TelegramAPI: telegramAPI,
		RedditAPI:   redditAPI,
	}
}

func (m *MessageHandler) HandleMessage(ctx context.Context, message *gotgbot.Message) error {
	log := logger.FromContext(ctx)
	log.Info("Handling message", "messageText", message.Text)
	defer log.Info("Finished handling message")

	personalID := viper.GetInt64(config.PersonalIDKey)
	if personalID != api.TelegramPublicPersonalID && message.From.Id != personalID {
		log = log.With("personalID", message.From.Id, "message", message.Text)
		log.Warn("Message is not from personal ID, ignoring")
		return nil
	}

	if message.Text == StartCommand {
		log.Info("Received start command")
		return m.TelegramAPI.SendMessage(ctx, message.Chat.Id, StartCommandResponse)
	}

	postURL, err := url.Parse(message.Text)
	if err != nil {
		_ = m.TelegramAPI.SendMessage(ctx, message.Chat.Id, "Failed to parse URL")
		err = fmt.Errorf("failed to parse URL\n  caused by: %w", err)
		return err
	}

	if !strings.Contains(postURL.Hostname(), "reddit.com") {
		log.Warn("URL is not from Reddit, ignoring")
		return nil
	}

	// get reddit api token to authorize download
	token, err := m.RedditAPI.GetToken(ctx)
	if err != nil {
		log.Warn("Failed to get Reddit API token, download will likely fail", "cause", err)
	} else {
		// we got the token, so we replace the hostname with oauth.reddit.com
		postURL.Host = "oauth.reddit.com"
	}

	videoName := uuid.New().String()
	videoPath := fmt.Sprintf("/tmp/%s.mp4", videoName)

	args := []string{
		"yt-dlp",
		"--add-header",
		fmt.Sprintf("Authorization: Bearer %s", token),
		"--add-header",
		fmt.Sprintf("User-Agent: %s", api.UserAgent),
		"-P",
		"/tmp",
		"-o",
		videoPath,
		"-v",
		postURL.String(),
	}

	// run yt-dlp to download the video
	cmd := exec.Command(args[0], args[1:]...)
	outB, err := cmd.Output()
	log.Info("yt-dlp output", "stdout", string(outB))

	if err != nil {
		log.Error("Failed to download video", "cause", err)
		_ = m.TelegramAPI.SendMessage(ctx, message.Chat.Id, "Failed to download video")
		return nil
	}

	// get video file
	fi, err := os.Stat(fmt.Sprintf("/tmp/%s.mp4", videoName))
	if err != nil {
		_ = m.TelegramAPI.SendMessage(ctx, message.Chat.Id, "Failed to get video file")
		err = fmt.Errorf("failed to get video file\n  caused by: %w", err)
		return err
	}
	fsize := fi.Size()
	if fsize > VideoSizeLimit {
		log.Warn("Video is too large", "size", fsize)

		// remove video file
		if err := os.Remove(fmt.Sprintf("/tmp/%s.mp4", videoName)); err != nil {
			log.Error("Failed to remove video file", "cause", err)
		}

		_ = m.TelegramAPI.SendMessage(ctx, message.Chat.Id, "Video is too large")
		return nil
	}

	// send video
	if err := m.TelegramAPI.SendVideo(ctx, message.Chat.Id, videoPath); err != nil {
		_ = m.TelegramAPI.SendMessage(ctx, message.Chat.Id, "Failed to send video")
		err = fmt.Errorf("failed to send video\n  caused by: %w", err)
		return err
	}

	return nil
}
