package fcm

import (
	"net/http"

	"github.com/pkg/errors"
)

// Notification specifies the predefined, user-visible key-value pairs of the
// notification payload.
type Notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	ChannelID    string `json:"android_channel_id,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Image        string `json:"image,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
}

// Message represents list of targets, options, and payload for HTTP JSON
// messages.
type Message struct {
	To                       string                 `json:"to,omitempty"`
	RegistrationIDs          []string               `json:"registration_ids,omitempty"`
	Condition                string                 `json:"condition,omitempty"`
	CollapseKey              string                 `json:"collapse_key,omitempty"`
	Priority                 string                 `json:"priority,omitempty"`
	ContentAvailable         bool                   `json:"content_available,omitempty"`
	MutableContent           bool                   `json:"mutable_content,omitempty"`
	DelayWhileIdle           bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive               *uint                  `json:"time_to_live,omitempty"`
	DeliveryReceiptRequested bool                   `json:"delivery_receipt_requested,omitempty"`
	DryRun                   bool                   `json:"dry_run,omitempty"`
	RestrictedPackageName    string                 `json:"restricted_package_name,omitempty"`
	Notification             *Notification          `json:"notification,omitempty"`
	Data                     map[string]interface{} `json:"data,omitempty"`
	Apns                     map[string]interface{} `json:"apns,omitempty"`
	Webpush                  map[string]interface{} `json:"webpush,omitempty"`
}

// NotifyDevice allows a client notify a single device
func (c *Client) NotifyDevice(deviceId string, optionalData map[string]interface{}, messageTitle string) (Response, error) {
	msg := Message{
		To: deviceId,
		Notification: &Notification{
			Title: messageTitle,
		},
		Data: optionalData,
	}
	url := "fcm/send"

	var resp Response
	if err := c.makeRequest(http.MethodPost, url, msg, &resp); err != nil {
		return Response{}, errors.Wrap(err, "error making fcm request to notify single device")
	}
	if resp.Error != nil {
		return resp, resp.Error
	}
	return resp, nil
}

// NotifyDevices allows a client notify multiple devices at a time
func (c *Client) NotifyDevices(deviceIds []string, optionalData map[string]interface{}, messageTitle string) (Response, error) {
	msg := Message{
		RegistrationIDs: deviceIds,
		Notification: &Notification{
			Title: messageTitle,
		},
		Data: optionalData,
	}
	url := "fcm/send"

	var resp Response
	if err := c.makeRequest(http.MethodPost, url, msg, &resp); err != nil {
		return Response{}, errors.Wrap(err, "error making fcm request to notify multiple devices")
	}
	if resp.Error != nil {
		return resp, resp.Error
	}
	return resp, nil
}

// SendCustomMessage allows you define message attributes as needed
func (c Client) SendCustomMessage(msg Message) (Response, error) {
	url := "fcm/send"

	var resp Response
	if err := c.makeRequest(http.MethodPost, url, msg, &resp); err != nil {
		return Response{}, errors.Wrap(err, "error making fcm request to notify single device")
	}
	if resp.Error != nil {
		return resp, resp.Error
	}
	return resp, nil
}
