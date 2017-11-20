package nexmo

import (
	"testing"
	"time"
)

func ensureClient(t *testing.T) *Client {
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	client, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}
	return client
}

func testSend(t *testing.T) *VerifyMessageResponse {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	client := ensureClient(t)

	message := &VerifyMessageRequest{
		Number:   TEST_PHONE_NUMBER,
		Brand:    TEST_FROM,
		SenderID: TEST_FROM,
	}

	messageResponse, err := client.Verify.Send(message)
	if err != nil {
		t.Error("Failed to send verification request with error:", err)
	}

	return messageResponse
}

func TestSend(t *testing.T) {
	messageResponse := testSend(t)
	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}

func TestSendCheck(t *testing.T) {
	// We need the request id, so we have to run this first.
	sendResponse := testSend(t)

	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation

	client := ensureClient(t)

	message := &VerifyCheckRequest{
		RequestID: sendResponse.RequestID,
		Code:      "1122", // Take a random code here, the number will not be verified properly though.
	}

	messageResponse, err := client.Verify.Check(message)
	if err != nil {
		t.Error("Failed to send verification check request with error:", err)
	}

	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}

func TestSendSearch(t *testing.T) {
	// We need the request id, so we have to run this first.
	sendResponse := testSend(t)

	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation

	client := ensureClient(t)

	message := &VerifySearchRequest{
		RequestID: sendResponse.RequestID,
	}

	messageResponse, err := client.Verify.Search(message)
	if err != nil {
		t.Error("Failed to send verification search request with error:", err)
	}

	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}

// TestControl checks both CmdCancel and CmdTriggerNextEvent event.
// This test causes 30 second sleep for cancel to be properly done.
func TestControl(t *testing.T) {
	client := ensureClient(t)

	testCases := map[string]struct {
		cmd     Cmd
		timeout time.Duration
	}{
		"cancel": {
			cmd:     CmdCancel,
			timeout: 30 * time.Second,
		},
		"trigger_next_event": {
			cmd: CmdTriggerNextEvent,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sendResp := testSend(t)

			time.Sleep(tc.timeout)

			req := &VerifyControlRequest{
				RequestID: sendResp.RequestID,
				Cmd:       tc.cmd,
			}
			controlResp, err := client.Verify.Control(req)
			if err != nil {
				t.Fatal("Failed to send a verification control with error:", err)
			}
			if controlResp.Status != ResponseSuccess {
				t.Errorf("Control status is not success. Got=%s. Error text=%s", controlResp.Status, controlResp.ErrorText)
			}
		})
	}
}
