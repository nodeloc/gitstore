package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/nodeloc/git-store/internal/config"
)

type AlipayService struct {
	config     *config.Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	gatewayURL string
}

type AlipayTradeRequest struct {
	OutTradeNo  string  // 商户订单号
	TotalAmount float64 // 订单总金额
	Subject     string  // 订单标题
	Body        string  // 订单描述
	ProductCode string  // 产品码
}

func NewAlipayService(cfg *config.Config) (*AlipayService, error) {
	privateKey, err := parsePrivateKey(cfg.AlipayPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse alipay private key: %w", err)
	}

	publicKey, err := parsePublicKey(cfg.AlipayPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse alipay public key: %w", err)
	}

	return &AlipayService{
		config:     cfg,
		privateKey: privateKey,
		publicKey:  publicKey,
		gatewayURL: "https://openapi.alipay.com/gateway.do",
	}, nil
}

func parsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	privateKeyStr = strings.TrimSpace(privateKeyStr)

	// If the key doesn't have PEM headers, add them
	if !strings.HasPrefix(privateKeyStr, "-----BEGIN") {
		privateKeyStr = "-----BEGIN RSA PRIVATE KEY-----\n" + privateKeyStr + "\n-----END RSA PRIVATE KEY-----"
	}

	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		privateKey = key.(*rsa.PrivateKey)
	}

	return privateKey, nil
}

func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	publicKeyStr = strings.TrimSpace(publicKeyStr)

	// If the key doesn't have PEM headers, add them
	if !strings.HasPrefix(publicKeyStr, "-----BEGIN") {
		publicKeyStr = "-----BEGIN PUBLIC KEY-----\n" + publicKeyStr + "\n-----END PUBLIC KEY-----"
	}

	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := pubInterface.(*rsa.PublicKey)
	return publicKey, nil
}

func (s *AlipayService) sign(content string) (string, error) {
	hashed := crypto.SHA256.New()
	hashed.Write([]byte(content))

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (s *AlipayService) verify(content, sign string) error {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	hashed := crypto.SHA256.New()
	hashed.Write([]byte(content))

	return rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hashed.Sum(nil), signBytes)
}

func (s *AlipayService) buildRequestParams(bizContent map[string]interface{}, method string) (url.Values, error) {
	params := url.Values{}
	params.Set("app_id", s.config.AlipayAppID)
	params.Set("method", method)
	params.Set("format", "JSON")
	params.Set("charset", "utf-8")
	params.Set("sign_type", "RSA2")
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("notify_url", s.config.AlipayNotifyURL)

	// Convert bizContent to JSON string
	bizContentJSON, err := jsonMarshal(bizContent)
	if err != nil {
		return nil, err
	}
	params.Set("biz_content", string(bizContentJSON))

	// Generate signature
	signContent := s.buildSignContent(params)
	sign, err := s.sign(signContent)
	if err != nil {
		return nil, err
	}
	params.Set("sign", sign)

	return params, nil
}

func (s *AlipayService) buildSignContent(params url.Values) string {
	var keys []string
	for key := range params {
		if key != "sign" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signParts []string
	for _, key := range keys {
		value := params.Get(key)
		if value != "" {
			signParts = append(signParts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(signParts, "&")
}

// TradePagePay creates a web payment request (电脑网站支付)
func (s *AlipayService) TradePagePay(req *AlipayTradeRequest) (string, error) {
	bizContent := map[string]interface{}{
		"out_trade_no": req.OutTradeNo,
		"total_amount": fmt.Sprintf("%.2f", req.TotalAmount),
		"subject":      req.Subject,
		"body":         req.Body,
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}

	params, err := s.buildRequestParams(bizContent, "alipay.trade.page.pay")
	if err != nil {
		return "", err
	}

	// Return the complete payment URL
	return fmt.Sprintf("%s?%s", s.gatewayURL, params.Encode()), nil
}

// TradeWapPay creates a mobile payment request (手机网站支付)
func (s *AlipayService) TradeWapPay(req *AlipayTradeRequest) (string, error) {
	bizContent := map[string]interface{}{
		"out_trade_no": req.OutTradeNo,
		"total_amount": fmt.Sprintf("%.2f", req.TotalAmount),
		"subject":      req.Subject,
		"body":         req.Body,
		"product_code": "QUICK_WAP_WAY",
	}

	params, err := s.buildRequestParams(bizContent, "alipay.trade.wap.pay")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?%s", s.gatewayURL, params.Encode()), nil
}

// VerifyNotify verifies the Alipay callback notification
func (s *AlipayService) VerifyNotify(params url.Values) error {
	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")

	signContent := s.buildSignContent(params)
	return s.verify(signContent, sign)
}

// Helper function to marshal JSON
func jsonMarshal(v interface{}) ([]byte, error) {
	// Use encoding/json
	return []byte(fmt.Sprintf("%v", v)), nil
}
