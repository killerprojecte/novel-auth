# Auth

[![GPL-3.0](https://img.shields.io/github/license/auto-novel/auth)](https://github.com/auto-novel/auth#license)

提供统一登录认证（SSO）服务。

## 部署 (WIP)

```bash
echo "REFRESH_TOKEN_SECRET=$(pwgen -s 64 1)" >> .env
echo "ACCESS_TOKEN_SECRET=$(pwgen -s 64 1)" >> .env
echo "POSTGRES_PASSWORD=$(pwgen -s 64 1)" >> .env
echo "MAILGUN_DOMAIN=verify.fishhawk.top" >> .env
echo "MAILGUN_APIKEY=" >> .env
```

## 开发

```bash
docker compose -f docker-compose.yml -f docker-compose.debug.yml up -d
docker compose -f docker-compose.yml -f docker-compose.debug.yml up -d --build # 重新构建镜像
```