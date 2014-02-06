package main

import (
	"fmt"
	"github.com/inconshreveable/go-tunnel"
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/tls"
	"net"
	"net/http"
	"os"
)

func main() {
	// parse command line opts
	opts, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set up logging
	if opts.logto != "" {
		log.LogTo(opts.logto)
	}

	// set up tls configuration
	tlsName, _, err := net.SplitHostPort(opts.serverAddr)
	if err != nil {
		fmt.Printf("Failed to parse server address: %v: %v\n", opts.serverAddr, err)
		os.Exit(1)
	}

	tlsCfg, err := tls.ClientTrusted(tlsName)
	if err != nil {
		fmt.Printf("Failed to create TLS configuration for name %v: %v\n", tlsName, err)
		os.Exit(1)
	}

	c, err := tunnel.DialTLSReconnecting("tcp", opts.serverAddr, tlsCfg, nil)
	if err != nil {
		fmt.Printf("Failed to setup tunnel connection to %v: %v\n", opts.serverAddr, err)
		os.Exit(1)
	}

	routes := make([]Route, len(opts.dirs))

	for i, d := range opts.dirs {
		httpOpts := &proto.HTTPOptions{Subdomain: d.subdomain, Auth: opts.auth}
		tun, err := c.ListenHTTPS(httpOpts, nil)
		if err != nil {
			fmt.Printf("Failed to listen on subdomain '%v': %v\n", d.subdomain, err)
			os.Exit(1)
		}

		path := d.path
		go func() {
			log.Info("Serving %s on %s", path, tun.Addr().String())
			fs := FileServer(Dir(path), opts.index, true, opts.tmpl)
			if err := http.Serve(tun, fs); err != nil {
				fmt.Printf("Failed to start static http server for directory %v: %v\n", path, err)
				os.Exit(1)
			}
		}()

		routes[i] = Route{
			addr: tun.Addr().String(),
			path: path,
		}
	}

	// run the console UI
	if opts.logto != "stdout" {
		ui(routes)
	}
}
