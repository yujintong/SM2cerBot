package cmd

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/acme/api"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
	"github.com/urfave/cli/v2"
)

// Flag names.
const (
	flgNoBundle                       = "no-bundle"
	flgMustStaple                     = "must-staple"
	flgNotBefore                      = "not-before"
	flgNotAfter                       = "not-after"
	flgPrivateKey                     = "private-key"
	flgPreferredChain                 = "preferred-chain"
	flgProfile                        = "profile"
	flgAlwaysDeactivateAuthorizations = "always-deactivate-authorizations"
	flgRunHook                        = "run-hook"
	flgRunHookTimeout                 = "run-hook-timeout"
)

func createRun() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Register an account, then create and install a certificate",
		Before: func(ctx *cli.Context) error {
			// we require either domains or csr, but not both
			hasDomains := len(ctx.StringSlice(flgDomains)) > 0

			hasCsr := ctx.String(flgCSR) != ""
			if hasDomains && hasCsr {
				log.Fatal("Please specify either --domains/-d or --csr/-c, but not both")
			}

			if !hasDomains && !hasCsr {
				log.Fatal("Please specify --domains/-d (or --csr/-c if you already have a CSR)")
			}

			return nil
		},
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flgRunHook,
				Usage: "Define a hook. The hook is executed when the certificates are effectively created.",
			},
			&cli.DurationFlag{
				Name:  flgRunHookTimeout,
				Usage: "Define the timeout for the hook execution.",
				Value: 2 * time.Minute,
			},
		},
	}
}

const rootPathWarningMessage = `!!!! HEADS UP !!!!

Your account credentials have been saved in your
configuration directory at "%s".

You should make a secure backup of this folder now. This
configuration directory will also contain certificates and
private keys obtained from the ACME server so making regular
backups of this folder is ideal.
`

func run(ctx *cli.Context) error {
	accountsStorage := NewAccountsStorage(ctx)

	account, keyType := setupAccount(ctx, accountsStorage)

	client := setupClient(ctx, account, keyType)

	if account.Registration == nil {
		reg, err := register(ctx, client)
		if err != nil {
			log.Fatalf("Could not complete registration\n\t%v", err)
		}

		account.Registration = reg
		if err = accountsStorage.Save(account); err != nil {
			log.Fatal(err)
		}

		fmt.Printf(rootPathWarningMessage, accountsStorage.GetRootPath())
	}

	certsStorage := NewCertificatesStorage(ctx)
	certsStorage.CreateRootFolder()

	obtainCertificate(ctx, client, certsStorage)

	meta := map[string]string{
		hookEnvAccountEmail: account.Email,
	}

	return launchHook(ctx.String(flgRunHook), ctx.Duration(flgRunHookTimeout), meta)
}

func register(ctx *cli.Context, client *lego.Client) (*registration.Resource, error) {
	accepted := true

	if ctx.Bool(flgEAB) {
		kid := ctx.String(flgKID)
		hmacEncoded := ctx.String(flgHMAC)

		if kid == "" || hmacEncoded == "" {
			log.Fatalf("Requires arguments --%s and --%s.", flgKID, flgHMAC)
		}

		return client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: accepted,
			Kid:                  kid,
			HmacEncoded:          hmacEncoded,
		})
	}

	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func obtainCertificate(ctx *cli.Context, client *lego.Client, certificatesStorage *CertificatesStorage) error {
	var wg sync.WaitGroup

	domains := ctx.StringSlice(flgDomains)

	for _, domain := range domains {
		wg.Add(1)
		go obtainDomainCertificate(domain, ctx, client, certificatesStorage, &wg)
	}

	wg.Wait()

	return nil
}

func obtainDomainCertificate(domain string, ctx *cli.Context, client *lego.Client, certificatesStorage *CertificatesStorage, wg *sync.WaitGroup) error {

	defer wg.Done()

	var challengeType string
	if ctx.Bool(flgHTTP) {
		challengeType = "http-01"
	}

	domaintype := "ip"
	ip := net.ParseIP(domain)
	if ip == nil {
		domaintype = "dns"
	}

	var identifiers = []acme.Identifier{
		{
			Type:          domaintype,
			Value:         domain,
			ChallengeType: challengeType,
		},
	}

	req := acme.CertRequest{
		CommonName: domain,
	}

	eccCsr, ecckey, err := acme.CreateCSR(acme.ECC, 256, req)
	if err != nil {
		log.Fatalf("Failed to create ECC CSR: %v", err)
	}

	sm2SignCsr, signkey, err := acme.CreateCSR(acme.SM2, 256, req)
	if err != nil {
		log.Fatalf("Failed to create SM2 CSR: %v", err)
	}

	sm2EncCsr, enckey, err := acme.CreateCSR(acme.SM2, 256, req)
	if err != nil {
		log.Fatalf("Failed to create SM2 CSR: %v", err)
	}

	request := api.OrderWithCsrOptions{
		Csr:         eccCsr,
		CsrSign:     sm2SignCsr,
		CsrEncrypt:  sm2EncCsr,
		Identifiers: identifiers,
	}

	order, err := client.Certificate.ObtainForSm2(domain, request)

	if err != nil {
		log.Logger.Fatalf("Failed to send order request: %v", err)
		return err
	}

	certs, err := client.Certificate.GetSm2OrderStatus(domain, order)
	if err != nil {
		log.Logger.Fatalf("Failed to get order status request: %v", err)
	}

	for k, v := range certs {
		// fmt.Printf("url:%s\n %s%s\n\n", k, string(v.Cert), string(v.Issuer))
		bys := append(v.Cert, v.Issuer...)

		switch {
		case k == "intcert":
			certificatesStorage.SaveCertificate(domain, "_ecc", bys, ecckey)
		case k == "sm2sign":
			certificatesStorage.SaveCertificate(domain, "_sign", bys, signkey)
		case k == "sm2enc":
			certificatesStorage.SaveCertificate(domain, "_enc", bys, enckey)
		}
	}

	return nil
}
