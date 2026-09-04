package myrasec

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// FetchDomainForSubdomainNameContext returns the Domain for the passed subdomain (name)
func (api *API) FetchDomainForSubdomainNameContext(ctx context.Context, subdomain string) (*Domain, error) {
	if IsGeneralDomainName(subdomain) {
		if strings.HasPrefix(subdomain, "ALL-") {
			id, err := ExtractDomainIdFromGeneralDomainName(subdomain)
			if err != nil {
				return nil, err
			}
			return api.GetDomainContext(ctx, id)
		}

		var parts []string
		name := RemoveTrailingDot(subdomain)
		parts = strings.Split(name, "ALL:")
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong format for ALL:<DOMAIN_NAME> annotation")
		}

		return api.FetchDomainContext(ctx, parts[1])
	}

	maxRetries := 2
	retries := 0
	for {
		subdomains, err := api.ListAllSubdomainsContext(ctx, map[string]string{"search": subdomain})
		if err != nil {
			return nil, err
		}

		domainNames := make(map[string]bool)
		for _, s := range subdomains {
			domainNames[s.DomainName] = true
		}

		for dn := range domainNames {
			domains, err := api.ListDomainsContext(ctx, map[string]string{"search": dn})
			if err != nil {
				return nil, err
			}

			for _, d := range domains {
				vhosts, err := api.ListAllSubdomainsForDomainContext(ctx, d.ID, map[string]string{"search": subdomain})
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

// FetchDomainForSubdomainName is equivalent to FetchDomainForSubdomainNameContext with context.Background().
//
// Deprecated: use FetchDomainForSubdomainNameContext.
func (api *API) FetchDomainForSubdomainName(subdomain string) (*Domain, error) {
	return api.FetchDomainForSubdomainNameContext(context.Background(), subdomain)
}

// FetchDomainContext returns the Domain for the passed domain (name)
func (api *API) FetchDomainContext(ctx context.Context, domain string) (*Domain, error) {
	maxRetries := 2
	retries := 0
	for {
		domains, err := api.ListDomainsContext(ctx, map[string]string{"search": domain})
		if err != nil {
			return nil, err
		}

		for _, d := range domains {
			if d.Name == domain {
				return &d, nil
			}
		}

		d, err := api.FetchDomainForSubdomainNameContext(ctx, domain)
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

// FetchDomain is equivalent to FetchDomainContext with context.Background().
//
// Deprecated: use FetchDomainContext.
func (api *API) FetchDomain(domain string) (*Domain, error) {
	return api.FetchDomainContext(context.Background(), domain)
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
	name = RemoveTrailingDot(strings.ToUpper(name))
	if strings.HasPrefix(name, "ALL:") {
		return true
	}

	if strings.HasPrefix(name, "ALL-") {
		parts := strings.Split(name, "ALL-")
		if len(parts) != 2 {
			return false
		}
		_, err := strconv.Atoi(parts[1])
		return err == nil
	}

	return false
}

// ExtractDomainIdFromGeneralDomainName extracts the domainID from the general domain name annotation (ALL-1234.)
func ExtractDomainIdFromGeneralDomainName(generalDomainName string) (int, error) {
	if !IsGeneralDomainName(generalDomainName) {
		return 0, fmt.Errorf("passed generalDomainName has the wrong format")
	}

	var parts []string
	name := RemoveTrailingDot(generalDomainName)
	if !strings.HasPrefix(name, "ALL-") {
		return 0, fmt.Errorf("wrong format for ALL-<DOMAIN_ID> annotation")
	}

	parts = strings.Split(name, "ALL-")
	if len(parts) != 2 {
		return 0, fmt.Errorf("wrong format for ALL-<DOMAIN_ID> annotation")
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return id, nil
}
