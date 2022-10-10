package myrasec

import (
	"fmt"
	"strconv"
	"strings"
)

// FetchDomainForSubdomainName returns the Domain for the passed subdomain (name)
func (api *API) FetchDomainForSubdomainName(subdomain string) (*Domain, error) {
	if IsGeneralDomainName(subdomain) {
		var parts []string
		name := RemoveTrailingDot(subdomain)
		if strings.HasPrefix(name, "ALL-") {
			parts = strings.Split(name, "ALL-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("wrong format for ALL-<DOMAIN_ID> annotation")
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			return api.GetDomain(id)
		}

		parts = strings.Split(name, "ALL:")
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong format for ALL:<DOMAIN_NAME> annotation")
		}

		return api.FetchDomain(parts[1])
	}

	maxRetries := 2
	retries := 0
	for {
		subdomains, err := api.ListAllSubdomains(map[string]string{"search": subdomain})
		if err != nil {
			return nil, err
		}

		domainNames := make(map[string]bool)
		for _, s := range subdomains {
			domainNames[s.DomainName] = true
		}

		for dn := range domainNames {
			domains, err := api.ListDomains(map[string]string{"search": dn})
			if err != nil {
				return nil, err
			}

			for _, d := range domains {
				vhosts, err := api.ListAllSubdomainsForDomain(d.ID, map[string]string{"search": subdomain})
				if err != nil {
					return nil, err
				}

				for _, vh := range vhosts {
					if EnsureTrailingDot(vh.Label) == EnsureTrailingDot(subdomain) {
						return &d, nil
					}
				}
			}
		}

		retries++
		if retries >= maxRetries {
			break
		}

		api.PruneCache()
	}
	return nil, fmt.Errorf("unable to find domain for passed subdomain")
}

// FetchDomain returns the Domain for the passed domain (name)
func (api *API) FetchDomain(domain string) (*Domain, error) {

	maxRetries := 2
	retries := 0
	for {
		domains, err := api.ListDomains(map[string]string{"search": domain})
		if err != nil {
			return nil, err
		}

		for _, d := range domains {
			if d.Name == domain {
				return &d, nil
			}
		}

		d, err := api.FetchDomainForSubdomainName(domain)
		if err != nil {
			return nil, fmt.Errorf("unable to find domain for passed domain name [%s]", domain)
		}
		if d != nil {
			return d, nil
		}

		retries++
		if retries >= maxRetries {
			break
		}

		api.PruneCache()
	}
	return nil, nil
}

// EnsureTrailingDot ensures and returns the passed subdomain with a trailing dot
func EnsureTrailingDot(subdomain string) string {
	return RemoveTrailingDot(subdomain) + "."
}

// RemoveTrailingDot removes and returns the trailing dot from the passed subdomain name
func RemoveTrailingDot(subdomain string) string {
	return strings.TrimRight(subdomain, ".")
}

// IsGeneralDomainName checks if the passed name starts with ALL- or ALL:
func IsGeneralDomainName(name string) bool {
	name = strings.ToUpper(name)
	return strings.HasPrefix(name, "ALL-") || strings.HasPrefix(name, "ALL:")
}
