<div align="center">
  <p>Automatic Certificates and HTTPS for everyone.</p>
</div>

# Lego

Let's Encrypt client and ACME library written in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/go-acme/lego/v4.svg)](https://pkg.go.dev/github.com/go-acme/lego/v4)
[![Build Status](https://github.com//go-acme/lego/workflows/Main/badge.svg?branch=master)](https://github.com//go-acme/lego/actions)
[![Docker Pulls](https://img.shields.io/docker/pulls/goacme/lego.svg)](https://hub.docker.com/r/goacme/lego/)

lego is an independent, free, and open-source project, if you value it, consider [supporting it](https://donate.ldez.dev)! ❤️

Everybody thinks that the others will donate, but in the end, nobody does.

So if you think that lego is worth it, please consider [donating](https://donate.ldez.dev).

## Features

- ACME v2 [RFC 8555](https://www.rfc-editor.org/rfc/rfc8555.html)
  - Support [RFC 8737](https://www.rfc-editor.org/rfc/rfc8737.html): TLS Application‑Layer Protocol Negotiation (ALPN) Challenge Extension
  - Support [RFC 8738](https://www.rfc-editor.org/rfc/rfc8738.html): certificates for IP addresses
  - Support [RFC 9773](https://www.rfc-editor.org/rfc/rfc9773.html): Renewal Information (ARI) Extension
  - Support [draft-ietf-acme-profiles-00](https://datatracker.ietf.org/doc/draft-ietf-acme-profiles/): Profiles Extension
- Comes with about [170 DNS providers](https://go-acme.github.io/lego/dns)
- Register with CA
- Obtain certificates, both from scratch or with an existing CSR
- Renew certificates
- Revoke certificates
- Robust implementation of ACME challenges:
  - HTTP (http-01)
  - DNS (dns-01)
  - TLS (tls-alpn-01)
- SAN certificate support
- [CNAME support](https://letsencrypt.org/2019/10/09/onboarding-your-customers-with-lets-encrypt-and-acme.html) by default
- [Custom challenge solvers](https://go-acme.github.io/lego/usage/library/writing-a-challenge-solver/)
- Certificate bundling
- OCSP helper function


If your DNS provider is not supported, please open an [issue](https://github.com/go-acme/lego/issues/new?assignees=&labels=enhancement%2C+new-provider&template=new_dns_provider.yml).


## BUILD

- make build


## USING

 - RUN
 - SM2cerBot --email=