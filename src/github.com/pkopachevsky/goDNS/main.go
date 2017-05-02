package main

import (
	"github.com/miekg/dns"
	"os"
	"net"
	"github.com/prometheus/common/log"
	"fmt"
)
func main() {
	domain := os.Args[1];
	fmt.Printf("Query for domain %s\n", domain)
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	m.RecursionDesired = true;

	fmt.Printf("Requesting %s:%s\n", config.Servers[0], config.Port)
	hostport := net.JoinHostPort(config.Servers[0], config.Port)
	r, _, err := c.Exchange(m, hostport)

	if r == nil {
		log.Fatalf("*** error: %s\n", err.Error())
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Fatalf("*** invalid answer name %s after MX query for %s\n", domain, domain)
	}

	for _,a := range r.Answer {
		fmt.Printf("%v\n", a)
	}
}