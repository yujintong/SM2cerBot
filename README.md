<div align="center">
  <p>国密ACME客户端</p>
</div>

# SM2cerBot

这是一个基于《自动化证书管理规范》密码行业标准草案V6版本开发的开源ACME客户端软件，默认自动对接零信技术提供的国密ACME公共服务，自动申请证书、取回证书、部署证书和续期证书。双算法SSL证书都是90天有效期证书，国际SSL证书全球信任，所有浏览器信任；国密SSL证书国密合规，所有国密浏览器都信任，推荐使用完全免费的零信浏览器，优先采用国密算法实现HTTPS加密。

本项目基于 [LEGO](https://github.com/go-acme/lego) 改造

## 编译

- make build
- 编译后文件在dist目录下

## 运行准备条件（nginx）
- 
- 1.在nginx的服务上添加location, 包括（80和443端口）
-  'location /.well-known/pki-validation/ {
-        root /opt/acme/webroot ;       # http.webroot 指定路径，自定义
-    }'
-  配置后重启nginx
- 
- 
-  2.添加run-hook脚本文件(hook.sh)
-  '
-  #!/bin/bash
-  sudo nginx -s reload
-  echo "hook file run success"
-  '
-  修改文件属性  chmod +x hook.sh
- 
- 
-  3.kid是acme的外部账号，请在网站 (https://zotrus.com) 上注册
- 
- 
-  4.hmac是acme的外部账号token，请在网站 (https://zotrus.com) 用户中心查看
- 

## 运行 
 - RUN
 - ./SM2cerBot --email=email@domain.com --eab  --kid=email@domain.com   --hmac=token    --domains=youdomain.com    --path=/opt/acme   --http    --http.webroot=/opt/acme/webroot  --run-hook=/workspace/acme/hook.sh  run

## 证书签发后配置
 - 修改证书引用地址(nginx), 证书签发后存放在path目录下 例如 /opt/acme/certificates
 - 
 -   ssl_certificate          /opt/acme/certificates/youdomain_ecc.crt;
 -   ssl_certificate_key      /opt/acme/certificates/youdomain_ecc_ecc.key;
 -   ssl_sign_certificate     /opt/acme/certificates/youdomain_ecc_sign.crt;
 -   ssl_sign_certificate_key /opt/acme/certificates/youdomain_ecc_sign.key;
 -   ssl_enc_certificate      /opt/acme/certificates/youdomain_ecc_enc.crt;
 -   ssl_enc_certificate_key  /opt/acme/certificates/sm2acmetest.cersign.cn_enc.key;
 - 
 - 
 
 ## 续订
 -  
 - Renew
 - ./SM2cerBot --email=email@domain.com --eab  --kid=email@domain.com   --hmac=token    --domains=youdomain.com    --path=/opt/acme    --http --days=80 --http.webroot=/opt/acme/webroot   --renew-hook=/workspace/acme/hook.sh renew
