package main

import (
	"github.com/inconshreveable/go-tunnel/tls"
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel"
	"net/http"
)

func main() {
	opts, err := parseArgs()
	if err != nil {
		panic(err)
	}

	if opts.logto != "" {
		log.LogTo(opts.logto)
	}

	tlsCfg, err := tls.ClientTrusted("v1.srvdir.net")
	if err != nil {
		panic(err)
	}

	c, err := tunnel.DialTLSReconnecting("tcp", opts.serverAddr, tlsCfg, nil)
	if err != nil {
		panic(err)
	}

	routes := make([]Route, len(opts.dirs))

	for i, d := range opts.dirs {
		httpOpts := &proto.HTTPOptions{Subdomain: d.subdomain, Auth: opts.auth}
		tun, err := c.BindHTTPS(httpOpts, nil)
		if err != nil {
			panic(err)
		}

		path := d.path
		go func() {
			log.Info("Serving %s on %s", path, tun.Addr().String())
			fs := FileServer(Dir(path), opts.index, true, opts.tmpl)
			if err := http.Serve(tun, fs); err != nil {
				panic(err)
			}
		}()

		routes[i] = Route{
			addr: tun.Addr().String(),
			path: path,
		}
	}

	if opts.logto != "stdout" {
		ui(routes)
	}
}
