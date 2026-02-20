package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type OdooService struct {
	url      string
	db       string
	user     string
	password string
	uid      int
	client   *http.Client
}

func NewOdooService(url, db, user, password string) *OdooService {
	return &OdooService{
		url:      url,
		db:       db,
		user:     user,
		password: password,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

type OdooRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      uint64      `json:"id"`
}

type OdooLoginResponse struct {
	Result int        `json:"result"`
	Error  *OdooError `json:"error,omitempty"`
}

type OdooError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Message string `json:"message"`
		Name    string `json:"name"`
		Debug   string `json:"debug"`
	} `json:"data"`
}

func (s *OdooService) Login() error {
	if s.url == "" {
		return fmt.Errorf("odoo url is not configured")
	}

	payload := OdooRequest{
		Jsonrpc: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "common",
			"method":  "authenticate",
			"args":    []interface{}{s.db, s.user, s.password, map[string]interface{}{}},
		},
		Id: 1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jsonrpc", s.url), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("odoo returned non-200 status code for login: %d", resp.StatusCode)
	}

	var loginResp OdooLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	if loginResp.Error != nil {
		return fmt.Errorf("odoo login error: %s", loginResp.Error.Data.Message)
	}

	if loginResp.Result == 0 {
		return fmt.Errorf("odoo login failed: invalid credentials")
	}

	s.uid = loginResp.Result
	return nil
}

type OdooCreateTicketResponse struct {
	Result int        `json:"result"` // ID of the created record
	Error  *OdooError `json:"error,omitempty"`
}

func (s *OdooService) CreateTicket(title string, description string, teamID int, userID int) (int, error) {
	if s.uid == 0 {
		return 0, fmt.Errorf("not logged into odoo")
	}

	// In Odoo, tickets are usually managed by the 'helpdesk.ticket' or 'project.task' module.
	// Since the codebase refers to 'ticketing management' we will try 'helpdesk.ticket'
	// If the user's Odoo doesn't use this, they can modify this code or let us know.
	model := "helpdesk.ticket"

	args := []interface{}{
		s.db,
		s.uid,
		s.password,
		model,
		"create",
		[]map[string]interface{}{
			{
				"name":        title,
				"description": description,
				"team_id":     teamID,
				"user_id":     userID,
			},
		},
	}

	payload := OdooRequest{
		Jsonrpc: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args":    args,
		},
		Id: 2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal create ticket payload: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jsonrpc", s.url), bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create create ticket request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute create ticket request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("odoo returned non-200 status code for create ticket: %d", resp.StatusCode)
	}

	var createResp OdooCreateTicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return 0, fmt.Errorf("failed to decode create ticket response: %w", err)
	}

	if createResp.Error != nil {
		return 0, fmt.Errorf("odoo created ticket error: %s", createResp.Error.Data.Message)
	}

	log.Printf("[INFO] Created Odoo ticket %d under %s", createResp.Result, model)
	return createResp.Result, nil
}
