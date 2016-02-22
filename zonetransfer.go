package sonar

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func ZoneTransfer(domain string) Results {
	results := NewResultSet()
	fqdn := dns.Fqdn(domain)

	servers, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server.Host, "53"))
		if err != nil {
			log.Println(err)
			continue
		}

		for envelope := range answerChan {
			if envelope.Error != nil {
				log.Println(envelope.Error)
				break
			}

			for _, rr := range envelope.RR {
				switch v := rr.(type) {
				case *dns.A:
					results.Add(strings.TrimRight(v.Header().Name, "."), v.A.String())
				case *dns.AAAA:
					results.Add(strings.TrimRight(v.Header().Name, "."), v.AAAA.String())
				default:
				}
			}
		}
	}

	return results.Results()
}
