# Auth

[![GPL-3.0](https://img.shields.io/github/license/auto-novel/auth)](https://github.com/auto-novel/auth#license)

提供统一登录认证（SSO）服务。

## 部署

```bash
# 下载项目
mkdir auth
cd auth
mkdir sql
curl -sSL "https://raw.githubusercontent.com/auto-novel/auth/refs/heads/main/docker-compose.yml" -o "./docker-compose.yml"
curl -sSL "https://raw.githubusercontent.com/auto-novel/auth/refs/heads/main/sql/init.sql" -o "./sql/init.sql"

# 配置环境变量
echo "REFRESH_TOKEN_SECRET=$(pwgen -s 64 1)" >> .env
echo "ACCESS_TOKEN_SECRET=$(pwgen -s 64 1)" >> .env
echo "POSTGRES_PASSWORD=$(pwgen -s 64 1)" >> .env
echo "SMTP_MAIL=no-reply@example.com" >> .env
echo "SMTP_SERVER=smtp.example.com:465" >> .env
echo "SMTP_PASSWORD=" >> .env

# 启动服务
docker compose up -d
```

## 开发

### Api

```bash
make start_debug        # 启动 docker compose, debug 模式
make start_release      # 启动 docker compose, release 模式
make stop               # 关闭 docker compose

make generate           # 生成 sql 代码
make integration_test   # 运行集成测试
```