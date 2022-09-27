package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func (twilio Config) getBasicAuthCredentials() (string, string) {
	if twilio.APIKeySid != "" {
		return twilio.APIKeySid, twilio.APIKeySecret
	}
	return twilio.AccountSid, twilio.AuthToken
}

func (c Client) makeRequest(formValues url.Values) (SmsResponse, error) {
	twilioUrl := c.Config.BaseUrl + "/Accounts/" + c.Config.AccountSid + "/Messages.json"

	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return SmsResponse{}, errors.Wrap(err, "unable to make request with form values")
	}
	req.SetBasicAuth(c.Config.getBasicAuthCredentials())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := c.Client.Do(req)
	if err != nil {
		return SmsResponse{}, errors.Wrap(err, "unable to do request")
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK && result.StatusCode != http.StatusCreated {
		exception := new(Exception)
		if err := json.NewDecoder(result.Body).Decode(exception); err != nil {

			bb, err := ioutil.ReadAll(result.Body)
			if err != nil {
				return SmsResponse{}, errors.Wrapf(err, "unable to unmarshal exception response and can't read body, status code=%d", result.StatusCode)
			}
			return SmsResponse{}, fmt.Errorf("unexpected response code %d, body %s", result.StatusCode, bb)
		}
		return SmsResponse{}, exception
	}

	smsResponse := new(SmsResponse)
	if err := json.NewDecoder(result.Body).Decode(smsResponse); err != nil {
		return SmsResponse{}, errors.Wrap(err, "unable to unmarshal response")
	}

	return *smsResponse, nil
}

// SendSMS sends an sms to a provided phone number
func (c Client) SendSMS(to, body string) (SmsResponse, error) {
	formValues := url.Values{
		"From": []string{c.Config.PhoneNumber},
		"To":   []string{to},
		"Body": []string{body},
	}
	return c.makeRequest(formValues)
}
