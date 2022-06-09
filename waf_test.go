package myrasec

import (
	"testing"
)

func TestFetchWAFRule(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/waf-rules/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "WAF Test", "logIdentifier": "waf_test", "subDomainName": "www.example.com", "domainId": 1, "direction": "in", "processNext": false, "description": "This is a testing WAF rule",
					"actions": [
						{"name": "Log", "type": "log", "availablePhases": 1}
					],
					"conditions": [
						{"alias": "Path", "availablePhases": 1, "category": "URL", "id": 1, "matchingType": "IREGEX", "name": "url", "value": "test"}
					],
					"expireDate": null
				}
			]}`,
			methods["fetchWAFRule"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rule, err := api.FetchWAFRule(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if rule.ID != 1 {
		t.Errorf("Expected to get WAF rule with ID [%d] but got [%d]", 1, rule.ID)
	}

	if rule.Name != "WAF Test" {
		t.Errorf("Expected to get WAF rule with Name [%s] but got [%s]", "WAF Test", rule.Name)
	}

	if rule.SubDomainName != "www.example.com" {
		t.Errorf("Expected to get WAF rule with SubDomainName [%s] but got [%s]", "www.example.com", rule.SubDomainName)
	}

	if rule.Direction != "in" {
		t.Errorf("Expected to get WAF rule with Direction [%s] but got [%s]", "in", rule.Direction)
	}

	if rule.ProcessNext != false {
		t.Errorf("Expected to get WAF rule with ProcessNext [%t] but got [%t]", false, rule.ProcessNext)
	}

	if rule.Description != "This is a testing WAF rule" {
		t.Errorf("Expected to get WAF rule with Description [%s] but got [%s]", "This is a testing WAF rule", rule.Description)
	}

	if len(rule.Actions) != 1 {
		t.Errorf("Expected to get WAF rule with one Action but got [%d] Actions", len(rule.Actions))
	}

	if len(rule.Conditions) != 1 {
		t.Errorf("Expected to get WAF rule with one Condition but got [%d] Conditions", len(rule.Conditions))
	}

	for _, a := range rule.Actions {
		if a.Name != "Log" {
			t.Errorf("Expected to get WAF rule with Action Name [%s] but got [%s]", "Log", a.Name)
		}

		if a.Type != "log" {
			t.Errorf("Expected to get WAF rule with Action Type [%s] but got [%s]", "log", a.Type)
		}

		if a.AvailablePhases != 1 {
			t.Errorf("Expected to get WAF rule with Action AvailablePhases [%d] but got [%d]", 1, a.AvailablePhases)
		}
	}

	for _, c := range rule.Conditions {
		if c.Alias != "Path" {
			t.Errorf("Expected to get WAF rule with Condition Alias [%s] but got [%s]", "Path", c.Alias)
		}

		if c.AvailablePhases != 1 {
			t.Errorf("Expected to get WAF rule with Condition AvailablePhases [%d] but got [%d]", 1, c.AvailablePhases)
		}

		if c.Category != "URL" {
			t.Errorf("Expected to get WAF rule with Condition Category [%s] but got [%s]", "URL", c.Category)
		}

		if c.MatchingType != "IREGEX" {
			t.Errorf("Expected to get WAF rule with Condition MatchingType [%s] but got [%s]", "IREGEX", c.MatchingType)
		}

		if c.Name != "url" {
			t.Errorf("Expected to get WAF rule with Condition Name [%s] but got [%s]", "url", c.Name)
		}

		if c.Value != "test" {
			t.Errorf("Expected to get WAF rule with Condition Value [%s] but got [%s]", "test", c.Value)
		}
	}

}

func TestListWAFRules(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/waf-rules",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "WAF Test", "logIdentifier": "waf_test", "subDomainName": "www.example.com", "domainId": 1, "direction": "in", "processNext": false, "description": "This is a testing WAF rule",
					"actions": [
						{"name": "Log", "type": "log", "availablePhases": 1}
					],
					"conditions": [
						{"alias": "Path", "availablePhases": 1, "category": "URL", "id": 1, "matchingType": "IREGEX", "name": "url", "value": "test"}
					],
					"expireDate": null
				}
			]}`,
			methods["listWAFRules"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rules, err := api.ListWAFRules(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	for _, rule := range rules {
		if rule.ID != 1 {
			t.Errorf("Expected to get WAF rule with ID [%d] but got [%d]", 1, rule.ID)
		}

		if rule.Name != "WAF Test" {
			t.Errorf("Expected to get WAF rule with Name [%s] but got [%s]", "WAF Test", rule.Name)
		}

		if rule.SubDomainName != "www.example.com" {
			t.Errorf("Expected to get WAF rule with SubDomainName [%s] but got [%s]", "www.example.com", rule.SubDomainName)
		}

		if rule.Direction != "in" {
			t.Errorf("Expected to get WAF rule with Direction [%s] but got [%s]", "in", rule.Direction)
		}

		if rule.ProcessNext != false {
			t.Errorf("Expected to get WAF rule with ProcessNext [%t] but got [%t]", false, rule.ProcessNext)
		}

		if rule.Description != "This is a testing WAF rule" {
			t.Errorf("Expected to get WAF rule with Description [%s] but got [%s]", "This is a testing WAF rule", rule.Description)
		}

		if len(rule.Actions) != 1 {
			t.Errorf("Expected to get WAF rule with one Action but got [%d] Actions", len(rule.Actions))
		}

		if len(rule.Conditions) != 1 {
			t.Errorf("Expected to get WAF rule with one Condition but got [%d] Conditions", len(rule.Conditions))
		}

		for _, a := range rule.Actions {
			if a.Name != "Log" {
				t.Errorf("Expected to get WAF rule with Action Name [%s] but got [%s]", "Log", a.Name)
			}

			if a.Type != "log" {
				t.Errorf("Expected to get WAF rule with Action Type [%s] but got [%s]", "log", a.Type)
			}

			if a.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF rule with Action AvailablePhases [%d] but got [%d]", 1, a.AvailablePhases)
			}
		}

		for _, c := range rule.Conditions {
			if c.Alias != "Path" {
				t.Errorf("Expected to get WAF rule with Condition Alias [%s] but got [%s]", "Path", c.Alias)
			}

			if c.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF rule with Condition AvailablePhases [%d] but got [%d]", 1, c.AvailablePhases)
			}

			if c.Category != "URL" {
				t.Errorf("Expected to get WAF rule with Condition Category [%s] but got [%s]", "URL", c.Category)
			}

			if c.MatchingType != "IREGEX" {
				t.Errorf("Expected to get WAF rule with Condition MatchingType [%s] but got [%s]", "IREGEX", c.MatchingType)
			}

			if c.Name != "url" {
				t.Errorf("Expected to get WAF rule with Condition Name [%s] but got [%s]", "url", c.Name)
			}

			if c.Value != "test" {
				t.Errorf("Expected to get WAF rule with Condition Value [%s] but got [%s]", "test", c.Value)
			}
		}

	}
}

