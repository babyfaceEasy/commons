package twilio

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "https://api.twilio.com/2010-04-01"
)

// Config is a representation of a twilio config
type Config struct {
	AccountSid   string
	AuthToken    string
	BaseUrl      string
	PhoneNumber  string
	APIKeySid    string
	APIKeySecret string
}

// Client is a representation of an twilio client
type Client struct {
	Config Config
	Client *http.Client
}

// ConfigFromEnvVars provides the default config from env vars
func ConfigFromEnvVars() Config {
	return Config{
		AccountSid:  os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
		PhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		BaseUrl:     baseURL,
	}
}

// NewClient creates an twilio client using configuration variables
func NewClient() Client {
	return Client{Config: ConfigFromEnvVars(), Client: &http.Client{Timeout: 30 * time.Second}}
}

// SmsResponse is returned after a text/sms message is posted to Twilio
type SmsResponse struct {
	Sid         string `json:"sid"`
	DateCreated string `json:"date_created"`
	DateUpdate  string `json:"date_updated"`
	DateSent    string `json:"date_sent"`
	AccountSid  string `json:"account_sid"`
	To          string `json:"to"`
	From        string `json:"from"`
	NumMedia    string `json:"num_media"`
	Body        string `json:"body"`
	Status      string `json:"status"`
	Direction   string `json:"direction"`
	ApiVersion  string `json:"api_version"`
	Price       string `json:"price"`
	Url         string `json:"uri"`
}

// Exception is a representation of a json object returned for twilio error
type Exception struct {
	Status   string `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

// Print the RESTException in a human-readable form.
func (r Exception) Error() string {
	var errorCode int

	if r.Code != errorCode {
		return fmt.Sprintf("Code %d: %s", r.Code, r.Message)
	} else if r.Status != "" {
		return fmt.Sprintf("Status %s: %s", r.Status, r.Message)
	}
	return r.Message
}
