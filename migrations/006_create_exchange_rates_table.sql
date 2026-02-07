-- 创建汇率表
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(18, 8) NOT NULL,
    source VARCHAR(50) DEFAULT 'exchangerate-api',
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_currency_pair UNIQUE (from_currency, to_currency)
);

CREATE INDEX idx_exchange_rates_currencies ON exchange_rates(from_currency, to_currency);
CREATE INDEX idx_exchange_rates_updated ON exchange_rates(last_updated);

-- 插入初始汇率（美元到人民币，示例值）
INSERT INTO exchange_rates (from_currency, to_currency, rate) VALUES
    ('USD', 'CNY', 7.20)
ON CONFLICT (from_currency, to_currency) DO NOTHING;

COMMENT ON TABLE exchange_rates IS '货币汇率表';
COMMENT ON COLUMN exchange_rates.from_currency IS '源货币代码（如 USD）';
COMMENT ON COLUMN exchange_rates.to_currency IS '目标货币代码（如 CNY）';
COMMENT ON COLUMN exchange_rates.rate IS '汇率（1 源货币 = rate 目标货币）';
COMMENT ON COLUMN exchange_rates.source IS '汇率来源（如 exchangerate-api）';
