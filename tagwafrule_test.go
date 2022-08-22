package myrasec

import (
	"testing"
)

func TestGetTagWAFRule(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tag/1/waf-rules/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "WAF Test", "logIdentifier": "waf_test", "subDomainName": "www.example.com", "domainId": 1, "direction": "in", "processNext": false, "description": "This is a testing WAF rule",
					"actions": [
						{"name": "Log", "type": "log", "availablePhases": 1}
					],
					"conditions": [
						{"alias": "Path", "availablePhases": 1, "category": "URL", "id": 1, "matchingType": "IREGEX", "name": "url", "value": "test"}
					],
					"expireDate": null, "tagId": 1
				}
			]}`,
			methods["getTagWAFRule"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rule, err := api.GetTagWAFRule(1, 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if rule.ID != 1 {
		t.Errorf("Expected to get WAF rule with ID [%d] but got [%d]", 1, rule.ID)
	}

	if rule.Name != "WAF Test" {
		t.Errorf("Expected to get WAF rule with Name [%s] but got [%s]", "WAF Test", rule.Name)
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

func TestListTagWAFRules(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tag/1/waf-rules",
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
			methods["listTagWAFRules"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rules, err := api.ListTagWAFRules(1)
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

		//if rule.SubDomainName != "www.example.com" {
		//	t.Errorf("Expected to get WAF rule with SubDomainName [%s] but got [%s]", "www.example.com", rule.SubDomainName)
		//}

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
