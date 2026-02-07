# Nginx Proxy Manager 配置指南

## 1. 首次登录

部署完成后，访问：
```
http://YOUR_SERVER_IP:81
```

默认账号：
- Email: `admin@example.com`
- Password: `changeme`

**⚠️ 登录后立即修改密码！**

---

## 2. 添加代理主机（Proxy Host）

### 步骤：

1. **点击 "Proxy Hosts" → "Add Proxy Host"**

2. **填写 Details 选项卡：**
   - **Domain Names**: `discourseplugin.com`（或你的域名）
   - **Scheme**: `http`
   - **Forward Hostname / IP**: `git-store-frontend`
   - **Forward Port**: `80`
   - **Cache Assets**: ✅ 勾选
   - **Block Common Exploits**: ✅ 勾选
   - **Websockets Support**: ✅ 勾选

3. **配置 SSL 选项卡：**
   - **SSL Certificate**: 选择 "Request a new SSL Certificate"
   - **Force SSL**: ✅ 勾选
   - **HTTP/2 Support**: ✅ 勾选
   - **HSTS Enabled**: ✅ 勾选
   - **Email**: 填入你的邮箱（用于 Let's Encrypt）
   - **I Agree to the Let's Encrypt Terms of Service**: ✅ 勾选

4. **点击 Save**

---

## 3. 添加多个域名（可选）

如果需要同时支持 `www.discourseplugin.com`：

1. 编辑现有 Proxy Host
2. 在 **Domain Names** 添加第二行：
   ```
   discourseplugin.com
   www.discourseplugin.com
   ```

---

## 4. 高级配置（Custom Nginx Configuration）

如果需要自定义配置，在 **Advanced** 选项卡添加：

```nginx
# 上传文件大小限制
client_max_body_size 50M;

# 自定义头部
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;

# 超时设置
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;
```

---

## 5. 验证配置

访问你的域名：
```
https://discourseplugin.com
```

检查：
- ✅ HTTPS 正常工作
- ✅ 自动跳转到 HTTPS
- ✅ SSL 证书有效
- ✅ API 请求正常 (`/api/health`)

---

## 6. 常见问题

### SSL 证书申请失败
- 确保域名 DNS 已正确解析到服务器 IP
- 确保 80 和 443 端口已开放
- 等待 DNS 传播（最多 24 小时）

### 502 Bad Gateway
- 检查后端服务是否正常：`docker ps`
- 查看日志：`docker logs git-store-backend`
- 确认 `Forward Hostname/IP` 是容器名称，不是 IP

### 无法访问管理界面
- 确保端口 81 已开放：`sudo ufw allow 81/tcp`
- 重启 NPM：`docker restart nginx-proxy-manager`

---

## 7. 备份配置

NPM 的所有配置存储在 Docker Volume 中：

```bash
# 备份
docker run --rm -v npm_data:/data -v $(pwd):/backup alpine tar czf /backup/npm-backup.tar.gz /data

# 恢复
docker run --rm -v npm_data:/data -v $(pwd):/backup alpine sh -c "cd / && tar xzf /backup/npm-backup.tar.gz"
```

---

## 8. 监控和日志

```bash
# 查看 NPM 日志
docker logs -f nginx-proxy-manager

# 查看访问日志
docker exec nginx-proxy-manager cat /data/logs/proxy-host-*.log

# 查看错误日志
docker exec nginx-proxy-manager cat /data/logs/error.log
```

---

## 9. 性能优化建议

在 NPM 的 **Advanced** 配置中添加：

```nginx
# 启用 Gzip
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/json;

# 浏览器缓存
location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

---

## 10. 更新 NPM

```bash
cd /opt/git-store
docker-compose -f docker-compose.prod.yml pull nginx-proxy-manager
docker-compose -f docker-compose.prod.yml up -d nginx-proxy-manager
```
