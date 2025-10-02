package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type MailmanRequest struct {
	ToAddresses  []string `json:"to_addresses"`
	CCAddresses  []string `json:"cc_addresses,omitempty"`
	BCCAddresses []string `json:"bcc_addresses,omitempty"`
	ReplyTo      []string `json:"reply_to,omitempty"`
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
	IsHTMLBody   bool     `json:"is_html_body"`
}

type MailmanClient struct {
	url    string
	apiKey string
	client *http.Client
}

func NewMailmanClient(url, apiKey string) *MailmanClient {
	return &MailmanClient{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (m *MailmanClient) SendEmail(req *MailmanRequest) error {
	if len(req.ToAddresses) == 0 {
		return fmt.Errorf("to_addresses is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if req.Body == "" {
		return fmt.Errorf("body is required")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := strings.TrimSuffix(m.url, "/")
	endpoint := baseURL + "/data/mailman"

	log.Printf("Mailman: Sending email to %s via %s", req.ToAddresses, endpoint)
	log.Printf("Mailman: Request payload: %s", string(jsonData))

	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("x-api-key", m.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(httpReq)
	if err != nil {
		log.Printf("Mailman: HTTP request failed: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Mailman: Response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Mailman: Error response body: %s", string(body))
		return fmt.Errorf("email service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
