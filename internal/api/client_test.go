package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://example.com"
	client := NewClient(baseURL)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestPingSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			t.Errorf("Expected path /ping, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Ping()

	if err != nil {
		t.Errorf("Ping should succeed, got error: %v", err)
	}
}

func TestPingFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Ping()

	if err == nil {
		t.Error("Ping should fail with server error")
	}
}

func TestGetChallengeSuccess(t *testing.T) {
	expectedResponse := ChallengeResponse{
		UUID: "test-uuid",
		SetX: "4",
		SetY: "4",
		SetZ: "5",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/challenge-me-easy" {
			t.Errorf("Expected path /challenge-me-easy, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var request ChallengeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if request.UUID != "test-uuid" || request.User != "test-user" {
			t.Errorf("Unexpected request values: %+v", request)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	response, err := client.GetChallenge("test-uuid", "test-user")

	if err != nil {
		t.Errorf("GetChallenge should succeed, got error: %v", err)
		return
	}

	if response.UUID != expectedResponse.UUID ||
		response.SetX != expectedResponse.SetX ||
		response.SetY != expectedResponse.SetY ||
		response.SetZ != expectedResponse.SetZ {
		t.Errorf("Response mismatch. Expected %+v, got %+v", expectedResponse, response)
	}
}

func TestGetChallengeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetChallenge("test-uuid", "test-user")

	if err == nil {
		t.Error("GetChallenge should fail with bad request")
	}
}

func TestSubmitSolutionSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/challenge-me-easy" {
			t.Errorf("Expected path /challenge-me-easy, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var request SolutionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		expectedUUID := "test-uuid"
		expectedResult := "0,1,2,4,6,8,9,10"
		expectedHash := "test-hash"

		if request.UUID != expectedUUID ||
			request.Result != expectedResult ||
			request.Hash != expectedHash {
			t.Errorf("Unexpected request values: %+v", request)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Solution accepted"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.SubmitSolution("test-uuid", "0,1,2,4,6,8,9,10", "test-hash")

	if err != nil {
		t.Errorf("SubmitSolution should succeed, got error: %v", err)
	}
}

func TestSubmitSolutionFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid solution"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.SubmitSolution("test-uuid", "wrong-result", "wrong-hash")

	if err == nil {
		t.Error("SubmitSolution should fail with bad request")
	}
}
