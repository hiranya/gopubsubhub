package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const errorErrorDetailsRequiredInBodyResponse = "Error response on subscription should return the errors in the response body"

func TestSubscribeSuccess(t *testing.T) {
	handler := mainHandler()
	data := url.Values{}
	data.Set("hub.mode", "subscribe")
	data.Set("hub.callback", "http://localhost:7000")
	data.Set("hub.topic", "http://localhost:8000")
	data.Set("hub.lease_seconds", "3600") // in seconds, 1-hour
	data.Set("hub.secret", "secret-string-1")

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Errorf("Subscribe request did not return %v", http.StatusAccepted)
	}
}

func TestSubscribeFailure(t *testing.T) {
	handler := mainHandler()
	data := url.Values{}
	data.Set("hub.mode", "subscribe")
	// ommitting required fields here to trigger a failure

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Incomplete subscribe request did not return %v", http.StatusBadRequest)
	}

	// testing for the existence of just one error string would do
	if !strings.Contains(w.Body.String(), errorRequiredFieldMissingHubCallback) {
		t.Errorf(errorErrorDetailsRequiredInBodyResponse)
	}
}

func TestUnsubscribe(t *testing.T) {
	handler := mainHandler()
	data := url.Values{}
	data.Set("hub.mode", "unsubscribe")

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Unsubscribe request did not return %v", http.StatusOK)
	}
}

func TestPublish(t *testing.T) {
	handler := mainHandler()
	data := url.Values{}
	data.Set("hub.mode", "publish")

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Publish request did not return %v", http.StatusOK)
	}
}

func logi(s string) {
	log.Info(">> Test Suite:", s)
}
