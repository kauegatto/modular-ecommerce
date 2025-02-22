package adapters

import (
	"bytes"
	"context"
	"ecommerce/Payment/Domain/models"
	"ecommerce/Payment/Domain/ports"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ERedeConfig struct {
	PV      string
	Token   string
	BaseURL string
	Timeout time.Duration
}

type ERedeProcessor struct {
	config ERedeConfig
	client *http.Client
}

func NewERedeProcessor(config ERedeConfig) *ERedeProcessor {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &ERedeProcessor{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

type transactionRequest struct {
	Capture         bool   `json:"capture"`
	Kind            string `json:"kind"`
	Reference       string `json:"reference"`
	Amount          int    `json:"amount"`
	CardHolder      string `json:"cardHolderName"`
	CardNumber      string `json:"cardNumber"`
	ExpirationMonth int    `json:"expirationMonth"`
	ExpirationYear  int    `json:"expirationYear"`
	SecurityCode    string `json:"securityCode"`
	SoftDescriptor  string `json:"softDescriptor"`
}

type transactionResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Reference     string `json:"reference"`
	TID           string `json:"tid"`
	NSU           string `json:"nsu"`
	Authorization string `json:"authorization"`
}

func (p *ERedeProcessor) Capture(ctx context.Context, card *models.Card, payment *models.Payment) (ports.CaptureTransactionResponse, error) {
	expDate := card.ExpirationDate()
	expMonth, expYear := parseExpirationDate(expDate.String())

	var kind string
	if strings.EqualFold(string(payment.Kind), string(models.PaymentKindCredit)) {
		kind = "credit"
	} else {
		kind = "debit"
	}

	cardDTO := models.NewCardDTO(card)

	request := transactionRequest{
		Capture:         true,
		Kind:            kind,
		Reference:       payment.OrderId,
		Amount:          int(payment.TotalPrice), // Assuming Money is in cents
		CardHolder:      cardDTO.CardHolder,
		CardNumber:      cardDTO.CardNumber,
		ExpirationMonth: expMonth,
		ExpirationYear:  expYear,
		SecurityCode:    cardDTO.SecurityCode,
		SoftDescriptor:  "Ecommerce Do Kaue :)",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return ports.CaptureTransactionResponse{}, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/v1/transactions", p.config.BaseURL),
		bytes.NewBuffer(payload))
	if err != nil {
		return ports.CaptureTransactionResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.SetBasicAuth(p.config.PV, p.config.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApplicationId", p.config.PV)

	resp, err := p.client.Do(req)
	if err != nil {
		return ports.CaptureTransactionResponse{}, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	var response transactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ports.CaptureTransactionResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ports.CaptureTransactionResponse{}, fmt.Errorf("transaction failed: %s - %s",
			response.ReturnCode, response.ReturnMessage)
	}

	payment.ExternalIntegratorID = response.TID

	return ports.CaptureTransactionResponse{
		ExternalTransactionId: response.TID,
	}, nil
}

type RefundRequest struct {
	Amount int         `json:"amount"`
	URLs   []RefundURL `json:"urls,omitempty"`
}

type RefundURL struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

type RefundResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}

func (p *ERedeProcessor) RequestCancellation(ctx context.Context, eRedeTID string, amount int) error {
	url := fmt.Sprintf("%s/v1/transactions/%s/refunds", p.config.BaseURL, eRedeTID)

	reqBody := RefundRequest{
		Amount: amount,
		URLs: []RefundURL{
			{
				Kind: "callback",
				URL:  "https://testerde.free.beeceptor.com",
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal refund request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(p.config.PV, p.config.Token)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refund request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response RefundResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if response.ReturnCode != "359" {
		return fmt.Errorf("refund failed: %s - %s",
			response.ReturnCode, response.ReturnMessage)
	}

	return nil
}

func parseExpirationDate(expDate string) (month, year int) {
	var m, y int
	fmt.Sscanf(expDate, "%d/%d", &m, &y)
	return m, y
}

type ERedeError struct {
	Code    string
	Message string
}

func (e ERedeError) Error() string {
	return fmt.Sprintf("e-rede error: %s - %s", e.Code, e.Message)
}
