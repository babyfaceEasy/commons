package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Config is a representation of config variables
type Config struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseURL"`
}

// Response represents the FCM server's response to the application
// server's sent message.
type Response struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`

	// Device Group HTTP Response
	FailedRegistrationIDs []string `json:"failed_registration_ids"`

	// Topic HTTP response
	MessageID         int64 `json:"message_id"`
	Error             error `json:"error"`
	ErrorResponseCode string
}

// Result represents the status of a processed message.
type Result struct {
	MessageID         string `json:"message_id"`
	RegistrationID    string `json:"registration_id"`
	Error             string `json:"error"`
	ErrorResponseCode string
}

// ConfigFromEnvVars provides the default config from env vars
func ConfigFromEnvVars() Config {
	baseURL := os.Getenv("FCM_URL")
	if baseURL == "" {
		baseURL = "https://fcm.googleapis.com"
	}
	return Config{
		APIKey:  os.Getenv("FCM_API_KEY"),
		BaseURL: baseURL,
	}
}

// Client is a representation of an fcm client
type Client struct {
	Config Config
	Client *http.Client
}

// NewClient creates an fcm client using configuration variables
func NewClient() Client {
	return Client{Config: ConfigFromEnvVars(), Client: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) makeRequest(method, rURL string, reqBody interface{}, resp interface{}) error {
	URL := fmt.Sprintf("%s/%s", c.Config.BaseURL, rURL)
	var body io.Reader
	if reqBody != nil {
		bb, err := json.Marshal(reqBody)
		if err != nil {
			return errors.Wrap(err, "client - unable to marshal request struct")
		}
		body = bytes.NewReader(bb)
	}
	req, err := http.NewRequest(method, URL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", c.Config.APIKey))
	if err != nil {
		return errors.Wrap(err, "client - unable to create request body")
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client - failed to execute request")
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != 204 {
		return errors.Errorf("invalid status code received, expected 200/204, got %v", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return errors.Wrap(err, "unable to unmarshal request body")
	}
	return nil
}
