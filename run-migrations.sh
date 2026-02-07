#!/bin/bash

set -e

echo "🔄 执行数据库迁移..."

# 配置变量
CONTAINER_NAME="git-store-db-prod"
DB_NAME=${DB_NAME:-git_store}
DB_USER=${DB_USER:-postgres}

# 检查容器是否运行
if ! docker ps | grep -q $CONTAINER_NAME; then
    echo "❌ 数据库容器未运行，请先启动服务"
    exit 1
fi

echo "📋 执行迁移文件..."

# 按顺序执行所有迁移文件
for sql_file in ./migrations/*.sql; do
    if [ -f "$sql_file" ]; then
        echo "   执行: $(basename $sql_file)"
        docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < "$sql_file" 2>&1 | grep -v "ERROR.*already exists\|ERROR.*duplicate key" || true
    fi
done

echo ""
echo "✅ 迁移执行完成！"
echo ""
echo "🔍 验证系统设置:"
docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "SELECT key, value FROM system_settings WHERE key IN ('site_name', 'site_subtitle', 'logo_url');"
