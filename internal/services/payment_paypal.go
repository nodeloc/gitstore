package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nodeloc/git-store/internal/config"
)

type PayPalService struct {
	config  *config.Config
	baseURL string
}

type PayPalOrder struct {
	ID            string                  `json:"id"`
	Status        string                  `json:"status"`
	Intent        string                  `json:"intent"`
	PurchaseUnits []PayPalPurchaseUnit    `json:"purchase_units"`
	Links         []PayPalLink            `json:"links"`
	Payer         *PayPalPayer            `json:"payer,omitempty"`
}

type PayPalPurchaseUnit struct {
	ReferenceID string        `json:"reference_id"`
	Amount      PayPalAmount  `json:"amount"`
	Description string        `json:"description"`
}

type PayPalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type PayPalPayer struct {
	EmailAddress string `json:"email_address"`
	PayerID      string `json:"payer_id"`
}

type PayPalAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewPayPalService(cfg *config.Config) *PayPalService {
	baseURL := "https://api-m.paypal.com"
	if cfg.PayPalMode == "sandbox" {
		baseURL = "https://api-m.sandbox.paypal.com"
	}

	return &PayPalService{
		config:  cfg,
		baseURL: baseURL,
	}
}

func (s *PayPalService) getAccessToken() (string, error) {
	url := fmt.Sprintf("%s/v1/oauth2/token", s.baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.config.PayPalClientID, s.config.PayPalClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp PayPalAccessToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (s *PayPalService) CreateOrder(amount, currency, description, referenceID string) (*PayPalOrder, error) {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": referenceID,
				"description":  description,
				"amount": map[string]string{
					"currency_code": currency,
					"value":         amount,
				},
			},
		},
		"application_context": map[string]string{
			"return_url": fmt.Sprintf("%s/payments/paypal/success", s.config.AppURL),
			"cancel_url": fmt.Sprintf("%s/payments/paypal/cancel", s.config.AppURL),
		},
	}

	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v2/checkout/orders", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var order PayPalOrder
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *PayPalService) CaptureOrder(orderID string) (*PayPalOrder, error) {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", s.baseURL, orderID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var order PayPalOrder
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *PayPalService) GetOrder(orderID string) (*PayPalOrder, error) {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s", s.baseURL, orderID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var order PayPalOrder
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
