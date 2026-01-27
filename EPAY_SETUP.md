# 易支付接入说明

本项目已集成易支付（码上付）接口，用于处理支付宝支付。

## 接口文档

- **API文档**: https://p.ma3fu.com/doc/pay_create.html
- **签名规则**: https://p.ma3fu.com/doc/sign_note.html
- **商户后台**: https://p.ma3fu.com/

## 配置说明

### 1. 环境变量配置

在 `.env` 文件中配置以下参数：

```bash
# 易支付配置
EPAY_PID=your_epay_pid                          # 商户ID（从商户后台获取）
EPAY_API_URL=https://p.ma3fu.com/api/pay/create # API地址（默认）
EPAY_PRIVATE_KEY=your_merchant_key              # 商户密钥（MD5或RSA格式）
EPAY_PUBLIC_KEY=platform_public_key             # 平台公钥（仅RSA签名需要）
EPAY_NOTIFY_URL=http://localhost:8080/api/webhooks/alipay  # 异步通知地址
```

**签名方式自动检测**：
- 如果 `EPAY_PRIVATE_KEY` 长度小于100字符，自动使用 **MD5签名**
- 如果 `EPAY_PRIVATE_KEY` 长度大于等于100字符，自动使用 **RSA签名**

### 2. 获取配置参数

1. 注册易支付商户账号: https://p.ma3fu.com/
2. 登录商户后台，获取商户ID（PID）
3. 获取密钥：
   - **MD5签名**：直接从商户后台复制密钥字符串（如：`5BBbibRC85RFAut29Ibmz9nCiZA9bzur`）
   - **RSA签名**：生成RSA密钥对，商户私钥用于签名，平台公钥用于验签

### 3. 密钥格式

#### MD5密钥格式（推荐用于测试）：
```
5BBbibRC85RFAut29Ibmz9nCiZA9bzur
```
直接使用从商户后台获取的密钥字符串，无需额外格式。

#### RSA私钥格式（用于生产环境）：
```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
...
-----END RSA PRIVATE KEY-----
```

#### RSA平台公钥格式：
```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
...
-----END PUBLIC KEY-----
```

**注意**: RSA密钥可以不包含 BEGIN/END 标签，代码会自动添加。

## 接口说明

### 1. 创建支付订单

**接口**: `POST /api/payments/alipay/create`

**请求参数**:
```json
{
  "order_id": "订单UUID"
}
```

**响应示例**:
```json
{
  "trade_no": "平台订单号",
  "pay_type": "jump",  // 支付类型：jump(跳转), qrcode(二维码), html(HTML代码)等
  "pay_info": "支付URL或二维码链接",
  "order_id": "订单UUID",
  "amount": 99.00
}
```

**支付类型说明**:
- `jump`: 直接跳转URL，前端使用 `window.location.href = pay_info` 跳转
- `qrcode`: 二维码链接，前端生成二维码展示
- `html`: HTML代码，前端渲染后自动跳转

### 2. 支付回调通知

**接口**: `POST /api/webhooks/alipay`

易支付会向此接口发送异步通知，包含以下参数：
- `trade_status`: 订单状态，`TRADE_SUCCESS` 或 `1` 表示支付成功
- `out_trade_no`: 商户订单号（订单UUID）
- `trade_no`: 平台订单号
- `money`: 支付金额
- `sign`: 签名字符串
- `sign_type`: 签名类型（RSA）

**处理流程**:
1. 验证签名
2. 检查订单状态，避免重复处理
3. 更新订单状态为已支付
4. 生成许可证（License）
5. 返回 `success` 字符串

## 前端集成示例

### 发起支付

```javascript
// 创建支付订单
const response = await fetch('/api/payments/alipay/create', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    order_id: orderId
  })
});

const result = await response.json();

// 根据支付类型处理
switch (result.pay_type) {
  case 'jump':
    // 直接跳转
    window.location.href = result.pay_info;
    break;
    
  case 'qrcode':
    // 显示二维码
    showQRCode(result.pay_info);
    break;
    
  case 'html':
    // 渲染HTML（会自动跳转）
    document.write(result.pay_info);
    break;
}
```

### 支付结果处理

用户支付完成后会跳转到 `return_url`（前端URL），前端需要：

1. 检查URL参数确认支付成功
2. 轮询订单状态或使用WebSocket实时更新
3. 显示支付成功页面

## 测试说明

