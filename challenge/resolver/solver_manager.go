package resolver

import (
	"time"

	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/acme/api"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/log"
)

type byType []acme.Challenge

func (a byType) Len() int           { return len(a) }
func (a byType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byType) Less(i, j int) bool { return a[i].Type > a[j].Type }

type SolverManager struct {
	core    *api.Core
	solvers map[challenge.Type]solver
}

func NewSolversManager(core *api.Core) *SolverManager {
	return &SolverManager{
		solvers: map[challenge.Type]solver{},
		core:    core,
	}
}

// SetHTTP01Provider specifies a custom provider p that can solve the given HTTP-01 challenge.
func (c *SolverManager) SetHTTP01Provider(p challenge.Provider, opts ...http01.ChallengeOption) error {
	c.solvers[challenge.HTTP01] = http01.NewChallenge(c.core, validate, p, opts...)
	return nil
}

// SetTLSALPN01Provider specifies a custom provider p that can solve the given TLS-ALPN-01 challenge.
func (c *SolverManager) SetTLSALPN01Provider(p challenge.Provider, opts ...tlsalpn01.ChallengeOption) error {
	c.solvers[challenge.TLSALPN01] = tlsalpn01.NewChallenge(c.core, validate, p, opts...)
	return nil
}

// SetDNS01Provider specifies a custom provider p that can solve the given DNS-01 challenge.
func (c *SolverManager) SetDNS01Provider(p challenge.Provider, opts ...dns01.ChallengeOption) error {
	c.solvers[challenge.DNS01] = dns01.NewChallenge(c.core, validate, p, opts...)
	return nil
}

// Remove removes a challenge type from the available solvers.
func (c *SolverManager) Remove(chlgType challenge.Type) {
	delete(c.solvers, chlgType)
}

// Checks all challenges from the server in order and returns the first matching solver.
func (c *SolverManager) chooseSolver(authz acme.Authorization) solver {
	// Allow to have a deterministic challenge order

	// sort.Sort(byType(authz.Challenges))

	domain := challenge.GetTargetedDomain(authz)
	for _, chlg := range authz.Challenges {
		if solvr, ok := c.solvers[challenge.Type(chlg.Type)]; ok {
			log.Infof("[%s] acme: use %s solver", domain, chlg.Type)
			return solvr
		}

		log.Infof("[%s] acme: Could not find solver for: %s", domain, chlg.Type)
	}

	return nil
}

func validate(core *api.Core, domain string, chlg acme.Challenge) error {

	icout := 180
	status := ""
	for {
		time.Sleep(30 * time.Second)
		o, err := core.Orders.GetOrderStatus(chlg.URL)
		if err != nil {
			log.Warnf("[%s] acme: get order status %s", domain, err)
		}

		if status != o.Status {
			log.Infof("[%s] acme: get order status %s", domain, o.Status)
			status = o.Status
		}

		if o.Status == "valid" {
			break
		}

		icout--
		if icout < 0 {
			log.Fatalf("[%s] acme: get order status timeout", domain)
		}
	}

	return nil
}
