package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Verification wraps a client to be able to use local verify methods.
type Verification struct {
	client *Client
}

// MarshalJSON returns a byte slice with the serialized JSON of the
// VerifyMessageRequest struct.
func (m *VerifyMessageRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
		VerifyMessageRequest
	}{
		ApiKey:               m.apiKey,
		ApiSecret:            m.apiSecret,
		VerifyMessageRequest: *m,
	})
}

// VerifyMessageRequest is the request struct for initiating the verification process
// for a phone number.
type VerifyMessageRequest struct {
	apiKey    string
	apiSecret string

	Number        string `json:"number"`
	Brand         string `json:"brand"`
	SenderID      string `json:"sender_id,omitempty"`
	Country       string `json:"country,omitempty"`
	Language      string `json:"lg,omitempty"`
	CodeLength    int    `json:"code_length,omitempty"`
	PINExpiry     int    `json:"pin_expiry,omitempty"`
	NextEventWait int    `json:"next_event_wait,omitempty"`
}

// VerifyMessageResponse is the struct for the response from the verify
// endpoint.
type VerifyMessageResponse struct {
	Status    ResponseCode `json:"status,string"`
	RequestID string       `json:"request_id"`
	ErrorText string       `json:"error_text"`
}

// Send makes the actual HTTP request to the endpoint and returns the
// response.
func (c *Verification) Send(m *VerifyMessageRequest) (*VerifyMessageResponse, error) {
	if len(m.Number) == 0 {
		return nil, errors.New("Invalid Number field specified")
	}

	if len(m.Brand) == 0 {
		return nil, errors.New("Invalid Brand field specified")
	}

	var verifyMessageResponse *VerifyMessageResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("Invalid message struct. Cannot convert to json.")
	}
	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HttpClient.Do(r)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifyMessageResponse)
	if err != nil {
		return nil, err
	}
	return verifyMessageResponse, nil
}

func (m *VerifyCheckRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
		VerifyCheckRequest
	}{
		ApiKey:             m.apiKey,
		ApiSecret:          m.apiSecret,
		VerifyCheckRequest: *m,
	})
}

type VerifyCheckRequest struct {
	apiKey    string
	apiSecret string

	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	IPAddress string `json:"ip_address,omitempty"`
}

type VerifyCheckResponse struct {
	Status    ResponseCode `json:"status,string"`
	EventID   string       `json:"event_id"`
	Price     string       `json:"price"`
	Currency  string       `json:"currency"`
	ErrorText string       `json:"error_text"`
}

// Check is a request to send a code to the Nexmo and verify a code
// for the request.
func (c *Verification) Check(m *VerifyCheckRequest) (*VerifyCheckResponse, error) {
	if len(m.RequestID) == 0 {
		return nil, errors.New("Invalid RequestID field specified")
	}

	if len(m.Code) == 0 {
		return nil, errors.New("Invalid Code field specified")
	}

	var verifyCheckResponse *VerifyCheckResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("Invalid message struct. Cannot convert to json.")
	}
	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/check/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HttpClient.Do(r)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifyCheckResponse)
	if err != nil {
		return nil, err
	}
	return verifyCheckResponse, nil
}

func (m *VerifySearchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
		VerifySearchRequest
	}{
		ApiKey:              m.apiKey,
		ApiSecret:           m.apiSecret,
		VerifySearchRequest: *m,
	})
}

type VerifySearchRequest struct {
	apiKey    string
	apiSecret string

	RequestID string `json:"request_id,omitempty"`
}

type VerifySearchResponse struct {
	RequestID      string   `json:"request_id"`
	AccountID      string   `json:"account_id"`
	Number         string   `json:"number"`
	SenderID       string   `json:"sender_id"`
	DateSubmitted  string   `json:"date_submitted"`
	DateFinalized  string   `json:"date_finalized"`
	FirstEventDate string   `json:"first_event_date"`
	LastEventDate  string   `json:"last_event_date"`
	Status         string   `json:"status"`
	Checks         []Checks `json:"checks"`
	Price          string   `json:"price"`
	Currency       string   `json:"currency"`
	ErrorText      string   `json:"error_text"`
}

type Checks struct {
	DateReceived string `json:"date_received"`
	Code         string `json:"code"`
	Status       string `json:"status"`
	IPAddress    string `json:"ip_address,omitempty"`
}

// Search sends the verify search request to Nexmo.
func (c *Verification) Search(m *VerifySearchRequest) (*VerifySearchResponse, error) {
	var verifySearchResponse *VerifySearchResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("Invalid message struct. Cannot convert to json.")
	}
	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/search/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HttpClient.Do(r)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifySearchResponse)
	if err != nil {
		return nil, err
	}
	return verifySearchResponse, nil
}

type VerifyControlRequest struct {
	apiKey    string
	apiSecret string

	RequestID string `json:"request_id"`
	Cmd       Cmd    `json:"cmd"`
}

// MarshalJSON returns a byte slice with the serialized JSON of the
// VerifyMessageRequest struct.
func (m *VerifyControlRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
		VerifyControlRequest
	}{
		ApiKey:               m.apiKey,
		ApiSecret:            m.apiSecret,
		VerifyControlRequest: *m,
	})
}

type Cmd int

const (
	CmdUnknown Cmd = iota // default value
	CmdCancel
	CmdTriggerNextEvent
)

var cmdToStringMap = map[Cmd]string{
	CmdCancel:           "cancel",
	CmdTriggerNextEvent: "trigger_next_event",
}

var stringToCmdMap = map[string]Cmd{
	"cancel":             CmdCancel,
	"trigger_next_event": CmdTriggerNextEvent,
}

func (c Cmd) String() string {
	return cmdToStringMap[c]
}

func (c Cmd) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Cmd) UnmarshalJSON(data []byte) error {
	str := string(data)
	v, ok := stringToCmdMap[str]
	if !ok {
		return fmt.Errorf("unknown value %q for cmd", data)
	}
	*c = v
	return nil
}

func isValidCmd(c Cmd) bool {
	_, ok := cmdToStringMap[c]
	return ok
}

type VerifyControlResponse struct {
	Status    ResponseCode `json:"status,string"`
	Command   Cmd          `json:"command"`
	ErrorText string       `json:"error_text"`
}

// Control executes a cmd for a given request.
// Cancel command can be sent only in 30 seconds after sending a verify request.
// Doc: https://developer.nexmo.com/api/verify#verify-control
func (c *Verification) Control(m *VerifyControlRequest) (*VerifyControlResponse, error) {
	if !isValidCmd(m.Cmd) {
		return nil, fmt.Errorf("invalid cmd [%d]", m.Cmd)
	}

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("invalid message struct(cannot convert to json)")
	}
	r, err := http.NewRequest("POST", apiRootv2+"/verify/control/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var controlResp VerifyControlResponse
	if err := json.NewDecoder(resp.Body).Decode(&controlResp); err != nil {
		return nil, err
	}
	return &controlResp, nil
}
