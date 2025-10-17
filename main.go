package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
)

type options struct {
	ProxyAddr       string `short:"a" long:"addr" description:"Proxy host address" default:"127.0.0.1:8080"`
	EdgeGridFile    string `short:"r" long:"file" description:"Location of EdgeGrid file" default:"~/.edgerc"`
	EdgeGridSection string `short:"s" long:"section" description:"Section of EdgeGrid file" default:"default"`
	AccountKey      string `short:"k" long:"key" env:"EDGEGRID_ACCOUNT_KEY" description:"Account switch key"`
	Host            string `long:"host" env:"EDGEGRID_HOST" description:"EdgeGrid Host"`
	ClientToken     string `long:"client-token" env:"EDGEGRID_CLIENT_TOKEN" description:"EdgeGrid ClientToken"`
	ClientSecret    string `long:"client-secret" env:"EDGEGRID_CLIENT_SECRET" description:"EdgeGrid ClientSecret"`
	AccessToken     string `long:"access-token" env:"EDGEGRID_ACCESS_TOKEN" description:"EdgeGrid AccessToken"`
	ProxyTLSCert    string `long:"tls-crt" description:"Proxy TLS/SSL certificate file path"`
	ProxyTLSKey     string `long:"tls-key" description:"Proxy TLS/SSL key file path"`

	proxyScheme string `no-flag:"true"`
}

func run(opts options) error {
	egpath, err := homedir.Expand(opts.EdgeGridFile)
	if err != nil {
		return err
	}

	var edgerc *edgegrid.Config
	if _, err := os.Stat(egpath); err == nil {
		log.Printf("Reading the credential in the '%s' section from %s", opts.EdgeGridSection, opts.EdgeGridFile)
		edgerc, err = edgegrid.New(
			edgegrid.WithFile(egpath),
			edgegrid.WithSection(opts.EdgeGridSection),
		)
		if err != nil {
			return err
		}
	} else {
		edgerc, _ = edgegrid.New()
		if opts.Host != "" {
			edgerc.Host = opts.Host
		}
		if opts.ClientToken != "" {
			edgerc.ClientToken = opts.ClientToken
		}
		if opts.ClientSecret != "" {
			edgerc.ClientSecret = opts.ClientSecret
		}
		if opts.AccessToken != "" {
			edgerc.AccessToken = opts.AccessToken
		}
	}
	if edgerc.Host == "" || edgerc.ClientToken == "" || edgerc.ClientSecret == "" || edgerc.AccessToken == "" {
		return fmt.Errorf("failed to load an edgegrid configuration")
	}

	if opts.AccountKey != "" {
		edgerc.AccountKey = opts.AccountKey
	}

	opts.proxyScheme = "http"
	if opts.ProxyTLSCert != "" && opts.ProxyTLSKey != "" {
		opts.proxyScheme = "https"
	}

	apiHost := &url.URL{Scheme: "https", Host: edgerc.Host}
	egproxy := httputil.NewSingleHostReverseProxy(apiHost)
	director := egproxy.Director

	egproxy.Director = func(req *http.Request) {
		req.Host = apiHost.Host
		director(req)

		//sign request
		edgerc.SignRequest(req)
		log.Printf("%s %s", req.Method, req.URL.String())
	}

	egproxy.ModifyResponse = func(resp *http.Response) error {
		loc := resp.Header.Get("Location")
		if loc == "" {
			return nil
		}

		u, err := url.Parse(loc)
		if err != nil {
			return nil
		}

		u.Scheme = opts.proxyScheme
		u.Host = opts.ProxyAddr

		//rewrite redirects
		resp.Header.Set("Location", u.String())
		return nil
	}

	log.Printf("EdgeGrid ClientToken: %s", edgerc.ClientToken)
	if opts.AccountKey != "" {
		log.Printf("EdgeGrid AccountSwitchKey: %s", edgerc.AccountKey)
	}

	log.Printf("Starting EdgeGrid proxy on %s://%s", opts.proxyScheme, opts.ProxyAddr)
	http.Handle("/", egproxy)

	if opts.proxyScheme == "https" {
		return http.ListenAndServeTLS(opts.ProxyAddr, opts.ProxyTLSCert, opts.ProxyTLSKey, nil)
	}
	return http.ListenAndServe(opts.ProxyAddr, nil)
}

func main() {
	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
