package services

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nodeloc/git-store/internal/config"
)

// EpayService 易支付服务
type AlipayService struct {
	config     *config.Config
	pid        string          // 商户ID
	key        string          // MD5密钥（用于MD5签名）
	privateKey *rsa.PrivateKey // 商户私钥（用于RSA签名）
	publicKey  *rsa.PublicKey  // 平台公钥（用于RSA验签）
	signType   string          // 签名类型：MD5或RSA
	apiURL     string          // API地址
	httpClient *http.Client    // HTTP客户端
}

// EpayCreateRequest 易支付统一下单请求
type AlipayTradeRequest struct {
	OutTradeNo  string  // 商户订单号
	TotalAmount float64 // 订单总金额
	Subject     string  // 商品名称
	Body        string  // 订单描述（可选）
	NotifyURL   string  // 异步通知地址
	ReturnURL   string  // 跳转通知地址
	ClientIP    string  // 用户IP地址
}

// EpayCreateResponse 易支付统一下单响应
type EpayCreateResponse struct {
	Code      int    `json:"code"`      // 返回状态码，0为成功，其它值为失败
	Msg       string `json:"msg"`       // 错误信息
	TradeNo   string `json:"trade_no"`  // 易支付订单号
	PayURL    string `json:"payurl"`    // 支付跳转URL（新版字段）
	QRCode    string `json:"qrcode"`    // 二维码链接
	URLScheme string `json:"urlscheme"` // 小程序跳转URL
	Sign      string `json:"sign"`      // 签名字符串
	SignType  string `json:"sign_type"` // 签名类型

	// 兼容旧字段
	PayType   string `json:"pay_type"` // 支付类型：jump/html/qrcode
	PayInfo   string `json:"pay_info"` // 支付信息（跳转URL、HTML代码或二维码链接）
	Timestamp string `json:"timestamp"`
}

