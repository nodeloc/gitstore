#!/bin/bash

set -e

echo "🚀 GitStore 生产部署脚本 (使用 Nginx Proxy Manager)"
echo "========================================"

# 配置变量
DOMAIN=${DOMAIN:-"discourseplugin.com"}
PROJECT_DIR="/opt/gitstore"

# 检查 root 权限
if [ "$EUID" -ne 0 ]; then 
    echo "❌ 请使用 root 或 sudo 运行此脚本"
    exit 1
fi

echo ""
echo "📋 部署配置:"
echo "   域名: $DOMAIN"
echo "   项目目录: $PROJECT_DIR"
echo ""

read -p "确认部署配置正确吗? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 部署已取消"
    exit 1
fi

# 1. 停止旧服务
echo ""
echo "1️⃣  停止旧服务..."
cd $PROJECT_DIR
docker-compose -f docker-compose.prod.yml down || true

# 2. 更新代码
echo ""
echo "2️⃣  更新代码..."
git pull origin main || echo "跳过 git pull"

# 3. 构建前端
echo ""
echo "3️⃣  构建前端..."
cd $PROJECT_DIR/frontend
npm install
npm run build

# 4. 启动所有服务
echo ""
echo "4️⃣  启动所有服务..."
cd $PROJECT_DIR
docker-compose -f docker-compose.prod.yml up -d --build

# 5. 等待服务启动
echo ""
echo "5️⃣  等待服务启动..."
sleep 15

# 6. 健康检查
echo ""
echo "6️⃣  健康检查..."
for i in {1..30}; do
    if docker exec gitstore-backend wget --quiet --tries=1 --spider http://localhost:8080/api/health 2>/dev/null; then
        echo "   ✅ 后端服务启动成功"
        break
    fi
    echo "   ⏳ 等待后端启动... ($i/30)"
    sleep 2
done

# 7. 显示状态和配置说明
echo ""
echo "✅ 部署完成!"
echo ""
echo "📊 服务状态:"
docker-compose -f docker-compose.prod.yml ps
echo ""
echo "🌐 访问地址:"
echo "   Nginx Proxy Manager 管理界面: http://YOUR_SERVER_IP:81"
echo "   默认登录: admin@example.com / changeme"
echo ""
echo "📝 后续配置步骤:"
echo "   1. 访问 http://YOUR_SERVER_IP:81 登录 NPM"
echo "   2. 修改默认密码"
echo "   3. 添加代理主机 (Proxy Hosts):"
echo "      - Domain: $DOMAIN"
echo "      - Forward Hostname/IP: gitstore-frontend"
echo "      - Forward Port: 80"
echo "      - 启用 SSL (Let's Encrypt)"
echo ""
echo "📝 查看日志:"
echo "   后端: docker-compose -f docker-compose.prod.yml logs -f backend"
echo "   前端: docker-compose -f docker-compose.prod.yml logs -f frontend"
echo "   数据库: docker-compose -f docker-compose.prod.yml logs -f postgres"
echo "   NPM: docker-compose -f docker-compose.prod.yml logs -f nginx-proxy-manager"
echo ""
echo "🔧 直接访问测试 (通过内部网络):"
echo "   docker exec gitstore-backend wget -O- http://localhost:8080/api/health"
