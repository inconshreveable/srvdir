package main

import (
	"flag"
)

type Options struct {
	httpAddr   string
	httpsAddr  string
	tunnelAddr string
	reverseProxyAddr string
	domain     string
	tunnelTLSCrt     string
	tunnelTLSKey     string
	httpsTLSCrt string
	httpsTLSKey string
	logto      string
}

func parseArgs() *Options {
	httpAddr := flag.String("httpAddr", ":80", "Public address for HTTP connections, empty string to disable")
	httpsAddr := flag.String("httpsAddr", ":443", "Public address listening for HTTPS connections, emptry string to disable")
	tunnelAddr := flag.String("tunnelAddr", ":4443", "Public address listening for ngrok client")
	reverseProxyAddr := flag.String("reverseProxyAddr", "", "Address to listen for HTTP/HTTPS connections from an upstream reverse proxy")
	domain := flag.String("domain", "srvdir.net", "Domain where the public folders are available")
	tunnelTLSCrt := flag.String("tunnelTLSCrt", "", "Path to a TLS certificate file")
	tunnelTLSKey := flag.String("tunnelTLSKey", "", "Path to a TLS key file")
	httpsTLSCrt := flag.String("httpsTLSCrt", "", "Path to a TLS certificate file for the HTTPS server")
	httpsTLSKey := flag.String("httpsTLSKey", "", "Path to a TLS key file for the HTTPS server")
	logto := flag.String("log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")

	flag.Parse()

	return &Options{
		httpAddr:   *httpAddr,
		httpsAddr:  *httpsAddr,
		tunnelAddr: *tunnelAddr,
		reverseProxyAddr: *reverseProxyAddr,
		domain:     *domain,
		tunnelTLSCrt:     *tunnelTLSCrt,
		tunnelTLSKey:     *tunnelTLSKey,
		httpsTLSCrt:     *httpsTLSCrt,
		httpsTLSKey:     *httpsTLSKey,
		logto:      *logto,
	}
}
