package acme

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
	gmx509 "github.com/tjfoc/gmsm/x509"
)

// KeyType 密钥类型
type KeyType int

const (
	RSA KeyType = iota
	ECC
	SM2
)

// CertRequest 证书请求参数
type CertRequest struct {
	CommonName         string
	Organization       []string
	OrganizationalUnit []string
	Country            []string
	Province           []string
	Locality           []string
	EmailAddresses     []string
	DNSNames           []string
	IPAddresses        []string
}

// GenerateKeyPair 生成密钥对
func GenerateKeyPair(keyType KeyType, bits int) (interface{}, error) {
	switch keyType {
	case RSA:
		return rsa.GenerateKey(rand.Reader, bits)
	case ECC:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case SM2:
		return sm2.GenerateKey(rand.Reader)
	default:
		return nil, fmt.Errorf("unsupported key type")
	}
}

// CreateCSR 创建证书签名请求
func CreateCSR(keyType KeyType, keyBits int, req CertRequest) (csrBaseUrl string, keyPrivate []byte, err error) {
	// 生成密钥对
	privateKey, err := GenerateKeyPair(keyType, keyBits)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	var csrDER []byte

	// 根据密钥类型创建 CSR
	switch keyType {
	case RSA, ECC:
		// 创建证书请求模板
		template := x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:         req.CommonName,
				Organization:       req.Organization,
				OrganizationalUnit: req.OrganizationalUnit,
				Country:            req.Country,
				Province:           req.Province,
				Locality:           req.Locality,
			},
			EmailAddresses: req.EmailAddresses,
			DNSNames:       req.DNSNames,

			// SignatureAlgorithm: ,
		}
		if keyType == RSA {
			template.SignatureAlgorithm = x509.SHA256WithRSA
		} else {
			template.SignatureAlgorithm = x509.ECDSAWithSHA256
		}

		// 使用标准库处理 RSA 和 ECC
		csrDER, err = x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create CSR: %v", err)
		}
	case SM2:
		// 创建证书请求模板
		template := gmx509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:         req.CommonName,
				Organization:       req.Organization,
				OrganizationalUnit: req.OrganizationalUnit,
				Country:            req.Country,
				Province:           req.Province,
				Locality:           req.Locality,
			},
			EmailAddresses:     req.EmailAddresses,
			DNSNames:           req.DNSNames,
			SignatureAlgorithm: gmx509.SM2WithSM3,
		}

		// 使用国密库处理 SM2
		sm2PrivateKey := privateKey.(*sm2.PrivateKey)
		csrDER, err = gmx509.CreateCertificateRequest(rand.Reader, &template, sm2PrivateKey)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create SM2 CSR: %v", err)
		}
	}

	// 编码 CSR 为 PEM 格式
	// csrbase := pem.EncodeToMemory(&pem.Block{
	// 	Type:  "CERTIFICATE REQUEST",
	// 	Bytes: csrDER,
	// })

	// fmt.Printf(string(csrbase))

	csrPEM := csrToBase64URL(csrDER)

	// 编码私钥为 PEM 格式
	privateKeyBytes, err := encodePrivateKey(privateKey, keyType)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encode private key: %v", err)
	}

	return csrPEM, privateKeyBytes, nil
}

func csrToBase64URL(csrDER []byte) string {
	// 直接对 DER 数据进行 Base64URL 编码
	return base64.URLEncoding.EncodeToString(csrDER)
}

// encodePrivateKey 编码私钥为 PEM 格式
func encodePrivateKey(privateKey interface{}, keyType KeyType) ([]byte, error) {
	var pemType string
	var keyBytes []byte
	var err error

	switch keyType {
	case RSA:
		pemType = "RSA PRIVATE KEY"
		rsaKey := privateKey.(*rsa.PrivateKey)
		keyBytes = x509.MarshalPKCS1PrivateKey(rsaKey)
	case ECC:
		pemType = "EC PRIVATE KEY"
		ecKey := privateKey.(*ecdsa.PrivateKey)
		keyBytes, err = x509.MarshalECPrivateKey(ecKey)
		if err != nil {
			return nil, err
		}
	case SM2:
		pemType = "SM2 PRIVATE KEY"
		sm2Key := privateKey.(*sm2.PrivateKey)
		keyBytes, err = gmx509.MarshalSm2UnecryptedPrivateKey(sm2Key)
		if err != nil {
			return nil, err
		}
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  pemType,
		Bytes: keyBytes,
	})

	return privateKeyPEM, nil
}
