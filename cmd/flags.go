package cmd

import (
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/urfave/cli/v2"
	"software.sslmate.com/src/go-pkcs12"
)

// Flag names.
const (
	flgDomains                  = "domains"
	flgServer                   = "server"
	flgAcceptTOS                = "accept-tos"
	flgEmail                    = "email"
	flgDisableCommonName        = "disable-cn"
	flgCSR                      = "csr"
	flgEAB                      = "eab"
	flgKID                      = "kid"
	flgHMAC                     = "hmac"
	flgKeyType                  = "key-type"
	flgFilename                 = "filename"
	flgPath                     = "path"
	flgHTTP                     = "http"
	flgHTTPPort                 = "http.port"
	flgHTTPDelay                = "http.delay"
	flgHTTPProxyHeader          = "http.proxy-header"
	flgHTTPWebroot              = "http.webroot"
	flgHTTPMemcachedHost        = "http.memcached-host"
	flgHTTPS3Bucket             = "http.s3-bucket"
	flgTLS                      = "tls"
	flgTLSPort                  = "tls.port"
	flgTLSDelay                 = "tls.delay"
	flgDNS                      = "dns"
	flgDNSDisableCP             = "dns.disable-cp"
	flgDNSPropagationWait       = "dns.propagation-wait"
	flgDNSPropagationDisableANS = "dns.propagation-disable-ans"
	flgDNSPropagationRNS        = "dns.propagation-rns"
	flgDNSResolvers             = "dns.resolvers"
	flgHTTPTimeout              = "http-timeout"
	flgTLSSkipVerify            = "tls-skip-verify"
	flgDNSTimeout               = "dns-timeout"
	flgPEM                      = "pem"
	flgPFX                      = "pfx"
	flgPFXPass                  = "pfx.pass"
	flgPFXFormat                = "pfx.format"
	flgCertTimeout              = "cert.timeout"
	flgOverallRequestLimit      = "overall-request-limit"
	flgUserAgent                = "user-agent"
)

const (
	envEAB         = "LEGO_EAB"
	envEABHMAC     = "LEGO_EAB_HMAC"
	envEABKID      = "LEGO_EAB_KID"
	envEmail       = "LEGO_EMAIL"
	envPath        = "LEGO_PATH"
	envPFX         = "LEGO_PFX"
	envPFXFormat   = "LEGO_PFX_FORMAT"
	envPFXPassword = "LEGO_PFX_PASSWORD"
	envServer      = "LEGO_SERVER"
)

func CreateFlags(defaultPath string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    flgDomains,
			Aliases: []string{"d"},
			Usage:   "Add a domain to the process. Can be specified multiple times.",
		},
		&cli.StringFlag{
			Name:    flgServer,
			Aliases: []string{"s"},
			EnvVars: []string{envServer},
			Usage:   "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
			Value:   lego.LEDirectoryProduction,
		},
		&cli.BoolFlag{
			Name:    flgAcceptTOS,
			Aliases: []string{"a"},
			Usage:   "By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.",
		},
		&cli.StringFlag{
			Name:    flgEmail,
			Aliases: []string{"m"},
			EnvVars: []string{envEmail},
			Usage:   "Email used for registration and recovery contact.",
		},
		&cli.BoolFlag{
			Name:  flgDisableCommonName,
			Usage: "Disable the use of the common name in the CSR.",
		},
		&cli.StringFlag{
			Name:    flgCSR,
			Aliases: []string{"c"},
			Usage:   "Certificate signing request filename, if an external CSR is to be used.",
		},
		&cli.BoolFlag{
			Name:    flgEAB,
			EnvVars: []string{envEAB},
			Usage:   "Use External Account Binding for account registration. Requires --kid and --hmac.",
		},
		&cli.StringFlag{
			Name:    flgKID,
			EnvVars: []string{envEABKID},
			Usage:   "Key identifier from External CA. Used for External Account Binding.",
		},
		&cli.StringFlag{
			Name:    flgHMAC,
			EnvVars: []string{envEABHMAC},
			Usage:   "MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.",
		},
		&cli.StringFlag{
			Name:    flgKeyType,
			Aliases: []string{"k"},
			Value:   "ec256",
			Usage:   "Key type to use for private keys. Supported: rsa2048, rsa3072, rsa4096, rsa8192, ec256, ec384.",
		},
		&cli.StringFlag{
			Name:  flgFilename,
			Usage: "(deprecated) Filename of the generated certificate.",
		},
		&cli.StringFlag{
			Name:    flgPath,
			EnvVars: []string{envPath},
			Usage:   "Directory to use for storing the data.",
			Value:   defaultPath,
		},
		&cli.BoolFlag{
			Name:  flgHTTP,
			Usage: "Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.",
		},
		&cli.DurationFlag{
			Name:  flgHTTPDelay,
			Usage: "Delay between the starts of the HTTP server (use for HTTP-01 based challenges) and the validation of the challenge.",
			Value: 0,
		},
		&cli.StringFlag{
			Name: flgHTTPWebroot,
			Usage: "Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file." +
				" This disables the built-in server and expects the given directory to be publicly served with access to .well-known/acme-challenge",
		},

		&cli.IntFlag{
			Name:  flgHTTPTimeout,
			Usage: "Set the HTTP timeout value to a specific value in seconds.",
		},
		&cli.IntFlag{
			Name:  flgDNSTimeout,
			Usage: "Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries.",
			Value: 10,
		},
		&cli.BoolFlag{
			Name:  flgPEM,
			Usage: "Generate an additional .pem (base64) file by concatenating the .key and .crt files together.",
		},
		&cli.BoolFlag{
			Name:    flgPFX,
			Usage:   "Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.",
			EnvVars: []string{envPFX},
		},
		&cli.StringFlag{
			Name:    flgPFXPass,
			Usage:   "The password used to encrypt the .pfx (PCKS#12) file.",
			Value:   pkcs12.DefaultPassword,
			EnvVars: []string{envPFXPassword},
		},
		&cli.StringFlag{
			Name:    flgPFXFormat,
			Usage:   "The encoding format to use when encrypting the .pfx (PCKS#12) file. Supported: RC2, DES, SHA256.",
			Value:   "RC2",
			EnvVars: []string{envPFXFormat},
		},
		&cli.IntFlag{
			Name:  flgCertTimeout,
			Usage: "Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates.",
			Value: 30,
		},
		&cli.IntFlag{
			Name:  flgOverallRequestLimit,
			Usage: "ACME overall requests limit.",
			Value: certificate.DefaultOverallRequestLimit,
		},
	}
}

func getTime(ctx *cli.Context, name string) time.Time {
	value := ctx.Timestamp(name)
	if value == nil {
		return time.Time{}
	}

	return *value
}
