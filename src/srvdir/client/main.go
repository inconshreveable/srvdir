package client

import (
	"fmt"
	"github.com/inconshreveable/go-tunnel"
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/tls"
	"net"
	"net/http"
	"os"
	srvdir_proto "srvdir/proto"
)

const version = "0.1"

func Main() {
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

	// load configuration file
	config, err := LoadConfiguration(opts.configPath)
	if err != nil {
		fmt.Printf("Error loading confuration: %v", err)
		os.Exit(1)
	}

	// override authtoken with command line option
	if opts.authtoken != "" {
		config.AuthToken = opts.authtoken
	}

	// connect to srvdir service
	authExtra := srvdir_proto.NewAuthExtra(config.AuthToken, version)
	c, err := tunnel.DialTLSReconnecting("tcp", opts.serverAddr, tlsCfg, authExtra)
	if err != nil {
		fmt.Printf("Failed to setup tunnel connection to %v: %v\n", opts.serverAddr, err)
		os.Exit(1)
	}

	if err := SaveAuthToken(opts.configPath, config.AuthToken); err != nil {
		log.Warn("Failed to save authtoken to config file: %v", err)
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
			fs := FileServer(Dir(path), opts.index, opts.readOnly, opts.tmpl)
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
	} else {
		select {}
	}
}
