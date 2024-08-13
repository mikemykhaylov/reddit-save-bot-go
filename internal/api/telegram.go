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
		httpClient: httpClient,
	}
}

func (t *TelegramAPI) SendMessage(ctx context.Context, chatID int64, text string) error {
	requestURL, _ := url.Parse(fmt.Sprintf("%s/sendMessage", t.baseURL))
	query := requestURL.Query()
	query.Set("chat_id", fmt.Sprintf("%d", chatID))
	query.Set("text", text)
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return err
	}

	res, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