func NewAlipayService(cfg *config.Config) (*AlipayService, error) {
	// 默认使用易支付API地址
	apiURL := cfg.AlipayAPIURL
	if apiURL == "" {
		apiURL = "https://p.ma3fu.com/mapi.php"
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	service := &AlipayService{
		config:     cfg,
		pid:        cfg.AlipayPID,
		apiURL:     apiURL,
		httpClient: httpClient,
	}

	// 检测密钥类型：如果密钥很短（<100字符），使用MD5签名；否则使用RSA签名
	privateKeyStr := strings.TrimSpace(cfg.AlipayPrivateKey)

	if len(privateKeyStr) < 100 {
		// MD5签名方式
		service.signType = "MD5"
		service.key = privateKeyStr
	} else {
		// RSA签名方式
		service.signType = "RSA"

		privateKey, err := parsePrivateKey(privateKeyStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse epay private key: %w", err)
		}
		service.privateKey = privateKey

		publicKey, err := parsePublicKey(cfg.AlipayPublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse epay public key: %w", err)
		}
		service.publicKey = publicKey
	}

	return service, nil
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
	if s.signType == "MD5" {
		// MD5签名：待签名字符串 + 密钥，结果小写
		signString := content + s.key
		h := md5.New()
		h.Write([]byte(signString))
		return hex.EncodeToString(h.Sum(nil)), nil // 小写
	}

	// RSA签名
	hashed := crypto.SHA256.New()
	hashed.Write([]byte(content))

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// verify 使用配置的签名类型验证（用于支付请求）
func (s *AlipayService) verify(content, sign string) error {
	return s.verifyWithType(content, sign, s.signType)
}

func (s *AlipayService) verifyWithType(content, sign, signType string) error {
	if signType == "MD5" {
		// MD5验签
		expectedSign, err := s.sign(content)
		if err != nil {
			return err
		}

		// 将签名转换为大写进行比较
		if strings.ToUpper(sign) != expectedSign {
			return fmt.Errorf("signature verification failed")
		}

		return nil
	}

	// RSA验签
	if s.publicKey == nil {
		return fmt.Errorf("RSA public key not configured")
	}

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	hashed := crypto.SHA256.New()
	hashed.Write([]byte(content))

	err = rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hashed.Sum(nil), signBytes)
	if err != nil {
		return fmt.Errorf("RSA signature verification failed: %v", err)
	}

	return nil
}

// buildSignContent 构建待签名字符串
// 按照ASCII码递增排序，剔除sign和sign_type字段
func (s *AlipayService) buildSignContent(params map[string]string) string {
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" && params[key] != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signParts []string
	for _, key := range keys {
		value := params[key]
		if value != "" {
			// 参数值不进行URL编码，直接拼接
			signParts = append(signParts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(signParts, "&")
}

// CreatePayment 创建易支付订单
func (s *AlipayService) CreatePayment(req *AlipayTradeRequest) (*EpayCreateResponse, error) {
	// 构建请求参数
	params := map[string]string{
		"pid":          s.pid,
		"method":       "web",                                    // 接口类型：web（根据device自动判断）
		"type":         "alipay",                                 // 支付方式：支付宝
		"out_trade_no": req.OutTradeNo,                           // 商户订单号
		"notify_url":   req.NotifyURL,                            // 异步通知地址
		"return_url":   req.ReturnURL,                            // 跳转通知地址
		"name":         req.Subject,                              // 商品名称
		"money":        fmt.Sprintf("%.2f", req.TotalAmount),     // 商品金额
		"clientip":     req.ClientIP,                             // 用户IP
		"device":       "pc",                                     // 设备类型：pc/mobile/wap
		"param":        req.Body,                                 // 业务扩展参数
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10), // 当前时间戳
		"sign_type":    s.signType,                               // 签名类型：MD5/RSA
	}

	// 生成签名
	signContent := s.buildSignContent(params)
	sign, err := s.sign(signContent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}
	params["sign"] = sign

	// 构建POST请求
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	// 构建完整的 API URL（在代码中拼接路径）
	createPaymentURL := s.apiURL + "/api/pay/create"

	resp, err := s.httpClient.PostForm(createPaymentURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var result EpayCreateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	// code=0 表示成功（易支付标准返回）
	if result.Code != 0 {
		return nil, fmt.Errorf("payment creation failed: %s", result.Msg)
	}

	return &result, nil
}

// VerifyNotify verifies the Alipay callback notification
func (s *AlipayService) VerifyNotify(params map[string]string) error {
	sign := params["sign"]
	if sign == "" {
		return fmt.Errorf("missing sign parameter")
	}

	// 易支付回调参数中可能包含 sign_type=RSA，但实际使用MD5签名
	// 强制使用配置的签名类型（MD5）进行验证
	signContent := s.buildSignContent(params)
	return s.verifyWithType(signContent, sign, s.signType)
}

// TradePagePay 创建网页支付（兼容旧接口）
func (s *AlipayService) TradePagePay(req *AlipayTradeRequest) (string, error) {
	if req.NotifyURL == "" {
		req.NotifyURL = s.config.AlipayNotifyURL
	}
	if req.ReturnURL == "" {
		req.ReturnURL = s.config.AppURL + "/payment/success"
	}
	if req.ClientIP == "" {
		req.ClientIP = "127.0.0.1"
	}

	result, err := s.CreatePayment(req)
	if err != nil {
		return "", err
	}

	// 根据pay_type返回不同的内容
	switch result.PayType {
	case "jump":
		// 直接跳转URL
		return result.PayInfo, nil
	case "html":
		// HTML代码，需要渲染
		return result.PayInfo, nil
	case "qrcode":
		// 二维码链接
		return result.PayInfo, nil
	default:
		return result.PayInfo, nil
	}
}

// TradeWapPay 创建移动支付（兼容旧接口）
func (s *AlipayService) TradeWapPay(req *AlipayTradeRequest) (string, error) {
	// 移动支付和网页支付使用相同的接口
	return s.TradePagePay(req)
}
