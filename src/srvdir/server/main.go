package server

import (
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel/server"
	"github.com/inconshreveable/go-tunnel/server/binder"
	tunneltls "github.com/inconshreveable/go-tunnel/tls"
	"time"
)

const (
	muxTimeout = 10 * time.Second
	version    = "0.1"
)

func Main() {
	opts := parseArgs()

	// set up logging
	log.LogTo(opts.logto)

	// load the tunnel TLS
	tunnelTLSConfig, err := tunneltls.ServerConfig(opts.tunnelTLSCrt, opts.tunnelTLSKey)
	if err != nil {
		panic(err)
	}

	// setup http and https binders
	var binders server.Binders = make(server.Binders)
	if opts.reverseProxyAddr != "" {
		if binders["http"], binders["https"], err = binder.NewReverseProxyBinder(opts.reverseProxyAddr, opts.domain, muxTimeout); err != nil {
			panic(err)
		}
	} else {
		if opts.httpAddr != "" {
			if binders["http"], err = binder.NewHTTPBinder(opts.httpAddr, opts.domain, muxTimeout); err != nil {
				panic(err)
			}
		}

		if opts.httpsAddr != "" {
			httpsTLSConfig, err := tunneltls.ServerConfig(opts.httpsTLSCrt, opts.httpsTLSKey)
			if err != nil {
				panic(err)
			}

			if binders["https"], err = binder.NewHTTPSBinder(opts.httpsAddr, opts.domain, muxTimeout, httpsTLSConfig); err != nil {
				panic(err)
			}
		}
	}

	server, err := server.ServeTLS("tcp", opts.tunnelAddr, tunnelTLSConfig, binders)
	if err != nil {
		panic(err)
	}

	server.SessionHooks = new(SessionHooks)

	server.Run()
}
