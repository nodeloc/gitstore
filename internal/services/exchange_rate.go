package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"gorm.io/gorm"
)

// ExchangeRateService 汇率服务
type ExchangeRateService struct {
	db     *gorm.DB
	config *config.Config
}

// NewExchangeRateService 创建汇率服务
func NewExchangeRateService(db *gorm.DB, cfg *config.Config) *ExchangeRateService {
	return &ExchangeRateService{
		db:     db,
		config: cfg,
	}
}

// ExchangeRateAPIResponse exchangerate-api 响应结构
type ExchangeRateAPIResponse struct {
	Result            string             `json:"result"`
	Documentation     string             `json:"documentation"`
	TermsOfUse        string             `json:"terms_of_use"`
	TimeLastUpdateUTC string             `json:"time_last_update_utc"`
	TimeNextUpdateUTC string             `json:"time_next_update_utc"`
	BaseCode          string             `json:"base_code"`
	ConversionRates   map[string]float64 `json:"conversion_rates"`
}

// UpdateExchangeRates 更新汇率（从免费 API 获取）
func (s *ExchangeRateService) UpdateExchangeRates() error {
	log.Println("📊 开始更新汇率...")

	// 使用免费的 exchangerate-api.com API
	// USD 作为基准货币
	url := "https://open.er-api.com/v6/latest/USD"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("获取汇率失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("汇率 API 返回错误状态: %d", resp.StatusCode)
	}

	var apiResp ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("解析汇率响应失败: %v", err)
	}

	if apiResp.Result != "success" {
		return fmt.Errorf("汇率 API 返回失败: %s", apiResp.Result)
	}

	// 更新数据库中的汇率
	// 主要更新 USD -> CNY (用于支付宝)
	now := time.Now()
	targetCurrencies := []string{"CNY", "EUR", "GBP", "JPY", "HKD"}

	for _, toCurrency := range targetCurrencies {
		rate, exists := apiResp.ConversionRates[toCurrency]
		if !exists {
			log.Printf("⚠️  未找到 USD -> %s 的汇率", toCurrency)
			continue
		}

		exchangeRate := models.ExchangeRate{
			FromCurrency: "USD",
			ToCurrency:   toCurrency,
			Rate:         rate,
			Source:       "exchangerate-api",
			LastUpdated:  now,
		}

		// 使用 UPSERT 更新或插入
		err := s.db.
			Where("from_currency = ? AND to_currency = ?", "USD", toCurrency).
			Assign(map[string]interface{}{
				"rate":         rate,
				"last_updated": now,
			}).
			FirstOrCreate(&exchangeRate).Error

		if err != nil {
			log.Printf("❌ 更新 USD -> %s 汇率失败: %v", toCurrency, err)
		} else {
			log.Printf("✅ 更新汇率: 1 USD = %.4f %s", rate, toCurrency)
		}
	}

	log.Println("📊 汇率更新完成")
	return nil
}

// GetExchangeRate 获取指定货币对的汇率
func (s *ExchangeRateService) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	// 如果是相同货币，返回 1
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	var rate models.ExchangeRate
	err := s.db.
		Where("from_currency = ? AND to_currency = ?", fromCurrency, toCurrency).
		First(&rate).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("未找到 %s -> %s 的汇率", fromCurrency, toCurrency)
		}
		return 0, err
	}

	// 检查汇率是否过期（超过24小时）
	if time.Since(rate.LastUpdated) > 24*time.Hour {
		log.Printf("⚠️  汇率数据已过期（最后更新: %s），建议更新", rate.LastUpdated.Format("2006-01-02 15:04:05"))
	}

	return rate.Rate, nil
}

// ConvertAmount 转换金额
func (s *ExchangeRateService) ConvertAmount(amount float64, fromCurrency, toCurrency string) (float64, error) {
	rate, err := s.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}

	convertedAmount := amount * rate
	return convertedAmount, nil
}
