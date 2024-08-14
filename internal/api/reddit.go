package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mikemykhaylov/reddit-save-bot-go/internal/logger"
)

var (
	RedditBaseURL string = "https://www.reddit.com/api/v1"
	UserAgent     string = "reddit-save-bot-go"
)

type RedditAPI struct {
	clientID     string
	clientSecret string

	userAgent  string
	baseURL    string
	httpClient *http.Client

	token          string
	tokenExpiresAt time.Time
}

type TokenGrant struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewRedditAPI(clientID, clientSecret string, httpClient *http.Client) *RedditAPI {
	return &RedditAPI{
		clientID:     clientID,
		clientSecret: clientSecret,
		userAgent:    UserAgent,
		baseURL:      RedditBaseURL,
		httpClient:   httpClient,
	}
}

func (r *RedditAPI) GetToken(ctx context.Context) (string, error) {
	log := logger.FromContext(ctx)

	errors := []string{}

	if r.clientID == "" {
		errors = append(errors, "clientID")
	}
	if r.clientSecret == "" {
		errors = append(errors, "clientSecret")
	}
	if len(errors) > 0 {
		err := fmt.Errorf("missing required fields: %s", strings.Join(errors, ", "))
		return "", err
	}

	if r.token != "" {
		if time.Now().Before(r.tokenExpiresAt) {
			return r.token, nil
		}
		log.Info("Refreshing Reddit access token")
	} else {
		log.Info("Fetching Reddit access token")
	}

	requestURL, _ := url.Parse(fmt.Sprintf("%s/access_token", r.baseURL))

	values := url.Values{}
	values.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, requestURL.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}

	// set basic auth
	req.SetBasicAuth(r.clientID, r.clientSecret)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", r.userAgent)

	res, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	// parse response
	var tokenGrant TokenGrant
	if err := json.NewDecoder(res.Body).Decode(&tokenGrant); err != nil {
		return "", err
	}

	r.token = tokenGrant.AccessToken
	r.tokenExpiresAt = time.Now().Add(time.Duration(tokenGrant.ExpiresIn) * time.Second)

	return tokenGrant.AccessToken, nil
}
