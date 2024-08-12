package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TelegramAPI struct {
	baseURL    string
	httpClient HTTPClient
}

func NewTelegramAPI(token string, httpClient HTTPClient) *TelegramAPI {
	return &TelegramAPI{
		baseURL:    "https://api.telegram.org/bot" + token,
		httpClient: &http.Client{},
	}
}

func (t *TelegramAPI) SendMessage(ctx context.Context, chatID int64, text string) error {
	methodURL, _ := url.Parse(fmt.Sprintf("%s/sendMessage", t.baseURL))
	query := methodURL.Query()
	query.Set("chat_id", fmt.Sprintf("%d", chatID))
	query.Set("text", text)
	methodURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, methodURL.String(), nil)
	if err != nil {
		return err
	}

	_, err = t.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
