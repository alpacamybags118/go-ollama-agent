package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) SendRequest(endpoint string, payload interface{}) (*http.Response, error) {
	url := c.BaseURL + endpoint
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

func (c *Client) GetResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

// CompletionRequest represents a request to the Ollama completion API
type CompletionRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// StreamingCompletionResponse represents each chunk in a streaming response
type StreamingCompletionResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

func (c *Client) CreateCompletionStream(req CompletionRequest) (chan StreamingCompletionResponse, chan error) {
	responseChannel := make(chan StreamingCompletionResponse)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(responseChannel)
		defer close(errorChannel)

		resp, err := c.SendRequest("/api/generate", req)
		if err != nil {
			errorChannel <- err
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			errorChannel <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			return
		}

		decoder := json.NewDecoder(resp.Body)

		for {
			var chunk StreamingCompletionResponse
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					return
				}
				errorChannel <- fmt.Errorf("error decoding stream: %w", err)
				return
			}

			// Send the chunk to the channel
			responseChannel <- chunk

			// Check if this is the end of the stream
			if chunk.Done {
				return
			}
		}
	}()

	return responseChannel, errorChannel
}
