package server

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	httpAddr         string
	httpsAddr        string
	tunnelAddr       string
	reverseProxyAddr string
	domain           string
	tunnelTLSCrt     string
	tunnelTLSKey     string
	httpsTLSCrt      string
	httpsTLSKey      string
	logto            string
}

func parseArgs() *Options {
	httpAddr := flag.String("httpAddr", ":80", "Public address for HTTP connections, empty string to disable")
	httpsAddr := flag.String("httpsAddr", ":443", "Public address listening for HTTPS connections, emptry string to disable")
	tunnelAddr := flag.String("tunnelAddr", ":4443", "Public address listening for srvdir client")
	reverseProxyAddr := flag.String("reverseProxyAddr", "", "Address to listen for HTTP/HTTPS connections from an upstream reverse proxy")
	domain := flag.String("domain", "srvdir.net", "Domain where the public folders are available")
	tunnelTLSCrt := flag.String("tunnelTLSCrt", "", "Path to a TLS certificate file")
	tunnelTLSKey := flag.String("tunnelTLSKey", "", "Path to a TLS key file")
	httpsTLSCrt := flag.String("httpsTLSCrt", "", "Path to a TLS certificate file for the HTTPS server")
	httpsTLSKey := flag.String("httpsTLSKey", "", "Path to a TLS key file for the HTTPS server")
	logto := flag.String("log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	v := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	return &Options{
		httpAddr:         *httpAddr,
		httpsAddr:        *httpsAddr,
		tunnelAddr:       *tunnelAddr,
		reverseProxyAddr: *reverseProxyAddr,
		domain:           *domain,
		tunnelTLSCrt:     *tunnelTLSCrt,
		tunnelTLSKey:     *tunnelTLSKey,
		httpsTLSCrt:      *httpsTLSCrt,
		httpsTLSKey:      *httpsTLSKey,
		logto:            *logto,
	}
}
