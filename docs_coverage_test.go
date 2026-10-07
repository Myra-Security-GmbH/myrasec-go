package myrasec

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"
)

const (
	documentedRoutesFile = "testdata/documented-routes.txt"
	pendingRoutesFile    = "testdata/pending-routes.txt"
)

var actionPlaceholders = strings.NewReplacer("%d", "*", "%s", "*", "%v", "*")

func TestEveryOperationIsDocumented(t *testing.T) {
	documented := readRoutes(t, documentedRoutesFile)
	pending := readRoutes(t, pendingRoutesFile)

	if len(documented) == 0 {
		t.Fatalf("%s lists no routes, run ./scripts/update-documented-routes.sh", documentedRoutesFile)
	}

	methods := initializeMethods()
	for _, name := range sortedOperationNames(methods) {
		route := documentedRoute(methods[name])
		_, isDocumented := documented[route]
		_, isPending := pending[route]

		if !isDocumented && !isPending {
			t.Errorf("SDK operation %q calls %q, which neither %s nor %s lists", name, route, documentedRoutesFile, pendingRoutesFile)
		}
	}
}

func TestPendingRoutesAreNotDocumentedYet(t *testing.T) {
	documented := readRoutes(t, documentedRoutesFile)

	for _, route := range sortedRoutes(readRoutes(t, pendingRoutesFile)) {
		if _, ok := documented[route]; ok {
			t.Errorf("%q is documented in production now, remove it from %s", route, pendingRoutesFile)
		}
	}
}

func TestPendingRoutesAreCalledBySDK(t *testing.T) {
	called := map[string]struct{}{}
	for _, method := range initializeMethods() {
		called[documentedRoute(method)] = struct{}{}
	}

	for _, route := range sortedRoutes(readRoutes(t, pendingRoutesFile)) {
		if _, ok := called[route]; !ok {
			t.Errorf("no SDK operation calls %q, remove it from %s", route, pendingRoutesFile)
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

func sortedOperationNames(methods map[string]APIMethod) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

func sortedRoutes(routes map[string]struct{}) []string {
	sorted := make([]string, 0, len(routes))
	for route := range routes {
		sorted = append(sorted, route)
	}
	sort.Strings(sorted)

	return sorted
}

func readRoutes(t *testing.T, path string) map[string]struct{} {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	defer file.Close()

	routes := map[string]struct{}{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "#")
		if line = strings.TrimSpace(line); line != "" {
			routes[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}

	return routes
}
