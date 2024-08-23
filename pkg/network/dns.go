package network

import (
	"errors"
	"net"
)

func CheckCNAMERecord(domain string) (bool, error) {
	cNames, err := net.LookupCNAME(domain)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			if !dnsErr.IsNotFound {
				return false, err
			}
		}
	}

	if len(cNames) > 0 && cNames != domain+"." {
		return true, nil
	}

	return false, nil
}