func TestListWAFActions(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/waf/actions",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"availablePhases": 1, "forceCustomValues": 0, "name": "Allow", "type": "allow"},
				{"availablePhases": 1, "forceCustomValues": 0, "name": "Block", "type": "block"},
				{"availablePhases": 3, "forceCustomValues": 2, "name": "Add header", "type": "add_header"}
			]}`,
			methods["listWAFActions"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	actions, err := api.ListWAFActions()
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(actions) != 3 {
		t.Errorf("Expected to get [%d] WAF Actions but got [%d]", 3, len(actions))
	}

	for _, a := range actions {
		if a.Type == "allow" {
			if a.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF Action with AvailablePhases [%d] but got [%d]", 1, a.AvailablePhases)
			}
			if a.ForceCustomValues != false {
				t.Errorf("Expected to get WAF Action with ForceCustomValues [%t] but got [%t]", false, a.ForceCustomValues)
			}
			if a.Name != "Allow" {
				t.Errorf("Expected to get WAF Action with Name [%s] but got [%s]", "Allow", a.Name)
			}
		}

		if a.Type == "block" {
			if a.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF Action with AvailablePhases [%d] but got [%d]", 1, a.AvailablePhases)
			}
			if a.ForceCustomValues != false {
				t.Errorf("Expected to get WAF Action with ForceCustomValues [%t] but got [%t]", false, a.ForceCustomValues)
			}
			if a.Name != "Block" {
				t.Errorf("Expected to get WAF Action with Name [%s] but got [%s]", "Block", a.Name)
			}
		}

		if a.Type == "add_header" {
			if a.AvailablePhases != 3 {
				t.Errorf("Expected to get WAF Action with AvailablePhases [%d] but got [%d]", 3, a.AvailablePhases)
			}
			if a.ForceCustomValues != false {
				t.Errorf("Expected to get WAF Action with ForceCustomValues [%t] but got [%t]", true, a.ForceCustomValues)
			}
			if a.Name != "Add header" {
				t.Errorf("Expected to get WAF Action with Name [%s] but got [%s]", "Add header", a.Name)
			}
		}

	}

}

func TestListWAFConditions(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/waf/conditions",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"alias": "Custom header", "availablePhases": 3, "category": "HEADER", "forceCustomValues": true, "name": "custom_header", "value": ""},
				{"alias": "Host header", "availablePhases": 1, "category": "HEADER", "forceCustomValues": false, "name": "host", "value": ""},
				{"alias": "User-Agent header", "availablePhases": 1, "category": "HEADER", "forceCustomValues": false, "name": "user_agent", "value": ""}
			]}`,
			methods["listWAFConditions"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	conditions, err := api.ListWAFConditions()
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(conditions) != 3 {
		t.Errorf("Expected to get [%d] WAF Conditions but got [%d]", 3, len(conditions))
	}

	for _, c := range conditions {
		if c.Name == "custom_header" {
			if c.AvailablePhases != 3 {
				t.Errorf("Expected to get WAF Condition with AvailablePhases [%d] but got [%d]", 3, c.AvailablePhases)
			}
			if c.ForceCustomValues != true {
				t.Errorf("Expected to get WAF Condition with ForceCustomValues [%t] but got [%t]", true, c.ForceCustomValues)
			}
			if c.Value != "" {
				t.Errorf("Expected to get WAF Condition with Value [%s] but got [%s]", "", c.Value)
			}
			if c.Alias != "Custom header" {
				t.Errorf("Expected to get WAF Condition with Alias [%s] but got [%s]", "Custom header", c.Alias)
			}
			if c.Category != "HEADER" {
				t.Errorf("Expected to get WAF Condition with Category [%s] but got [%s]", "HEADER", c.Category)
			}
		}

		if c.Name == "host" {
			if c.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF Condition with AvailablePhases [%d] but got [%d]", 1, c.AvailablePhases)
			}
			if c.ForceCustomValues != false {
				t.Errorf("Expected to get WAF Condition with ForceCustomValues [%t] but got [%t]", false, c.ForceCustomValues)
			}
			if c.Value != "" {
				t.Errorf("Expected to get WAF Condition with Value [%s] but got [%s]", "", c.Value)
			}
			if c.Alias != "Host header" {
				t.Errorf("Expected to get WAF Condition with Alias [%s] but got [%s]", "Host header", c.Alias)
			}
			if c.Category != "HEADER" {
				t.Errorf("Expected to get WAF Condition with Category [%s] but got [%s]", "HEADER", c.Category)
			}
		}

		if c.Name == "user_agent" {
			if c.AvailablePhases != 1 {
				t.Errorf("Expected to get WAF Condition with AvailablePhases [%d] but got [%d]", 1, c.AvailablePhases)
			}
			if c.ForceCustomValues != false {
				t.Errorf("Expected to get WAF Condition with ForceCustomValues [%t] but got [%t]", false, c.ForceCustomValues)
			}
			if c.Value != "" {
				t.Errorf("Expected to get WAF Condition with Value [%s] but got [%s]", "", c.Value)
			}
			if c.Alias != "User-Agent header" {
				t.Errorf("Expected to get WAF Condition with Alias [%s] but got [%s]", "User-Agent header", c.Alias)
			}
			if c.Category != "HEADER" {
				t.Errorf("Expected to get WAF Condition with Category [%s] but got [%s]", "HEADER", c.Category)
			}
		}
	}
}