### 1. 本地测试

需要配置公网回调地址，可以使用：
- [ngrok](https://ngrok.com/): `ngrok http 8080`
- [frp](https://github.com/fatedier/frp): 内网穿透工具

配置示例：
```bash
# 使用ngrok
ngrok http 8080
# 获得公网地址，如: https://abc123.ngrok.io

# 更新.env中的回调地址
EPAY_NOTIFY_URL=https://abc123.ngrok.io/api/webhooks/alipay
```

### 2. 查看日志

支付回调会记录详细日志：
```bash
# 成功日志
Payment successful - Order: xxx-xxx-xxx, TradeNo: 20240124xxx

# 失败日志  
Failed to verify Alipay signature: xxx
```

## 签名规则

### MD5签名方式

**签名步骤**（发送请求时）:
1. 获取所有非空参数，剔除 `sign` 和 `sign_type`
2. 按照键名ASCII码递增排序
3. 拼接成 `key1=value1&key2=value2` 格式
4. 在字符串末尾追加商户密钥：`待签名字符串 + 密钥`
5. 对结果字符串进行MD5哈希
6. 转为大写十六进制字符串

**验签步骤**（接收回调时）:
1. 按照相同方式生成待签名字符串
2. 对比计算出的签名与回调中的签名是否一致

示例：
```
参数: pid=1001&type=alipay&money=199.00&...
密钥: 5BBbibRC85RFAut29Ibmz9nCiZA9bzur
待签名: pid=1001&type=alipay&money=199.00&...5BBbibRC85RFAut29Ibmz9nCiZA9bzur
签名: MD5(待签名字符串) 转大写
```

### RSA签名方式

**签名步骤**（发送请求时）:

1. 获取所有非空参数，剔除 `sign` 和 `sign_type`
2. 按照键名ASCII码递增排序
3. 拼接成 `key1=value1&key2=value2` 格式
4. 使用商户私钥进行RSA签名（SHA256WithRSA）
5. Base64编码签名结果

### 验签步骤（接收回调时）:

1. 获取所有参数，剔除 `sign` 和 `sign_type`
2. 按照键名ASCII码递增排序
3. 拼接成待验签字符串
4. 使用平台公钥验证签名（SHA256WithRSA）

## 注意事项

1. **密钥安全**: 商户私钥必须妥善保管，不要提交到代码仓库
2. **回调验签**: 必须验证回调签名，防止伪造通知
3. **幂等性**: 支付回调可能重复发送，需要做好幂等性处理
4. **订单状态**: 检查 `trade_status` 字段，只有 `TRADE_SUCCESS` 或 `1` 才表示支付成功
5. **金额校验**: 回调中应该校验支付金额是否与订单金额一致（当前未实现，可根据需要添加）
6. **超时处理**: 订单应设置超时时间，过期未支付自动关闭

## 故障排查

### 1. 签名验证失败

**原因**:
- 商户私钥或平台公钥配置错误
- 密钥格式不正确
- 参数顺序错误

**解决方法**:
- 检查密钥是否完整，包括 BEGIN/END 标签
- 确认密钥是否匹配（商户后台-密钥管理）
- 查看日志中的待签名字符串

### 2. 回调未收到

**原因**:
- 回调地址不可访问
- 防火墙拦截
- 本地开发未配置内网穿透

**解决方法**:
- 确认 `EPAY_NOTIFY_URL` 是公网可访问地址
- 检查服务器防火墙设置
- 使用 ngrok/frp 进行内网穿透

### 3. 许可证生成失败

**原因**:
- 用户没有绑定GitHub账户
- 数据库外键约束

**解决方法**:
- 确保用户已登录并绑定GitHub账户
- 检查数据库日志
- 查看 `github_accounts` 表数据

## 相关文件

- **支付服务**: `internal/services/payment_alipay.go`
- **支付处理器**: `internal/handlers/handlers_all.go`
- **配置加载**: `internal/config/config.go`
- **环境配置**: `.env`
- **辅助函数**: `internal/utils/helpers.go`

## 后续优化建议

1. **金额校验**: 在回调中添加订单金额验证
2. **邮件通知**: 支付成功后发送邮件通知
3. **日志完善**: 记录更详细的支付流程日志
4. **监控告警**: 添加支付失败率监控
5. **退款功能**: 实现退款接口（如果易支付支持）
6. **订单超时**: 实现订单超时自动取消机制
