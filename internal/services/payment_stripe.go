package services

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
	"github.com/nodeloc/git-store/internal/config"
)

type StripeService struct {
	config *config.Config
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeSecretKey
	return &StripeService{
		config: cfg,
	}
}

// SetupWebhook 检查并创建Stripe Webhook端点
func (s *StripeService) SetupWebhook() error {
	if s.config.AppURL == "" {
		log.Println("[Stripe] APP_URL not configured, skipping webhook setup")
		return nil
	}

	webhookURL := s.config.AppURL + "/api/webhooks/stripe"
	
	// 检查是否已存在相同URL的Webhook
	params := &stripe.WebhookEndpointListParams{}
	iter := webhookendpoint.List(params)
	
	for iter.Next() {
		we := iter.WebhookEndpoint()
		if we.URL == webhookURL {
			log.Printf("[Stripe] Webhook already exists: %s (ID: %s)", webhookURL, we.ID)
			return nil
		}
	}
	
	if err := iter.Err(); err != nil {
		log.Printf("[Stripe] Failed to list webhook endpoints: %v", err)
		return err
	}
	
	// 创建新的Webhook端点
	createParams := &stripe.WebhookEndpointParams{
		URL: stripe.String(webhookURL),
		EnabledEvents: stripe.StringSlice([]string{
			"payment_intent.succeeded",
			"payment_intent.payment_failed",
		}),
		APIVersion: stripe.String("2023-10-16"),
		Description: stripe.String("Auto-created by gitstore"),
	}
	
	endpoint, err := webhookendpoint.New(createParams)
	if err != nil {
		log.Printf("[Stripe] Failed to create webhook endpoint: %v", err)
		return err
	}
	
	log.Printf("[Stripe] ✅ Created webhook endpoint:")
	log.Printf("[Stripe]   - URL: %s", endpoint.URL)
	log.Printf("[Stripe]   - ID: %s", endpoint.ID)
	log.Printf("[Stripe]   - Secret: %s", endpoint.Secret)
	log.Printf("[Stripe]   ⚠️  Please update STRIPE_WEBHOOK_SECRET in .env with the secret above")
	
	return nil
}


type PaymentIntentRequest struct {
	Amount      int64             // Amount in cents
	Currency    string            // e.g., "usd"
	Description string
	Metadata    map[string]string
}

func (s *StripeService) CreatePaymentIntent(req *PaymentIntentRequest) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(req.Amount),
		Currency:    stripe.String(req.Currency),
		Description: stripe.String(req.Description),
	}

	// Add metadata
	for key, value := range req.Metadata {
		params.AddMetadata(key, value)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

func (s *StripeService) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

func (s *StripeService) ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Confirm(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return pi, nil
}

func (s *StripeService) CancelPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Cancel(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	return pi, nil
}

func (s *StripeService) VerifyWebhookSignature(payload []byte, signature string) (stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		s.config.StripeWebhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		return event, fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	return event, nil
}
