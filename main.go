package main

import (
	"flag"
	"log"
	"sync"

	"github.com/mailhog/MailHog-MTA/backend/auth"
	"github.com/mailhog/MailHog-MTA/backend/delivery"
	"github.com/mailhog/MailHog-MTA/backend/resolver"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/MailHog-MTA/smtp"
)

var conf *config.Config
var wg sync.WaitGroup

func configure() {
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()
}

func main() {
	configure()

	for _, s := range conf.Servers {
		wg.Add(1)
		go func(s *config.Server) {
			defer wg.Done()
			err := newServer(conf, s)
			if err != nil {
				log.Fatal(err)
			}
		}(s)
	}

	wg.Wait()
}

func newServer(cfg *config.Config, server *config.Server) error {
	s := &smtp.Server{
		BindAddr:        server.BindAddr,
		Hostname:        server.Hostname,
		PolicySet:       server.PolicySet,
		AuthBackend:     auth.Load(cfg, server),
		DeliveryBackend: delivery.Load(cfg, server),
		ResolverBackend: resolver.Load(cfg, server),
	}

	return s.Listen()
}
