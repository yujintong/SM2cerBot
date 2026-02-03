# SM2cerBot - 国密ACME客户端

这是一个基于《自动化证书管理规范》密码行业标准草案V6版本开发的开源ACME客户端软件，默认自动对接零信技术提供的国密ACME公共服务，自动申请证书、取回证书、部署证书和续期证书。

零信国密ACME公共服务网址:https://my.zotrus.com/cnacme

## 特性

- 双算法SSL证书都是90天有效期证书
- 国际SSL证书全球信任，所有浏览器信任
- 国密SSL证书国密合规，所有国密浏览器都信任
- 推荐使用完全免费的零信浏览器，优先采用国密算法实现HTTPS加密

## 项目来源

本项目基于 [LEGO](https://github.com/go-acme/lego) 改造

## 编译

```bash
make build
```

编译后的文件在 `dist` 目录下

## 运行准备条件

### Nginx 配置

在 Nginx 的服务上添加 location（包括 80 和 443 端口）：

```nginx
location /.well-known/pki-validation/ {
    root /opt/acme/webroot;  # http.webroot 指定路径，自定义
}
```

配置后重启 Nginx。

### Hook 脚本

添加 `run-hook` 脚本文件 (`hook.sh`)：

```bash
#!/bin/bash
sudo nginx -s reload
echo "hook file run success"
```

修改文件属性：

```bash
chmod +x hook.sh
```

### ACME 账户配置

- `kid` 是 ACME 的外部账号，是您在网站 [ZoTrus](https://zotrus.com) 上注册的邮箱
- `hmac` 是 ACME 的外部账号 token，您可在网站 [ZoTrus](https://zotrus.com) 用户中心获取

## 运行

### 申请证书

```bash
./SM2cerBot \
  --email=email@domain.com \
  --eab \
  --kid=email@domain.com \
  --hmac=token \
  --domains=youdomain.com \
  --path=/opt/acme \
  --http \
  --http.webroot=/opt/acme/webroot \
  --run-hook=/workspace/acme/hook.sh \
  run
```

### 续订证书

```bash
./SM2cerBot \
  --email=email@domain.com \
  --eab \
  --kid=email@domain.com \
  --hmac=token \
  --domains=youdomain.com \
  --path=/opt/acme \
  --http \
  --days=80 \
  --http.webroot=/opt/acme/webroot \
  --renew-hook=/workspace/acme/hook.sh \
  renew
```
