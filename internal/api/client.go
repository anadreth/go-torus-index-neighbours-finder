package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type ChallengeRequest struct {
	UUID string `json:"uuid"`
	User string `json:"user"`
}

type ChallengeResponse struct {
	UUID string `json:"uuid"`
	SetX string `json:"set_x"`
	SetY string `json:"set_y"`
	SetZ string `json:"set_z"`
}

type SolutionRequest struct {
	UUID   string `json:"uuid"`
	Result string `json:"result"`
	Hash   string `json:"hash"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type debugTransport struct {
	http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s\n", reqDump)
	
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	
	respDump, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("RESPONSE:\n%s\n", respDump)
	
	return resp, err
}

func NewClient(baseURL string) *Client {
	transport := http.DefaultTransport
	if os.Getenv("DEBUG_HTTP") != "" {
		transport = &debugTransport{RoundTripper: transport}
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

func (c *Client) Ping() error {
	url := fmt.Sprintf("%s/ping", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to ping server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ping failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetChallenge(uuid, user string) (*ChallengeResponse, error) {
	url := fmt.Sprintf("%s/challenge-me-easy", c.baseURL)

	request := ChallengeRequest{
		UUID: uuid,
		User: user,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send challenge request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("challenge request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChallengeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode challenge response: %w", err)
	}

	return &response, nil
}

func (c *Client) SubmitSolution(uuid, result, hash string) error {
	url := fmt.Sprintf("%s/challenge-me-easy", c.baseURL)

	request := SolutionRequest{
		UUID:   uuid,
		Result: result,
		Hash:   hash,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal solution: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send solution: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solution submission failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Solution submitted successfully. Response: %s\n", string(body))
	return nil
}
