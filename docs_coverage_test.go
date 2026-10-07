package myrasec

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"
)

const documentedRoutesFile = "testdata/documented-routes.txt"

var actionPlaceholders = strings.NewReplacer("%d", "*", "%s", "*", "%v", "*")

func TestEveryOperationIsDocumented(t *testing.T) {
	documented := readDocumentedRoutes(t)
	methods := initializeMethods()

	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		route := documentedRoute(methods[name])
		if _, ok := documented[route]; !ok {
			t.Errorf("SDK operation %q calls %q, which %s does not list", name, route, documentedRoutesFile)
		}
	}
}

func TestDocumentedRoute(t *testing.T) {
	tests := map[string]struct {
		method   APIMethod
		expected string
	}{
		"no placeholder":                  {APIMethod{Method: "GET", Action: "domains"}, "GET /domains"},
		"numeric and string placeholders": {APIMethod{Method: "PUT", Action: "domain/%d/%s/settings"}, "PUT /domain/*/*/settings"},
		"query string is dropped":         {APIMethod{Method: "GET", Action: "waiting-rooms?domainId=%d"}, "GET /waiting-rooms"},
		"generic placeholder":             {APIMethod{Method: "DELETE", Action: "tag/%v"}, "DELETE /tag/*"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if route := documentedRoute(tc.method); route != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, route)
			}
		})
	}
}

func documentedRoute(method APIMethod) string {
	action, _, _ := strings.Cut(method.Action, "?")

	return method.Method + " /" + actionPlaceholders.Replace(action)
}

func readDocumentedRoutes(t *testing.T) map[string]struct{} {
	t.Helper()

	file, err := os.Open(documentedRoutesFile)
	if err != nil {
		t.Fatalf("cannot read %s: %v", documentedRoutesFile, err)
	}
	defer file.Close()

	routes := map[string]struct{}{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			routes[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("cannot read %s: %v", documentedRoutesFile, err)
	}

	if len(routes) == 0 {
		t.Fatalf("%s lists no routes, run ./scripts/update-documented-routes.sh", documentedRoutesFile)
	}

	return routes
}
